package disperser

import (
	"context"

	"github.com/0glabs/0g-da-client/common"
	"github.com/0glabs/0g-da-client/core"
	pb "github.com/0glabs/0g-da-client/disperser/api/grpc/signer"
)

type SignerClient interface {
	BatchSign(ctx context.Context, addr string, data []*pb.SignRequest, log common.Logger) ([]*core.Signature, error)
}
