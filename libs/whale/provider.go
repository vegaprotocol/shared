package whale

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/slack-go/slack"

	"code.vegaprotocol.io/shared/libs/num"
	"code.vegaprotocol.io/shared/libs/types"
	"code.vegaprotocol.io/shared/libs/whale/config"
	"code.vegaprotocol.io/vega/logging"
	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
)

type Provider struct {
	node             dataNode
	erc20            erc20Service
	account          types.AccountStream
	faucet           faucetClient
	slack            slacker
	ownerPrivateKeys map[string]string
	pendingDeposits  map[string]pendingDeposit
	mu               sync.Mutex
	topUpChan        chan types.TopUpRequest
	callTimeout      time.Duration
	log              *logging.Logger
}

type slacker struct {
	*slack.Client // TODO: abstract this out
	channelID     string
	enabled       bool
}

type pendingDeposit struct {
	amount    *num.Uint
	timestamp string
}

func NewProvider(
	log *logging.Logger,
	node dataNode,
	erc20 erc20Service,
	faucet faucetClient,
	account types.AccountStream,
	config *config.WhaleConfig,
) *Provider {
	p := &Provider{
		node:             node,
		erc20:            erc20,
		faucet:           faucet,
		account:          account,
		ownerPrivateKeys: config.OwnerPrivateKeys,
		topUpChan:        make(chan types.TopUpRequest),
		callTimeout:      time.Duration(config.SyncTimeoutSec) * time.Second,
		slack: slacker{
			Client:    slack.New(config.SlackConfig.BotToken, slack.OptionAppLevelToken(config.SlackConfig.AppToken)),
			channelID: config.SlackConfig.ChannelID,
			enabled:   config.SlackConfig.Enabled,
		},
		log: log.Named("WhaleProvider"),
	}

	go func() {
		for req := range p.topUpChan {
			req.ErrResp <- p.handleTopUp(req.Ctx, req.ReceiverName, req.ReceiverAddress, req.AssetID, req.Amount)
		}
	}()
	return p
}

func (p *Provider) TopUpChan() chan types.TopUpRequest {
	return p.topUpChan
}

func (p *Provider) handleTopUp(ctx context.Context, receiverName, receiverAddress, assetID string, amount *num.Uint) error {
	var err error
	defer func() {
		if err == nil || p.slack.enabled {
			if err := p.account.WaitForTopUpToFinalise(ctx, receiverAddress, assetID, amount, 0); err != nil {
				p.log.With(
					logging.String("receiver_address", receiverAddress),
					logging.String("asset_id", assetID),
					logging.String("amount", amount.String()),
				).Error("failed to finalise top up", logging.Error(err))
			}
		}
	}()

	// TODO: remove deposit slack request, once deposited
	if p.slack.enabled {
		if existDeposit, ok := p.getPendingDeposit(assetID); ok {
			newTimestamp, err := p.updateDan(ctx, assetID, receiverAddress, existDeposit.timestamp, existDeposit.amount)
			if err != nil {
				return fmt.Errorf("failed to update slack message: %s", err)
			}
			existDeposit.timestamp = newTimestamp
			existDeposit.amount = amount.Add(amount, existDeposit.amount)
			p.setPendingDeposit(assetID, existDeposit)
			return nil
		}
	}

	err = p.deposit(ctx, "Whale", receiverAddress, assetID, amount)
	if err == nil {
		return nil
	}

	p.log.With(
		logging.String("receiverName", receiverName),
		logging.String("receiverAddress", receiverAddress),
	).Warningf("Failed to deposit: %s", err)

	deposit := pendingDeposit{
		amount: amount,
	}

	if !p.slack.enabled {
		return fmt.Errorf("failed to deposit: %w", err)
	}

	p.log.Debug("Fallback to slacking Dan...")

	deposit.timestamp, err = p.slackDan(ctx, assetID, receiverAddress, amount)
	if err != nil {
		p.log.Error("Failed to slack Dan", logging.Error(err))
		return err
	}
	p.setPendingDeposit(assetID, deposit)
	return nil
}

func (p *Provider) deposit(ctx context.Context, receiverName, receiverAddress, assetID string, amount *num.Uint) error {
	asset, err := p.node.AssetByID(ctx, &dataapipb.GetAssetRequest{
		AssetId: assetID,
	})
	if err != nil {
		return fmt.Errorf("failed to get asset id: %w", err)
	}

	erc20 := asset.Details.GetErc20()
	if erc20 == nil {
		return fmt.Errorf("unsupported asset type")
	}

	ownerKey, err := p.getOwnerKeyForAsset(asset.Id)
	if err != nil {
		return fmt.Errorf("failed to get owner key: %w", err)
	}

	contractAddress := asset.Details.GetErc20().ContractAddress

	added, err := p.erc20.Deposit(
		ctx,
		ownerKey.privateKey,
		ownerKey.address,
		contractAddress,
		receiverAddress,
		amount,
	)
	if err != nil {
		return fmt.Errorf("failed to deposit %s %s coins to address '%s', name '%s': %w", amount.String(), asset.Details.Symbol, receiverAddress, receiverName, err)
	}

	if added.Int().LT(amount.Int()) {
		return fmt.Errorf("deposited less than requested amount")
	}

	return nil
}

