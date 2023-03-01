package whale

import (
	"context"

	"code.vegaprotocol.io/shared/libs/cache"
	"code.vegaprotocol.io/shared/libs/num"
	"code.vegaprotocol.io/vega/protos/vega"
)

type erc20Service interface {
	Stake(ctx context.Context, ownerPrivateKey, ownerAddress, vegaTokenAddress, vegaPubKey string, amount *num.Uint) (*num.Uint, error)
	Deposit(ctx context.Context, ownerPrivateKey, ownerAddress, tokenAddress, vegaPubKey string, amount *num.Uint) (*num.Uint, error)
}

type faucetClient interface {
	Mint(ctx context.Context, amount string, asset, party string) (bool, error)
}

type accountService interface {
	EnsureBalance(ctx context.Context, asset *vega.Asset, balanceFn func(cache.Balance) *num.Uint, targetAmount *num.Uint, dp, scale uint64, from string) error
	EnsureStake(ctx context.Context, receiverName, receiverPubKey string, asset *vega.Asset, targetAmount *num.Uint, from string) error
	Stake(ctx context.Context, receiverName, receiverPubKey string, asset *vega.Asset, amount *num.Uint, from string) error
}
