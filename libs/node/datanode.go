package node

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	e "code.vegaprotocol.io/shared/libs/errors"
	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
)

// DataNode stores state for a Vega Data node.
type DataNode struct {
	hosts       []string // format: host:port
	callTimeout time.Duration
	conn        *grpc.ClientConn
	mu          sync.RWMutex
	wg          sync.WaitGroup
	once        sync.Once
}

// NewDataNode returns a new node.
func NewDataNode(hosts []string, callTimeoutMil int) *DataNode {
	return &DataNode{
		hosts:       hosts,
		callTimeout: time.Duration(callTimeoutMil) * time.Millisecond,
	}
}

// MustDialConnection tries to establish a connection to one of the nodes from a list of locations.
// It is idempotent, where each call will block the caller until a connection is established.
func (n *DataNode) MustDialConnection(ctx context.Context) {
	n.once.Do(func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		n.wg.Add(len(n.hosts))

		for _, h := range n.hosts {
			go func(host string) {
				defer func() {
					cancel()
					n.wg.Done()
				}()
				n.dialNode(ctx, host)
			}(h)
		}
		n.wg.Wait()
		n.mu.Lock()
		defer n.mu.Unlock()

		if n.conn == nil {
			log.Fatalf("Failed to connect to DataNode")
		}
	})

	n.wg.Wait()
	n.once = sync.Once{}
}

func (n *DataNode) dialNode(ctx context.Context, host string) {
	conn, err := grpc.DialContext(
		ctx,
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		if err != context.Canceled {
			log.Printf("Failed to dial node '%s': %s\n", host, err)
		}
		return
	}

	n.mu.Lock()
	n.conn = conn
	n.mu.Unlock()
}

func (n *DataNode) Target() string {
	return n.conn.Target()
}

// === CoreService ===

// SubmitTransaction submits a signed v2 transaction.
func (n *DataNode) SubmitTransaction(ctx context.Context, req *vegaapipb.SubmitTransactionRequest) (*vegaapipb.SubmitTransactionResponse, error) {
	msg := "gRPC call failed: SubmitTransaction: %w"
	if n == nil {
		return nil, fmt.Errorf(msg, e.ErrNil)
	}

	if n.conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(msg, e.ErrConnectionNotReady)
	}

	c := vegaapipb.NewCoreServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(ctx, n.callTimeout)
	defer cancel()

	response, err := c.SubmitTransaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf(msg, e.ErrorDetail(err))
	}

	return response, nil
}

// LastBlockData gets the latest blockchain data, height, hash and pow parameters.
func (n *DataNode) LastBlockData(ctx context.Context) (*vegaapipb.LastBlockHeightResponse, error) {
	msg := "gRPC call failed: LastBlockData: %w"
	if n == nil {
		return nil, fmt.Errorf(msg, e.ErrNil)
	}

	if n.conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(msg, e.ErrConnectionNotReady)
	}

	c := vegaapipb.NewCoreServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(ctx, n.callTimeout)
	defer cancel()

	var response *vegaapipb.LastBlockHeightResponse

	response, err := c.LastBlockHeight(ctx, &vegaapipb.LastBlockHeightRequest{})
	if err != nil {
		err = fmt.Errorf(msg, e.ErrorDetail(err))
	}

	return response, err
}

// ObserveEventBus opens a stream.
func (n *DataNode) ObserveEventBus(ctx context.Context) (vegaapipb.CoreService_ObserveEventBusClient, error) {
	msg := "gRPC call failed: ObserveEventBus: %w"
	if n == nil {
		return nil, fmt.Errorf(msg, e.ErrNil)
	}

	if n.conn == nil || n.conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(msg, e.ErrConnectionNotReady)
	}

	c := vegaapipb.NewCoreServiceClient(n.conn)
	// no timeout on streams
	client, err := c.ObserveEventBus(ctx)
	if err != nil {
		return nil, fmt.Errorf(msg, e.ErrorDetail(err))
	}

	return client, nil
}

// === TradingDataService ===

