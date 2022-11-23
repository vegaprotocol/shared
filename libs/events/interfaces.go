package events

import (
	"context"

	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
)

type busStreamer interface {
	MustDialConnection(ctx context.Context)
	ObserveEventBus(ctx context.Context) (client vegaapipb.CoreService_ObserveEventBusClient, err error)
}
