package types

import (
	"context"
	"time"

	"code.vegaprotocol.io/shared/libs/num"
)

type AccountStream interface {
	WaitForTopUpToFinalise(ctx context.Context, receiverKey, assetID string, amount *num.Uint, timeout time.Duration) error
	WaitForStakeLinkingToFinalise(ctx context.Context, receiverKey string) error
}
