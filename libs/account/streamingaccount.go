package account

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/slices"

	"code.vegaprotocol.io/shared/libs/cache"
	"code.vegaprotocol.io/shared/libs/events"
	"code.vegaprotocol.io/shared/libs/num"
	"code.vegaprotocol.io/shared/libs/types"
	"code.vegaprotocol.io/vega/logging"
	dataapipb "code.vegaprotocol.io/vega/protos/data-node/api/v2"
	"code.vegaprotocol.io/vega/protos/vega"
	coreapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
	eventspb "code.vegaprotocol.io/vega/protos/vega/events/v1"
)

type account struct {
	name          string
	log           *logging.Logger
	node          dataNode
	balanceStores map[string]*balanceStores // pubKey: balanceStore
	busEvProc     busEventer

	mu              sync.Mutex
	waitingDeposits map[string]*num.Uint
}

func NewStream(log *logging.Logger, name string, node dataNode, pauseCh chan types.PauseSignal) *account {
	return &account{
		name:            name,
		log:             log.Named("AccountStreamer"),
		node:            node,
		waitingDeposits: make(map[string]*num.Uint),
		busEvProc:       events.NewBusEventProcessor(log, node, events.WithPauseCh(pauseCh)),
		balanceStores:   make(map[string]*balanceStores),
	}
}

func (a *account) GetBalances(ctx context.Context, assetID string, pubKey string) (balanceStore, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	stores, okAcc := a.balanceStores[pubKey]
	if !okAcc {
		stores = &balanceStores{
			balanceStores: make(map[string]balanceStore),
		}
		a.balanceStores[pubKey] = stores
	} else {
		store, okAss := stores.get(assetID)
		if okAss {
			return store, nil
		}
	}

	accounts, err := a.node.PartyAccounts(ctx, &dataapipb.ListAccountsRequest{
		Filter: &dataapipb.AccountFilter{
			PartyIds: []string{pubKey},
			AssetId:  assetID,
		},
	})
	if err != nil {
		return nil, err
	}

	store := cache.NewBalanceStore()
	stores.set(assetID, store)
	a.balanceStores[pubKey] = stores

	for _, acc := range accounts {
		if err = a.setBalanceByType(acc.Type, acc.Balance, store); err != nil {
			a.log.Error("failed to set account balance",
				logging.String("accountType", acc.Type.String()),
				logging.Error(err))
		}
	}

	a.subscribeToAccountEvents(ctx, pubKey)

	return store, nil
}

func (a *account) GetStake(ctx context.Context, pubKey string) (*num.Uint, error) {
	partyStakeResp, err := a.node.PartyStake(ctx, &dataapipb.GetStakeRequest{
		PartyId: pubKey,
	})
	if err != nil {
		return nil, err
	}

	stake, overflow := num.UintFromString(partyStakeResp.CurrentStakeAvailable, 10)
	if overflow {
		return nil, fmt.Errorf("failed to convert stake to uint: %w", err)
	}

	return stake, nil
}

func (a *account) subscribeToAccountEvents(ctx context.Context, pubKey string) {
	req := &coreapipb.ObserveEventBusRequest{
		Type: []eventspb.BusEventType{
			eventspb.BusEventType_BUS_EVENT_TYPE_ACCOUNT,
		},
		PartyId: pubKey,
	}

	proc := func(rsp *coreapipb.ObserveEventBusResponse) (bool, error) {
		for _, event := range rsp.Events {
			acct := event.GetAccount()
			// filter out any that are for different assets
			store, err := a.GetBalances(ctx, acct.Asset, pubKey)
			if err != nil {
				a.log.Error("failed to init balance store", logging.Error(err))
				return true, err
			}

			if err := a.setBalanceByType(acct.Type, acct.Balance, store); err != nil {
				a.log.Error("failed to set account balance",
					logging.String("accountType", acct.Type.String()),
					logging.Error(err))
			}
		}
		return false, nil
	}

	a.busEvProc.ProcessEvents(context.Background(), "AccountData: "+a.name, req, proc)
}

