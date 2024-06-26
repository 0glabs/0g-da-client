package mock

import (
	"context"

	"github.com/0glabs/0g-da-client/common"
)

type NoopRatelimiter struct {
}

func (r *NoopRatelimiter) AllowRequest(ctx context.Context, retrieverID string, blobSize uint, rate common.RateParam) (bool, error) {
	return true, nil
}