// PartyAccounts returns accounts for the given party.
func (n *DataNode) PartyAccounts(ctx context.Context, req *dataapipb.ListAccountsRequest) ([]*dataapipb.AccountBalance, error) {
	msg := "gRPC call failed (data-node): PartyAccounts: %w"
	if n == nil {
		return nil, fmt.Errorf(msg, e.ErrNil)
	}

	if n.conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(msg, e.ErrConnectionNotReady)
	}

	c := dataapipb.NewTradingDataServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(ctx, n.callTimeout)
	defer cancel()

	response, err := c.ListAccounts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf(msg, e.ErrorDetail(err))
	}

	var accounts []*dataapipb.AccountBalance
	for _, a := range response.Accounts.Edges {
		accounts = append(accounts, a.Node)
	}

	return accounts, nil
}

func (n *DataNode) PartyStake(ctx context.Context, req *dataapipb.GetStakeRequest) (*dataapipb.GetStakeResponse, error) {
	msg := "gRPC call failed (data-node): PartyStake: %w"
	if n == nil {
		return nil, fmt.Errorf(msg, e.ErrNil)
	}

	if n.conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(msg, e.ErrConnectionNotReady)
	}

	c := dataapipb.NewTradingDataServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(ctx, n.callTimeout)
	defer cancel()

	response, err := c.GetStake(ctx, req)
	if err != nil {
		return nil, fmt.Errorf(msg, e.ErrorDetail(err))
	}

	return response, nil
}

// MarketDataByID returns market data for the specified market.
func (n *DataNode) MarketDataByID(ctx context.Context, req *dataapipb.GetLatestMarketDataRequest) (*vega.MarketData, error) {
	msg := "gRPC call failed (data-node): MarketDataByID: %w"
	if n == nil {
		return nil, fmt.Errorf(msg, e.ErrNil)
	}

	if n.conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(msg, e.ErrConnectionNotReady)
	}

	c := dataapipb.NewTradingDataServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(ctx, n.callTimeout)
	defer cancel()

	response, err := c.GetLatestMarketData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf(msg, e.ErrorDetail(err))
	}

	return response.GetMarketData(), nil
}

// Markets returns all markets.
func (n *DataNode) Markets(ctx context.Context, req *dataapipb.ListMarketsRequest) ([]*vega.Market, error) {
	msg := "gRPC call failed (data-node): Markets: %w"
	if n == nil {
		return nil, fmt.Errorf(msg, e.ErrNil)
	}

	if n.conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(msg, e.ErrConnectionNotReady)
	}

	c := dataapipb.NewTradingDataServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(ctx, n.callTimeout)
	defer cancel()

	response, err := c.ListMarkets(ctx, req)
	if err != nil {
		return nil, fmt.Errorf(msg, e.ErrorDetail(err))
	}

	var markets []*vega.Market
	for _, m := range response.Markets.Edges {
		markets = append(markets, m.Node)
	}

	return markets, nil
}

// PositionsByParty returns positions for the given party.
func (n *DataNode) PositionsByParty(ctx context.Context, req *dataapipb.ListPositionsRequest) ([]*vega.Position, error) {
	msg := "gRPC call failed (data-node): PositionsByParty: %w"
	if n == nil {
		return nil, fmt.Errorf(msg, e.ErrNil)
	}

	if n.conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(msg, e.ErrConnectionNotReady)
	}

	c := dataapipb.NewTradingDataServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(ctx, n.callTimeout)
	defer cancel()

	response, err := c.ListPositions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf(msg, e.ErrorDetail(err))
	}

	var positions []*vega.Position
	for _, p := range response.Positions.Edges {
		positions = append(positions, p.Node)
	}

	return positions, nil
}

// AssetByID returns the specified asset.
func (n *DataNode) AssetByID(ctx context.Context, req *dataapipb.GetAssetRequest) (*vega.Asset, error) {
	msg := "gRPC call failed (data-node): AssetByID: %w"
	if n == nil {
		return nil, fmt.Errorf(msg, e.ErrNil)
	}

	if n.conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf(msg, e.ErrConnectionNotReady)
	}

	c := dataapipb.NewTradingDataServiceClient(n.conn)
	ctx, cancel := context.WithTimeout(ctx, n.callTimeout)
	defer cancel()

	response, err := c.GetAsset(ctx, req)
	if err != nil {
		return nil, fmt.Errorf(msg, e.ErrorDetail(err))
	}

	return response.GetAsset(), nil
}

func (n *DataNode) WaitForStateChange(ctx context.Context, state connectivity.State) bool {
	return n.conn.WaitForStateChange(ctx, state)
}