func (p *Provider) getPendingDeposit(assetID string) (pendingDeposit, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.pendingDeposits == nil {
		p.pendingDeposits = make(map[string]pendingDeposit)
		return pendingDeposit{}, false
	}

	pending, ok := p.pendingDeposits[assetID]
	return pending, ok
}

func (p *Provider) setPendingDeposit(assetID string, pending pendingDeposit) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.pendingDeposits == nil {
		p.pendingDeposits = make(map[string]pendingDeposit)
	}

	p.pendingDeposits[assetID] = pending
}

func (p *Provider) Stake(ctx context.Context, _, receiverAddress, assetID string, amount *num.Uint, _ string) error {
	asset, err := p.node.AssetByID(ctx, &dataapipb.GetAssetRequest{
		AssetId: assetID,
	})
	if err != nil {
		return fmt.Errorf("failed to get asset id: %w", err)
	}
	erc20 := asset.Details.GetErc20()
	if erc20 == nil {
		return fmt.Errorf("asset is not erc20")
	}

	ownerKey, err := p.getOwnerKeyForAsset(asset.Id)
	if err != nil {
		return fmt.Errorf("failed to get owner for key '%s': %w", receiverAddress, err)
	}

	contractAddress := asset.Details.GetErc20().ContractAddress

	added, err := p.erc20.Stake(ctx, ownerKey.privateKey, ownerKey.address, contractAddress, receiverAddress, amount)
	if err != nil {
		return fmt.Errorf("failed to stake Vega token for '%s': %w", receiverAddress, err)
	}

	if added.Int().LT(amount.Int()) {
		return fmt.Errorf("staked less than requested amount")
	}

	return nil
}

type key struct {
	privateKey string
	address    string
}

func (p *Provider) getOwnerKeyForAsset(assetID string) (*key, error) {
	ownerPrivateKey, ok := p.ownerPrivateKeys[assetID]
	if !ok {
		return nil, fmt.Errorf("owner private key not configured for asset '%s'", assetID)
	}

	address, err := addressFromPrivateKey(ownerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get address from private key: %w", err)
	}

	return &key{
		privateKey: ownerPrivateKey,
		address:    address,
	}, nil
}

func addressFromPrivateKey(privateKey string) (string, error) {
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to convert owner private key hash into ECDSA: %w", err)
	}

	publicKeyECDSA, ok := key.Public().(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address, nil
}

const msgTemplate = `Hi @here! Whale wallet account with pub key %s needs %s coins of assetID %s, so that it can feed the hungry bots.`

func (p *Provider) slackDan(ctx context.Context, assetID, walletPubKey string, amount *num.Uint) (string, error) {
	p.log.With(
		logging.String("assetID", assetID),
		logging.String("walletPubKey", walletPubKey),
		logging.String("amount", amount.String()),
	).Debug("Slack post @hungry-bots")

	message := fmt.Sprintf(msgTemplate, walletPubKey, amount.String(), assetID)

	respChannel, respTimestamp, err := p.slack.PostMessageContext(
		ctx,
		p.slack.channelID,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		return "", err
	}

	p.log.With(
		logging.String("channel", respChannel),
		logging.String("timestamp", respTimestamp),
	).Debug("Slack message successfully sent")

	time.Sleep(time.Second * 5)

	_, _, _ = p.slack.PostMessageContext(
		ctx,
		p.slack.channelID,
		slack.MsgOptionText("I can wait...", false),
	)
	return respTimestamp, nil
}

func (p *Provider) updateDan(ctx context.Context, assetID, walletPubKey, oldTimestamp string, amount *num.Uint) (string, error) {
	p.log.With(
		logging.String("assetID", assetID),
		logging.String("walletPubKey", walletPubKey),
		logging.String("amount", amount.String()),
	).Debug("Slack update @hungry-bots")

	message := fmt.Sprintf(msgTemplate, walletPubKey, amount.String(), assetID)

	respChannel, respTimestamp, _, err := p.slack.UpdateMessageContext(
		ctx,
		p.slack.channelID,
		oldTimestamp,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		return "", err
	}

	p.log.With(
		logging.String("channel", respChannel),
		logging.String("timestamp", respTimestamp),
	).Debug("Slack message successfully updated ")
	return respTimestamp, nil
}
