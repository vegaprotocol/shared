package account

import (
	"context"
	"fmt"
	"math"

	"code.vegaprotocol.io/shared/libs/cache"
	"code.vegaprotocol.io/shared/libs/num"
	"code.vegaprotocol.io/shared/libs/types"
	"code.vegaprotocol.io/vega/logging"
)

type Service struct {
	name          string
	pubKey        string
	stores        map[string]balanceStore
	accountStream accountStream
	coinProvider  CoinProvider
	log           *logging.Logger
}

func NewService(log *logging.Logger, name, pubKey string, accountStream accountStream, coinProvider CoinProvider) *Service {
	return &Service{
		name:          name,
		pubKey:        pubKey,
		stores:        make(map[string]balanceStore),
		accountStream: accountStream,
		coinProvider:  coinProvider,
		log:           log.Named("AccountService"),
	}
}

func (a *Service) EnsureBalance(ctx context.Context, assetID string, balanceFn func(cache.Balance) *num.Uint, targetAmount *num.Uint, dp, scale uint64, from string) error {
	store, err := a.getStore(ctx, assetID)
	if err != nil {
		return err
	}

	// for liquidity provision and placing orders, we need only General account balance
	// for liquidity increase, we need both Bond and General account balance
	balance := balanceFn(store.Balance())

	asset, err := a.accountStream.AssetByID(ctx, assetID)
	if err != nil {
		return fmt.Errorf("failed to get asset by id: %w", err)
	}

	// if asset decimal places is higher than market decimal places, we need to scale up the amount by the difference
	if assetDP := asset.Details.Decimals; dp > 0 && assetDP > dp {
		dpDiff := assetDP - dp
		targetAmount = num.Zero().Mul(targetAmount, num.NewUint(uint64(math.Pow10(int(dpDiff)))))
	}

	if balance.GTE(targetAmount) {
		return nil
	}

	askAmount := targetAmount.Clone()

	if scale > 1 {
		askAmount = num.Zero().Mul(targetAmount, num.NewUint(scale))
	}

	a.log.With(
		logging.String("name", a.name),
		logging.String("partyId", a.pubKey),
		logging.String("asset", assetID),
		logging.String("balance", balance.String()),
		logging.String("targetAmount", targetAmount.String()),
		logging.String("askAmount", askAmount.String()),
	).Debugf("%s: Account balance is less than target amount, depositing...", from)

	if err = a.topUp(ctx, assetID, askAmount); err != nil {
		return fmt.Errorf("failed to top up: %w", err)
	}

	return nil
}

func (a *Service) topUp(ctx context.Context, assetID string, askAmount *num.Uint) error {
	errCh := make(chan error)

	a.coinProvider.TopUpChan() <- types.TopUpRequest{
		Ctx:             ctx,
		ReceiverAddress: a.pubKey,
		ReceiverName:    a.name,
		AssetID:         assetID,
		Amount:          askAmount,
		ErrResp:         errCh,
	}

	if err := <-errCh; err != nil {
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

	a.log.With(
		logging.String("name", a.name),
		logging.String("receiverName", receiverName),
		logging.String("receiverPubKey", receiverPubKey),
		logging.String("partyId", a.pubKey),
		logging.String("stake", stake.String()),
		logging.String("targetAmount", targetAmount.String()),
	).Debugf("%s: Account Stake balance is less than target amount, staking...", from)

	if err = a.coinProvider.Stake(ctx, receiverName, receiverPubKey, assetID, targetAmount, from); err != nil {
		return fmt.Errorf("failed to stake: %w", err)
	}

	return nil
}

func (a *Service) Stake(ctx context.Context, receiverName, receiverPubKey, assetID string, amount *num.Uint, from string) error {
	return a.coinProvider.Stake(ctx, receiverName, receiverPubKey, assetID, amount, from)
}

func (a *Service) Balance(ctx context.Context, assetID string) cache.Balance {
	store, err := a.getStore(ctx, assetID)
	if err != nil {
		a.log.Error("failed to get balance store", logging.Error(err))
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
