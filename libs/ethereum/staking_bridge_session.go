package ethereum

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"code.vegaprotocol.io/shared/libs/ethereum/generated"
)

type StakingBridgeSession struct {
	generated.StakingBridgeSession
	syncTimeout time.Duration
	address     common.Address
}

func (ss StakingBridgeSession) Address() common.Address {
	return ss.address
}
