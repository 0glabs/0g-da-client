package signer

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/0glabs/0g-da-client/common"
	"github.com/0glabs/0g-da-client/core"
	"github.com/0glabs/0g-da-client/disperser"
	pb "github.com/0glabs/0g-da-client/disperser/api/grpc/signer"
	bn "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const ipv4WithPortPattern = `\b(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(?::\d{1,5})\b`
const ipv4Pattern = `\b(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`
const portPattern = `\b(\d{1,5})\b`

type client struct {
	timeout   time.Duration
	ipv4Regex *regexp.Regexp
}

func NewSignerClient(timeout time.Duration) (disperser.SignerClient, error) {
	regex := regexp.MustCompile(ipv4WithPortPattern)

	return client{
		timeout:   timeout,
		ipv4Regex: regex,
	}, nil
}

func (c client) BatchSign(ctx context.Context, addr string, data []*pb.SignRequest, log common.Logger) ([]*core.Signature, error) {
	matches := c.ipv4Regex.FindAllString(addr, -1)
	if len(matches) != 1 {
		formattedAddr := ""
		prefix := "http://"
		if strings.HasPrefix(strings.ToLower(addr), prefix) {
			addr = addr[len(prefix):]
		}

		idx := strings.Index(addr, ":")
		if idx != -1 {
			ipv4Reg := regexp.MustCompile(ipv4Pattern)
			matches := ipv4Reg.FindAllString(addr[:idx], -1)
			if len(matches) == 1 {
				formattedAddr = matches[0]

				portReg := regexp.MustCompile(portPattern)
				matches := portReg.FindAllString(addr[idx+1:], -1)
				if len(matches) == 1 {
					formattedAddr += ":" + matches[0]
				} else {
					formattedAddr = ""
				}
			}
		}

		if formattedAddr == "" {
			return nil, fmt.Errorf("signer addr is not correct: %v", addr)
		}

		addr = formattedAddr
	} else {
		addr = matches[0]
	}

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
