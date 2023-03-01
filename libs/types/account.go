package types

import (
	"context"
	"time"

	"code.vegaprotocol.io/shared/libs/cache"
	"code.vegaprotocol.io/shared/libs/num"
	"code.vegaprotocol.io/vega/protos/vega"
)

type AccountStream interface {
	WaitForTopUpToFinalise(ctx context.Context, receiverKey, assetID string, amount *num.Uint, timeout time.Duration) error
	WaitForStakeLinkingToFinalise(ctx context.Context, receiverKey string) error
	AssetByID(ctx context.Context, assetID string) (*vega.Asset, error)
	GetBalances(ctx context.Context, assetID string, pubKey string) (BalanceStore, error)
	GetStake(ctx context.Context, pubKey string) (*num.Uint, error)
}

type BalanceStore interface {
	Balance() cache.Balance
	BalanceSet(sets ...func(*cache.Balance))
}

type AssetStore interface {
	Assets() []vega.Asset
	AssetSet(sets ...func(*[]vega.Asset))
}
