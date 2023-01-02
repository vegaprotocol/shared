package whale

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"code.vegaprotocol.io/shared/libs/cache"
	"code.vegaprotocol.io/shared/libs/num"
	"code.vegaprotocol.io/shared/libs/types"
	"code.vegaprotocol.io/shared/libs/wallet"
	"code.vegaprotocol.io/shared/libs/whale/config"
	vtypes "code.vegaprotocol.io/vega/core/types"
	"code.vegaprotocol.io/vega/logging"
	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	commV1 "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	"code.vegaprotocol.io/vega/protos/vega/wallet/v1"
)

type Service struct {
	node          dataNode
	wallet        wallet.WalletV2
	account       accountService
	accountStream types.AccountStream
	faucet        faucetClient

	topUpChan    chan types.TopUpRequest
	walletConfig *config.WhaleConfig
	log          *logging.Logger
}

func NewService(
	log *logging.Logger,
	dataNode dataNode,
	wallet wallet.WalletV2,
	account accountService,
	accountStream types.AccountStream,
	faucet faucetClient,
	config *config.WhaleConfig,
) *Service {
	w := &Service{
		node:   dataNode,
		wallet: wallet,
		faucet: faucet,

		topUpChan:     make(chan types.TopUpRequest),
		account:       account,
		accountStream: accountStream,
		walletConfig:  config,
		log:           log.Named("Whale"),
	}
	go func() {
		for req := range w.topUpChan {
			req.ErrResp <- w.handleTopUp(
				req.Ctx,
				req.ReceiverName,
				req.ReceiverAddress,
				req.AssetID,
				req.Amount,
				req.From,
			)
			close(req.ErrResp)
		}
	}()
	return w
}

func (w *Service) TopUpChan() chan types.TopUpRequest {
	return w.topUpChan
}

func (w *Service) handleTopUp(ctx context.Context, receiverName, receiverAddress, assetID string, amount *num.Uint, from string) error {
	w.log.With(logging.String("receiverName", receiverName)).Debug("Top up...")

	if assetID == "" {
		return fmt.Errorf("assetID is empty for bot '%s'", receiverName)
	}

	if receiverAddress == w.wallet.PublicKey() {
		return fmt.Errorf("whale and bot address cannot be the same")
	}

	if err := w.topUp(ctx, receiverName, receiverAddress, assetID, amount, from); err != nil {
		return fmt.Errorf("failed to top up: %w", err)
	}

	w.log.With(
		logging.String("receiverName", receiverName),
		logging.String("receiverPubKey", receiverAddress),
		logging.AssetID(assetID),
		logging.String("amount", amount.String()),
	).Debug("Top-up sent")

	w.log.With(logging.String("name", receiverName)).Debugf("%s: Waiting for top-up...", from)

	if err := w.accountStream.WaitForTopUpToFinalise(ctx, receiverAddress, assetID, amount, 0); err != nil {
		return fmt.Errorf("failed to wait for top-up to finalise: %w", err)
	}
	w.log.With(logging.String("name", receiverName)).Debugf("%s: Top-up complete", from)

	return nil
}

func (w *Service) topUp(ctx context.Context, receiverName string, receiverAddress string, assetID string, amount *num.Uint, from string) error {
	asset, err := w.node.AssetByID(ctx, &dataapipb.GetAssetRequest{
		AssetId: assetID,
	})
	if err != nil {
		return fmt.Errorf("failed to get asset by id: %w", err)
	}

	ensureAmount := num.Zero().Mul(amount, num.NewUint(30))

	if builtin := asset.Details.GetBuiltinAsset(); builtin != nil {
		if err := w.depositBuiltin(ctx, assetID, receiverAddress, ensureAmount, builtin); err != nil {
			return errors.Wrap(err, "failed to deposit builtin")
		}
		return nil
	}

	// dp is 0 because the amount had already been corrected for the DP
	if err := w.account.EnsureBalance(ctx, assetID, cache.General, ensureAmount, 0, 100, from+">receiverNameWhale"); err != nil {
		return fmt.Errorf("failed to ensure enough funds: %w", err)
	}

	_, err = w.wallet.SendTransaction(ctx, &v1.SubmitTransactionRequest{
		Command: &v1.SubmitTransactionRequest_Transfer{
			Transfer: &commV1.Transfer{
				FromAccountType: vtypes.AccountTypeGeneral,
				To:              receiverAddress,
				ToAccountType:   vtypes.AccountTypeGeneral,
				Asset:           assetID,
				Amount:          amount.String(),
				Reference:       fmt.Sprintf("Bot '%s' Top-Up", receiverName),
				Kind:            &commV1.Transfer_OneOff{OneOff: &commV1.OneOffTransfer{}},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to top-up bot '%s': %w", receiverName, err)
	}
	return nil
}

func (w *Service) depositBuiltin(ctx context.Context, assetID, pubKey string, amount *num.Uint, builtin *vega.BuiltinAsset) error {
	maxFaucet, err := num.ConvertUint256(builtin.MaxFaucetAmountMint)
	if err != nil {
		return fmt.Errorf("failed to convert max faucet amount: %w", err)
	}

	if maxFaucet.GT(amount) {
		if ok, err := w.faucet.Mint(ctx, maxFaucet.String(), assetID, pubKey); err != nil {
			return fmt.Errorf("failed to mint: %w", err)
		} else if !ok {
			return fmt.Errorf("failed to mint")
		}
		return nil
	}

	times := int(new(num.Uint).Div(amount, maxFaucet).Uint64() + 1)
	totalMinted := new(num.Uint)

	// TODO: limit the time here!

	for i := 0; i < times; i++ {
		if ok, err := w.faucet.Mint(ctx, maxFaucet.String(), assetID, pubKey); err != nil {
			return fmt.Errorf("failed to mint: %w", err)
		} else if !ok {
			return fmt.Errorf("failed to mint")
		}

		totalMinted.Add(totalMinted, maxFaucet)

		time.Sleep(w.walletConfig.FaucetRateLimit)
		w.log.With(
			logging.AssetID(assetID),
			logging.PartyID(pubKey),
		).Infof("Minted %s out of %s for %s", totalMinted, amount, assetID)
	}

	return nil
}

func (w *Service) Stake(ctx context.Context, receiverName, receiverAddress, assetID string, amount *num.Uint, from string) error {
	w.log.With(logging.String("receiverAddress", receiverAddress)).Debug("Staking...")

	if err := w.account.Stake(ctx, receiverName, receiverAddress, assetID, amount, from); err != nil {
		return fmt.Errorf("failed to stake: %w", err)
	}

	w.log.With(
		logging.String("receiverName", receiverName),
		logging.String("receiverPubKey", receiverAddress),
		logging.String("targetAmount", amount.String()),
	).Debugf("%s: Waiting for staking...", from)

	if err := w.accountStream.WaitForStakeLinkingToFinalise(ctx, receiverAddress); err != nil {
		return fmt.Errorf("failed to finalise stake: %w", err)
	}

	return nil
}
