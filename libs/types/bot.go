package types

import (
	"context"

	"code.vegaprotocol.io/shared/libs/num"
)

type PauseSignal struct {
	From  string
	Pause bool
}

type TopUpRequest struct {
	Ctx             context.Context
	ReceiverName    string
	ReceiverAddress string
	AssetID         string
	Amount          *num.Uint
	From            string
	ErrResp         chan error
}
