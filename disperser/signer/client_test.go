package signer

import (
	"encoding/hex"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBigEndian(t *testing.T) {
	i, err := hex.DecodeString("b5129dd545319cdef628303e95b3134fd9c6d2b1f025eaba80ca556b661ede0e76daeab01e1bae6aa33379a3af63786ca36e7385472a8839a62e250a2b2a3ca9")
	assert.Nil(t, err, "decode should success")

	// Define test cases with input and expected output
	testCases := []struct {
		input          []byte
		expectedOutput []byte
		expectedError  error
	}{
		{
			i,
			[]byte{0xe, 0xde, 0x1e, 0x66, 0x6b, 0x55, 0xca, 0x80, 0xba, 0xea, 0x25, 0xf0, 0xb1, 0xd2, 0xc6, 0xd9, 0x4f, 0x13, 0xb3, 0x95, 0x3e, 0x30, 0x28, 0xf6, 0xde, 0x9c, 0x31, 0x45, 0xd5, 0x9d, 0x12, 0xb5, 0x29, 0x3c, 0x2a, 0x2b, 0xa, 0x25, 0x2e, 0xa6, 0x39, 0x88, 0x2a, 0x47, 0x85, 0x73, 0x6e, 0xa3, 0x6c, 0x78, 0x63, 0xaf, 0xa3, 0x79, 0x33, 0xa3, 0x6a, 0xae, 0x1b, 0x1e, 0xb0, 0xea, 0xda, 0x76},
			nil,
		},
		{[]byte{1, 2, 3, 4, 5, 6}, nil, io.ErrShortBuffer},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		// Call the function with test input
		output, err := toBigEndian(tc.input)

		// Check if the output and error match the expected values
		assert.Equal(t, tc.expectedError, err, "Error should match")
		assert.Equal(t, tc.expectedOutput, output, "Output should match")
	}
}
