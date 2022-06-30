package ethereum

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type ClientStakingBridgeSession struct {
	StakingBridgeSession
	syncTimeout time.Duration
	address     common.Address
}

func (ss ClientStakingBridgeSession) Address() common.Address {
	return ss.address
}
