package batcher

import (
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/0glabs/0g-da-client/core"
	"github.com/stretchr/testify/assert"
)

func TestGetHash(t *testing.T) {
	dataRoot := sha256.Sum256([]byte("dataRoot"))
	epoch := big.NewInt(123)
	quorumId := big.NewInt(456)
	erasureCommitment := core.NewG1Point(new(big.Int).SetUint64(1), new(big.Int).SetUint64(1))

	expectedHash := [32]byte{0xde, 0x7b, 0xb4, 0x32, 0xe4, 0xff, 0xf3, 0xff, 0xbd, 0x59, 0x3c, 0x99, 0x6a, 0x9a, 0x60, 0x62, 0x6d, 0x24, 0xa4, 0xaa, 0xc0, 0xa5, 0xd0, 0xbb, 0x49, 0x47, 0x66, 0x48, 0x92, 0x42, 0x91, 0xe}

	resultHash, err := getHash(dataRoot, epoch, quorumId, erasureCommitment)
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, resultHash, "Hashes should match")
}
