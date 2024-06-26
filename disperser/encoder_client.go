package disperser

import (
	"context"

	"github.com/0glabs/0g-da-client/common"
	"github.com/0glabs/0g-da-client/core"
)

type EncoderClient interface {
	EncodeBlob(ctx context.Context, data []byte, log common.Logger) (*core.BlobCommitments, error)
}
