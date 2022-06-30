package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
)

var defaultSyncDuration = time.Second * 5

type Client struct {
	*ethclient.Client
	chainID *big.Int
}

func NewClient(ctx context.Context, ethereumAddress string, chainID int64) (*Client, error) {
	addr, err := url.Parse(ethereumAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Ethereum address: %w", err)
	}

	addr.Scheme = "ws"

	client, err := ethclient.DialContext(ctx, addr.String())
	if err != nil {
		return nil, fmt.Errorf("failed to dial Ethereum client: %s", err)
	}

	return &Client{
		chainID: big.NewInt(chainID),
		Client:  client,
	}, nil
}

func (ec *Client) NewERC20BridgeSession(
	ctx context.Context,
	contractOwnerPrivateKey string,
	bridgeAddress common.Address,
	syncTimeout *time.Duration,
) (*ClientERC20BridgeSession, error) {
	privateKey, err := crypto.HexToECDSA(contractOwnerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert erc20 bridge contract owner private key hash into ECDSA: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, ec.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create erc20 bridge contract authentication: %w", err)
	}

	bridge, err := NewERC20Bridge(bridgeAddress, ec.Client)
	if err != nil {
		return nil, fmt.Errorf("failed creating erc20 bridge contract for address %q: %w", bridgeAddress, err)
	}

	if syncTimeout == nil {
		syncTimeout = &defaultSyncDuration
	}

	return &ClientERC20BridgeSession{
		ERC20BridgeSession: ERC20BridgeSession{
			Contract: bridge,
			CallOpts: bind.CallOpts{
				From:    auth.From,
				Context: ctx,
			},
			TransactOpts: *auth,
		},
		syncTimeout: *syncTimeout,
		address:     bridgeAddress,
	}, nil
}

func (ec *Client) NewStakingBridgeSession(
	ctx context.Context,
	contractOwnerPrivateKey string,
	bridgeAddress common.Address,
	syncTimeout *time.Duration,
) (*ClientStakingBridgeSession, error) {
	privateKey, err := crypto.HexToECDSA(contractOwnerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert staking bridge contract owner private key hash into ECDSA: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, ec.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create staking bridge contract authentication: %w", err)
	}

	bridge, err := NewStakingBridge(bridgeAddress, ec.Client)
	if err != nil {
		return nil, fmt.Errorf("failed creating staking bridge contract for address %q: %w", bridgeAddress, err)
	}

	if syncTimeout == nil {
		syncTimeout = &defaultSyncDuration
	}

	return &ClientStakingBridgeSession{
		StakingBridgeSession: StakingBridgeSession{
			Contract: bridge,
			CallOpts: bind.CallOpts{
				From:    auth.From,
				Context: ctx,
			},
			TransactOpts: *auth,
		},
		syncTimeout: *syncTimeout,
		address:     bridgeAddress,
	}, nil
}

func (ec *Client) NewBaseTokenSession(
	ctx context.Context,
	contractOwnerPrivateKey string,
	tokenAddress common.Address,
	syncTimeout *time.Duration,
) (*ClientBaseTokenSession, error) {
	privateKey, err := crypto.HexToECDSA(contractOwnerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert base token contract owner private key hash into ECDSA: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, ec.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create base token contract authentication: %w", err)
	}

	token, err := NewBaseToken(tokenAddress, ec.Client)
	if err != nil {
		return nil, fmt.Errorf("failed creating base token contract for address %q: %w", tokenAddress, err)
	}

	if syncTimeout == nil {
		syncTimeout = &defaultSyncDuration
	}

	return &ClientBaseTokenSession{
		BaseTokenSession: BaseTokenSession{
			Contract: token,
			CallOpts: bind.CallOpts{
				From:    auth.From,
				Context: ctx,
			},
			TransactOpts: *auth,
		},
		syncTimeout: *syncTimeout,
		address:     tokenAddress,
	}, nil
}

func wait[T any](sink chan T, sub event.Subscription, tx *types.Transaction, timeout time.Duration) (*types.Transaction, error) {
	select {
	case <-sink:
		return tx, nil
	case err := <-sub.Err():
		return nil, err
	case <-time.After(timeout):
		return nil, fmt.Errorf("transaction time has timed out")
	}
}

func StringToByte32Array(str string) [32]byte {
	value := [32]byte{}
	copy(value[:], []byte(str))

	return value
}
