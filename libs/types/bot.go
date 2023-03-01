package types

import (
	"context"

	"code.vegaprotocol.io/shared/libs/num"
	"code.vegaprotocol.io/vega/protos/vega"
)

type PauseSignal struct {
	From  string
	Pause bool
}

type TopUpRequest struct {
	Ctx             context.Context
	ReceiverName    string
	ReceiverAddress string
	Asset           *vega.Asset
	Amount          *num.Uint
	From            string
	ErrResp         chan error
}
