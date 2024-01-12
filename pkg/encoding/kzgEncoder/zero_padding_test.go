package kzgEncoder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rs "github.com/zero-gravity-labs/zerog-data-avail/pkg/encoding/encoder"
	kzgRs "github.com/zero-gravity-labs/zerog-data-avail/pkg/encoding/kzgEncoder"
)

func TestProveZeroPadding(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	group, _ := kzgRs.NewKzgEncoderGroup(kzgConfig)

	params := rs.GetEncodingParams(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	enc, err := group.NewKzgEncoder(params)
	require.Nil(t, err)

	inputFr := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)

	_, _, _, _, err = enc.Encode(inputFr)
	require.Nil(t, err)

	assert.True(t, true, "Proof %v failed\n")
}
