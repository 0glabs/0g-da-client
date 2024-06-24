package disperser

import (
	"context"

	"github.com/0glabs/0g-data-avail/common"
	"github.com/0glabs/0g-data-avail/core"
)

type EncoderClient interface {
	EncodeBlob(ctx context.Context, data []byte, log common.Logger) (*core.BlobCommitments, error)
}