func (a *account) setBalanceByType(accountType vega.AccountType, balanceStr string, store balanceStore) error {
	balance, err := num.ConvertUint256(balanceStr)
	if err != nil {
		return fmt.Errorf("failed to convert account balance: %w", err)
	}

	store.BalanceSet(cache.SetBalanceByType(accountType, balance))
	return nil
}

func (a *account) AssetByID(ctx context.Context, assetID string) (*vega.Asset, error) {
	return a.node.AssetByID(ctx, &dataapipb.GetAssetRequest{
		AssetId: assetID,
	})
}

// WaitForTopUpToFinalise is a blocking call that waits for the top-up finalise event to be received.
func (a *account) WaitForTopUpToFinalise(
	ctx context.Context,
	receiverPubKey,
	assetID string,
	expectAmount *num.Uint,
	timeout time.Duration,
) error {
	if exist, ok := a.getWaitingDeposit(assetID); ok {
		if !expectAmount.EQ(exist) {
			a.setWaitingDeposit(assetID, expectAmount)
		}
		return nil
	}

	req := &coreapipb.ObserveEventBusRequest{
		Type: []eventspb.BusEventType{
			eventspb.BusEventType_BUS_EVENT_TYPE_DEPOSIT,
			eventspb.BusEventType_BUS_EVENT_TYPE_TRANSFER,
			eventspb.BusEventType_BUS_EVENT_TYPE_ACCOUNT,
		},
		// PartyId: receiverPubKey, TODO: ??
		// AssetId: assetID, TODO: ??
	}

	proc := func(rsp *coreapipb.ObserveEventBusResponse) (bool, error) {
		for _, event := range rsp.Events {
			var (
				status   string
				reason   string
				partyId  string
				asset    string
				balance  string
				from     string
				amount   string
				okStatus []string
			)
			switch event.Type {
			case eventspb.BusEventType_BUS_EVENT_TYPE_DEPOSIT:
				depEvt := event.GetDeposit()
				status = depEvt.Status.String()
				partyId = depEvt.PartyId
				asset = depEvt.Asset
				amount = depEvt.Amount
				okStatus = []string{
					vega.Deposit_STATUS_FINALIZED.String(),
					vega.Deposit_STATUS_OPEN.String(),
				}
			case eventspb.BusEventType_BUS_EVENT_TYPE_TRANSFER:
				trfEvt := event.GetTransfer()
				status = trfEvt.Status.String()
				if trfEvt.Reason != nil {
					reason = *trfEvt.Reason
				}
				from = trfEvt.From
				partyId = trfEvt.To
				asset = trfEvt.Asset
				amount = trfEvt.Amount
				okStatus = []string{
					eventspb.Transfer_STATUS_DONE.String(),
					eventspb.Transfer_STATUS_PENDING.String(),
				}
			case eventspb.BusEventType_BUS_EVENT_TYPE_ACCOUNT:
				accEvt := event.GetAccount()
				partyId = accEvt.Owner
				asset = accEvt.Asset
				// we only care about the balance of the general account
				if accEvt.Type == vega.AccountType_ACCOUNT_TYPE_GENERAL {
					balance = accEvt.GetBalance()
				}
			}

			// filter out any that are for different assets, or not finalized
			if partyId != receiverPubKey || asset != assetID {
				continue
			}

			// if it's a deposit or transfer event, check if it failed
			if status != "" && !slices.Contains(okStatus, status) {
				genBalance := ""
				nodeBalance := ""

				if from != "" {
					balances, ok := a.balanceStores[from].get(assetID)
					if ok {
						genBalance = cache.General(balances.Balance()).String()
					}
				}
				a.log.With(
					logging.String("status", status),
					logging.String("reason", reason),
					logging.String("partyId", partyId),
					logging.String("from", from),
					logging.String("asset", asset),
					logging.String("balance", balance),
					logging.String("amount", amount),
					logging.String("expectAmount", expectAmount.String()),
					logging.String("balance.general", genBalance),
					logging.String("balance.node", nodeBalance),
				).Error("deposit or transfer failed")
				return true, fmt.Errorf("transfer %s failed: %s: %s", event.Id, status, reason)
			}

			// only check the deposited amount for account events
			if event.Type != eventspb.BusEventType_BUS_EVENT_TYPE_ACCOUNT {
				continue
			}

			// TODO: either check 0 or status
			if balance == "" || balance == "0" {
				continue
			}

			gotAmount, err := num.ConvertUint256(balance)
			if err != nil {
				return false, fmt.Errorf("failed to parse top-up expectAmount %s: %w", expectAmount.String(), err)
			}

			expect, ok := a.getWaitingDeposit(assetID)
			if !ok {
				expect = expectAmount.Clone()
				a.setWaitingDeposit(assetID, expect)
			}

			if gotAmount.GTE(expect) {
				a.log.With(
					logging.String("name", a.name),
					logging.String("partyId", receiverPubKey),
					logging.String("balance", gotAmount.String()),
				).Info("TopUp finalised")
				a.deleteWaitingDeposit(assetID)
				if _, err = a.GetBalances(ctx, assetID, receiverPubKey); err != nil {
					a.log.Error("failed to set balance after top-up", logging.Error(err))
				}
				return true, nil
			} else if !gotAmount.IsZero() {
				a.log.With(
					logging.String("name", a.name),
					logging.String("partyId", receiverPubKey),
					logging.String("gotAmount", gotAmount.String()),
					logging.String("targetAmount", expect.String()),
				).Info("Received funds, but balance is less than expected")
				// if we received fewer funds than expected, keep waiting (e.g. faucet tops-up in multiple iterations)
			}
		}
		return false, nil
	}

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	errCh := a.busEvProc.ProcessEvents(ctx, "TopUpData: "+a.name, req, proc)
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return fmt.Errorf("timed out waiting for top-up event")
	}
}

