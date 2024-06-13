package signer

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/0glabs/0g-data-avail/common"
	"github.com/0glabs/0g-data-avail/core"
	"github.com/0glabs/0g-data-avail/disperser"
	pb "github.com/0glabs/0g-data-avail/disperser/api/grpc/signer"
	bn "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
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

func (c client) BatchSign(ctx context.Context, addr string, data []*pb.SignRequest, log common.Logger) ([]*core.Signature, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	conn, err := grpc.DialContext(
		ctxWithTimeout,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)), // 1 GiB
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial signer: %w", err)
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
	signatures := make([]*core.Signature, len(data))
	for i := 0; i < len(data); i++ {
		signature := sigBytes[i]
		signature, err := toBigEndian(signature)
		if err != nil {
			return nil, err
		}
		point, err := new(core.Signature).Deserialize(signature)
		if err != nil {
			return nil, err
		}

		signatures[i] = &core.Signature{G1Point: point}
	}

	return signatures, nil
}

func toBigEndian(b []byte) ([]byte, error) {
	if len(b) != bn.SizeOfG1AffineUncompressed {
		return nil, io.ErrShortBuffer
	}

	b[bn.SizeOfG1AffineUncompressed-1] &= 63
	for i := 0; i < fp.Bytes/2; i++ {
		b[i], b[fp.Bytes-i-1] = b[fp.Bytes-i-1], b[i]
	}

	for i := fp.Bytes; i < fp.Bytes+fp.Bytes/2; i++ {
		b[i], b[len(b)-(i-fp.Bytes)-1] = b[len(b)-(i-fp.Bytes)-1], b[i]
	}

	return b, nil
}
