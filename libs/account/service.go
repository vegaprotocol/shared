package account

import (
	"context"
	"fmt"
	"math"

	log "github.com/sirupsen/logrus"

	"code.vegaprotocol.io/shared/libs/cache"
	"code.vegaprotocol.io/shared/libs/num"
	"code.vegaprotocol.io/shared/libs/types"
)

type Service struct {
	name          string
	pubKey        string
	assetID       string
	stores        map[string]balanceStore
	accountStream accountStream
	coinProvider  CoinProvider
	log           *log.Entry
}

func NewService(name, pubKey, assetID string, accountStream accountStream, coinProvider CoinProvider) *Service {
	return &Service{
		name:          name,
		pubKey:        pubKey,
		assetID:       assetID,
		stores:        make(map[string]balanceStore),
		accountStream: accountStream,
		coinProvider:  coinProvider,
		log:           log.WithField("component", "AccountService"),
	}
}

func (a *Service) EnsureBalance(ctx context.Context, assetID string, balanceFn func(cache.Balance) *num.Uint, targetAmount *num.Uint, dp, scale uint64, from string) error {
	store, err := a.getStore(ctx, assetID)
	if err != nil {
		return err
	}

	// or liquidity provision and placing orders, we need only General account balance
	// for liquidity increase, we need both Bond and General account balance
	balance := balanceFn(store.Balance())

	if balance.GTE(targetAmount) {
		return nil
	}

	askAmount := num.Zero().Mul(targetAmount, num.NewUint(scale))
	asset, err := a.accountStream.AssetByID(ctx, assetID)
	if err != nil {
		return fmt.Errorf("failed to get asset by id: %w", err)
	}

	if assetDP := asset.Details.Decimals; dp > 0 && assetDP > dp {
		dpDiff := assetDP - dp
		askAmount = askAmount.Div(askAmount, num.NewUint(uint64(math.Pow10(int(dpDiff)))))
	}

	a.log.WithFields(
		log.Fields{
			"name":         a.name,
			"partyId":      a.pubKey,
			"asset":        assetID,
			"balance":      balance.String(),
			"targetAmount": targetAmount.String(),
			"askAmount":    askAmount.String(),
		}).Debugf("%s: Account balance is less than target amount, depositing...", from)

	errCh := make(chan error)

	a.coinProvider.TopUpChan() <- types.TopUpRequest{
		Ctx:             ctx,
		ReceiverAddress: a.pubKey,
		ReceiverName:    a.name,
		AssetID:         assetID,
		Amount:          askAmount,
		ErrResp:         errCh,
	}

	if err = <-errCh; err != nil {
		return fmt.Errorf("failed to deposit: %w", err)
	}

	return nil
}

func (a *Service) EnsureStake(ctx context.Context, receiverName, receiverPubKey, assetID string, targetAmount *num.Uint, from string) error {
	if receiverPubKey == "" {
		return fmt.Errorf("receiver public key is empty")
	}

	stake, err := a.accountStream.GetStake(ctx, receiverPubKey)
	if err != nil {
		return err
	}

	if stake.GT(targetAmount) {
		return nil
	}

	a.log.WithFields(
		log.Fields{
			"name":           a.name,
			"receiverName":   receiverName,
			"receiverPubKey": receiverPubKey,
			"partyId":        a.pubKey,
			"stake":          stake.String(),
			"targetAmount":   targetAmount.String(),
		}).Debugf("%s: Account Stake balance is less than target amount, staking...", from)

	if err = a.coinProvider.Stake(ctx, receiverName, receiverPubKey, assetID, targetAmount, from); err != nil {
		return fmt.Errorf("failed to stake: %w", err)
	}

	return nil
}

func (a *Service) Stake(ctx context.Context, receiverName, receiverPubKey, assetID string, amount *num.Uint, from string) error {
	return a.coinProvider.Stake(ctx, receiverName, receiverPubKey, assetID, amount, from)
}

func (a *Service) Balance(ctx context.Context) cache.Balance {
	store, err := a.getStore(ctx, a.assetID)
	if err != nil {
		a.log.WithError(err).Error("failed to get balance store")
		return cache.Balance{}
	}
	return store.Balance()
}

func (a *Service) getStore(ctx context.Context, assetID string) (balanceStore, error) {
	var err error
	store, ok := a.stores[assetID]
	if !ok {
		store, err = a.accountStream.GetBalances(ctx, assetID, a.pubKey)
		if err != nil {
			return nil, fmt.Errorf("failed to initialise balances for '%s': %w", assetID, err)
		}

		a.stores[assetID] = store
	}

	return store, nil
}
