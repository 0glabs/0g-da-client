package encoder

import (
	"context"
	"fmt"
	"time"

	"github.com/0glabs/0g-data-avail/core"
	"github.com/0glabs/0g-data-avail/disperser"
	pb "github.com/0glabs/0g-data-avail/disperser/api/grpc/encoder"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	addr    string
	timeout time.Duration
}

func NewEncoderClient(addr string, timeout time.Duration) (disperser.EncoderClient, error) {
	return client{
		addr:    addr,
		timeout: timeout,
	}, nil
}

func ExtendedMatrixFromReply(reply *pb.EncodeBlobReply, blobLength uint) (*core.ExtendedMatrix, error) {
	if len(reply.Chunks) != int(core.CoeffSize*reply.Cols*reply.Rows) {
		return nil, fmt.Errorf("encoded matrix data length mismatch with rows x cols")
	}
	if len(reply.Commitment) != int(core.CommitmentSize*reply.Rows) {
		return nil, fmt.Errorf("commitment length mismatch with rows")
	}
	chunksIndex := 0
	commitmentsIndex := 0
	commitments := make([]core.Commitment, 0)
	rows := make([]core.EncodedRow, 0)
	for i := 0; i < int(reply.Rows); i++ {
		row := make([]core.Coeff, 0)
		for j := 0; j < int(reply.Cols); j++ {
			var coeff core.Coeff
			copy(coeff[:], reply.Chunks[chunksIndex:chunksIndex+core.CoeffSize])
			row = append(row, coeff)
			chunksIndex += core.CoeffSize
		}
		rows = append(rows, row)

		var commitment core.Commitment
		copy(commitment[:], reply.Commitment[commitmentsIndex:commitmentsIndex+core.CommitmentSize])
		commitments = append(commitments, commitment)
		commitmentsIndex += core.CommitmentSize
	}
	return &core.ExtendedMatrix{
		Length:      blobLength,
		Rows:        rows,
		Commitments: commitments,
	}, nil
}

func (c client) EncodeBlob(ctx context.Context, data []byte, dims core.MatrixDimsions) (*core.ExtendedMatrix, error) {
	conn, err := grpc.Dial(
		c.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)), // 1 GiB
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial encoder: %w", err)
	}
	defer conn.Close()

	encoder := pb.NewEncoderClient(conn)
	reply, err := encoder.EncodeBlob(ctx, &pb.EncodeBlobRequest{
		Data: data,
		Cols: uint32(dims.Cols),
	})
	if err != nil {
		return nil, err
	}
	return ExtendedMatrixFromReply(reply, core.GetBlobLength(uint(len(data))))
}
