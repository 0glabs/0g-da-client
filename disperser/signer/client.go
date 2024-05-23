package signer

import (
	"context"
	"fmt"
	"time"

	"github.com/0glabs/0g-data-avail/core"
	"github.com/0glabs/0g-data-avail/disperser"
	pb "github.com/0glabs/0g-data-avail/disperser/api/grpc/signer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	timeout time.Duration
}

func NewSignerClient(timeout time.Duration) (disperser.SignerClient, error) {
	return client{
		timeout: timeout,
	}, nil
}

func (c client) BatchSign(ctx context.Context, addr string, data []*pb.SignRequest) ([]*core.Signature, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)), // 1 GiB
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial encoder: %w", err)
	}
	defer conn.Close()

	signer := pb.NewSignerClient(conn)
	// requests := make([]*pb.SignRequest, 0, len(data))
	// for i, req := range data {
	// 	requests[i] = &pb.SignRequest{
	// 		Epoch:             req.Epoch,
	// 		ErasureCommitment: req.ErasureCommitment,
	// 		StorageRoot:       req.StorageRoot,
	// 		EncodedSlice:      req.EncodedSlice,
	// 	}
	// }

	reply, err := signer.BatchSign(ctx, &pb.BatchSignRequest{
		Requests: data,
	})
	if err != nil {
		return nil, err
	}

	sigBytes := reply.GetSignatures()
	signatures := make([]*core.Signature, 0, len(data))
	for i := 0; i < len(data); i++ {
		point, err := new(core.Signature).Deserialize(sigBytes[i])
		if err != nil {
			return nil, err
		}

		signatures[i] = &core.Signature{G1Point: point}
	}

	return signatures, nil
}
