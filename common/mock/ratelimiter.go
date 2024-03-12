package mock

import (
	"context"

	"github.com/0glabs/0g-data-avail/common"
)

type NoopRatelimiter struct {
}

func (r *NoopRatelimiter) AllowRequest(ctx context.Context, retrieverID string, blobSize uint, rate common.RateParam) (bool, error) {
	return true, nil
}
