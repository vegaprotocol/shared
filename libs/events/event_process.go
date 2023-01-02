package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	e "code.vegaprotocol.io/shared/libs/errors"
	"code.vegaprotocol.io/shared/libs/types"
	"code.vegaprotocol.io/vega/logging"
	coreapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
)

type busEventProcessor struct {
	node    busStreamer
	log     *logging.Logger
	pauseCh chan types.PauseSignal
}

func NewBusEventProcessor(log *logging.Logger, node busStreamer, opts ...Option) *busEventProcessor {
	b := &busEventProcessor{
		node: node,
		log:  log.Named("EventProcessor"),
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

type Option func(*busEventProcessor)

func WithPauseCh(ch chan types.PauseSignal) Option {
	return func(b *busEventProcessor) {
		b.pauseCh = ch
	}
}

func (b *busEventProcessor) ProcessEvents(
	ctx context.Context,
	name string,
	req *coreapipb.ObserveEventBusRequest,
	process func(*coreapipb.ObserveEventBusResponse) (bool, error),
) <-chan error {
	errCh := make(chan error)

	var stop bool
	go func() {
		defer func() {
			close(errCh)
		}()
		for s := b.mustGetStream(ctx, name, req); !stop; {
			select {
			case <-ctx.Done():
				return
			default:
				if s == nil {
					return
				}

				rsp, err := s.Recv()
				if err != nil {
					if ctx.Err() == context.DeadlineExceeded {
						return
					}

					b.log.With(
						logging.String("name", name),
					).Warn("Stream closed, resubscribing...", logging.Error(err))

					b.pause(true, name)
					s = b.mustGetStream(ctx, name, req)
					b.pause(false, name)
					continue
				}

				stop, err = process(rsp)
				if err != nil {
					b.log.With(
						logging.String("name", name),
					).Warn("Unable to process event")
					select {
					case errCh <- err:
					default:
					}
				}
			}
		}
	}()
	return errCh
}

func (b *busEventProcessor) mustGetStream(
	ctx context.Context,
	name string,
	req *coreapipb.ObserveEventBusRequest,
) coreapipb.CoreService_ObserveEventBusClient {
	var (
		s   coreapipb.CoreService_ObserveEventBusClient
		err error
	)

	attempt := 0
	sleepTime := time.Second * 3

	for s, err = b.getStream(ctx, req); err != nil; s, err = b.getStream(ctx, req) {
		if errors.Unwrap(err).Error() == e.ErrConnectionNotReady.Error() {
			b.log.With(
				logging.String("name", name),
				logging.Int("attempt", attempt),
			).Warn("Node is not ready, reconnecting", logging.Error(err))

			b.node.MustDialConnection(ctx)

			b.log.With(
				logging.String("name", name),
				logging.Int("attempt", attempt),
			).Debug("Node reconnected, reattempting to subscribe to stream")
		} else if ctx.Err() == context.DeadlineExceeded {
			b.log.With(
				logging.String("name", name),
			).Warn("Deadline exceeded. Stopping event processor")

			break
		} else {
			attempt++

			b.log.With(
				logging.String("name", name),
				logging.Int("attempt", attempt),
				logging.Duration("sleep_time", sleepTime),
			).Error("Failed to subscribe to stream, retrying...", logging.Error(err))

			time.Sleep(sleepTime)
		}
	}

	return s
}

func (b *busEventProcessor) getStream(ctx context.Context, req *coreapipb.ObserveEventBusRequest) (coreapipb.CoreService_ObserveEventBusClient, error) {
	s, err := b.node.ObserveEventBus(ctx)
	if err != nil {
		return nil, err
	}
	// Then we subscribe to the data
	if err = s.SendMsg(req); err != nil {
		return nil, fmt.Errorf("failed to send event bus request for stream: %w", err)
	}
	return s, nil
}

func (b *busEventProcessor) pause(p bool, name string) {
	if b.pauseCh == nil {
		return
	}
	select {
	case b.pauseCh <- types.PauseSignal{From: name, Pause: p}:
	default:
	}
}
