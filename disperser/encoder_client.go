package disperser

import (
	"context"

	"github.com/zero-gravity-labs/zgda/core"
)

type EncoderClient interface {
	EncodeBlob(ctx context.Context, data []byte, encodingParams core.EncodingParams) (*core.BlobCommitments, []*core.Chunk, error)
}
