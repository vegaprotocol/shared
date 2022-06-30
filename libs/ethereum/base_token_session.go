package ethereum

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ClientBaseTokenSession struct {
	BaseTokenSession
	syncTimeout time.Duration
	address     common.Address
}

func (ts ClientBaseTokenSession) Address() common.Address {
	return ts.address
}

func (ts ClientBaseTokenSession) ApproveSync(spender common.Address, value *big.Int) (*types.Transaction, error) {
	sink := make(chan *BaseTokenApproval)

	sub, err := ts.Contract.WatchApproval(&bind.WatchOpts{}, sink, []common.Address{ts.CallOpts.From}, []common.Address{spender})
	if err != nil {
		return nil, fmt.Errorf("failed to watch for approval: %w", err)
	}
	defer sub.Unsubscribe()

	tx, err := ts.Approve(spender, value)
	if err != nil {
		return nil, err
	}

	return wait(sink, sub, tx, ts.syncTimeout)
}

func (ts ClientBaseTokenSession) TransferSync(recipient common.Address, value *big.Int) (*types.Transaction, error) {
	sink := make(chan *BaseTokenTransfer)

	sub, err := ts.Contract.WatchTransfer(&bind.WatchOpts{}, sink, []common.Address{ts.CallOpts.From}, []common.Address{recipient})
	if err != nil {
		return nil, fmt.Errorf("failed to watch for transfer: %w", err)
	}
	defer sub.Unsubscribe()

	tx, err := ts.Transfer(recipient, value)
	if err != nil {
		return nil, err
	}

	return wait(sink, sub, tx, ts.syncTimeout)
}

func (ts ClientBaseTokenSession) TransferFromSync(sender common.Address, recipient common.Address, value *big.Int) (*types.Transaction, error) {
	sink := make(chan *BaseTokenTransfer)

	sub, err := ts.Contract.WatchTransfer(&bind.WatchOpts{}, sink, []common.Address{sender}, []common.Address{recipient})
	if err != nil {
		return nil, fmt.Errorf("failed to watch for transfer: %w", err)
	}
	defer sub.Unsubscribe()

	tx, err := ts.TransferFrom(sender, recipient, value)
	if err != nil {
		return nil, err
	}

	return wait(sink, sub, tx, ts.syncTimeout)
}

func (ts ClientBaseTokenSession) MintSync(to common.Address, amount *big.Int) (*types.Transaction, error) {
	origBalance, err := ts.BalanceOf(to)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance for %s: %s", to, err)
	}

	tx, err := ts.Mint(to, amount)
	if err != nil {
		return nil, err
	}

	timeout := time.After(ts.syncTimeout)

	for {
		balance, err := ts.BalanceOf(to)
		if err != nil {
			return nil, fmt.Errorf("failed to get balance for %s: %s", to, err)
		}

		if balance.Cmp(origBalance) == 1 {
			return tx, nil
		}

		select {
		case <-timeout:
			return nil, fmt.Errorf("minting of token %s to %s has timed out", ts.Address(), to)
		case <-time.After(time.Second * 1):
		}
	}
}