func (a *account) getWaitingDeposit(assetID string) (*num.Uint, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	req, ok := a.waitingDeposits[assetID]
	if ok {
		return req.Clone(), ok
	}
	return nil, false
}

func (a *account) setWaitingDeposit(assetID string, amount *num.Uint) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.waitingDeposits[assetID] = amount.Clone()
}

func (a *account) deleteWaitingDeposit(assetID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.waitingDeposits, assetID)
}

func (a *account) WaitForStakeLinkingToFinalise(ctx context.Context, pubKey string) error {
	req := &coreapipb.ObserveEventBusRequest{
		Type: []eventspb.BusEventType{eventspb.BusEventType_BUS_EVENT_TYPE_STAKE_LINKING},
	}

	proc := func(rsp *coreapipb.ObserveEventBusResponse) (bool, error) {
		for _, event := range rsp.GetEvents() {
			stake := event.GetStakeLinking()
			if stake.Party != pubKey {
				continue
			}

			if stake.Status != eventspb.StakeLinking_STATUS_ACCEPTED {
				if stake.Status == eventspb.StakeLinking_STATUS_PENDING {
					continue
				} else {
					return true, fmt.Errorf("stake linking failed: %s", stake.Status.String())
				}
			}
			a.log.With(
				logging.String("name", a.name),
				logging.String("partyId", stake.Party),
				logging.String("stakeID", stake.Id),
			).Info("Received stake linking")
			return true, nil
		}
		return false, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*450)
	defer cancel()

	errCh := a.busEvProc.ProcessEvents(ctx, "StakeLinking: "+a.name, req, proc)
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return fmt.Errorf("timed out waiting for top-up event")
	}
}

type balanceStores struct {
	mu            sync.Mutex
	balanceStores map[string]balanceStore
}

func (b *balanceStores) get(assetID string) (balanceStore, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	store, ok := b.balanceStores[assetID]
	return store, ok
}

func (b *balanceStores) set(assetID string, store balanceStore) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.balanceStores[assetID] = store
}
