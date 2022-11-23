package account

import (
	"context"
	"time"

	"code.vegaprotocol.io/shared/libs/cache"
	"code.vegaprotocol.io/shared/libs/num"
	"code.vegaprotocol.io/shared/libs/types"
	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
)

type dataNode interface {
	busStreamer
	PartyAccounts(ctx context.Context, req *dataapipb.ListAccountsRequest) ([]*dataapipb.AccountBalance, error)
	PartyStake(ctx context.Context, req *dataapipb.GetStakeRequest) (response *dataapipb.GetStakeResponse, err error)
}

type busStreamer interface {
	MustDialConnection(ctx context.Context)
	ObserveEventBus(ctx context.Context) (client vegaapipb.CoreService_ObserveEventBusClient, err error)
}

type CoinProvider interface {
	TopUpAsync(ctx context.Context, receiverName, receiverAddress, assetID string, amount *num.Uint) error
	StakeAsync(ctx context.Context, receiverAddress, assetID string, amount *num.Uint) error
}

type accountStream interface {
	init(pubKey string, pauseCh chan types.PauseSignal)
	GetBalances(ctx context.Context, assetID string) (balanceStore, error)
	GetStake(ctx context.Context) (*num.Uint, error)
	WaitForStakeLinking(ctx context.Context, pubKey string) error
	WaitForTopUpToFinalise(ctx context.Context, walletPubKey, assetID string, amount *num.Uint, timeout time.Duration) error
}

type balanceStore interface {
	Balance() cache.Balance
	BalanceSet(sets ...func(*cache.Balance))
}

type busEventer interface {
	ProcessEvents(ctx context.Context, name string, req *vegaapipb.ObserveEventBusRequest, process func(*vegaapipb.ObserveEventBusResponse) (bool, error)) <-chan error
}
