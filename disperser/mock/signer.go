package mock

import (
	"context"

	"github.com/0glabs/0g-da-client/common"
	"github.com/0glabs/0g-da-client/core"
	"github.com/0glabs/0g-da-client/disperser"
	pb "github.com/0glabs/0g-da-client/disperser/api/grpc/signer"
	"github.com/stretchr/testify/mock"
)

type MockSignerClient struct {
	mock.Mock
}

var _ disperser.SignerClient = (*MockSignerClient)(nil)

func NewMockSignerClient() *MockSignerClient {
	return &MockSignerClient{}
}

func (m *MockSignerClient) BatchSign(ctx context.Context, addr string, data []*pb.SignRequest, log common.Logger) ([]*core.Signature, error) {
	args := m.Called(ctx, addr, data, log)
	var signatures []*core.Signature
	if args.Get(0) != nil {
		signatures = args.Get(0).([]*core.Signature)
	}

	return signatures, args.Error(1)
}
