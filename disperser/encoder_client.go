package disperser

import (
	"context"

	"github.com/0glabs/0g-data-avail/core"
)

type EncoderClient interface {
	EncodeBlob(ctx context.Context, data []byte, dims core.MatrixDimsions) (*core.ExtendedMatrix, error)
}
