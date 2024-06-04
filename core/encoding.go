package core

import (
	"math"
)

const (
	EntrySize       = 256 //256B
	EntryPerSegment = 1024
	SegmentSize     = EntrySize * EntryPerSegment // 256KB
	ScalarSize      = 31
	CoeffSize       = 32
	CommitmentSize  = 48
	MaxCols         = 1024
	MaxRows         = 1024
	MaxBlobSize     = MaxCols * MaxRows * ScalarSize
)

type MatrixDimsions struct {
	Rows uint `json:"rows"`
	Cols uint `json:"cols"`
}

// Encoder is responsible for encoding
type Encoder interface {
	// Encode takes in a blob and returns the commitments and encoded matrix
	Encode(data []byte, dims MatrixDimsions) (*ExtendedMatrix, error)
}

func NextPowerOf2(d uint64) uint64 {
	nextPower := math.Ceil(math.Log2(float64(d)))
	return uint64(math.Pow(2.0, nextPower))
}

// GetBlobLength converts from blob size in bytes to blob size in symbols
func GetBlobLength(blobSize uint) uint {
	return (blobSize + ScalarSize - 1) / ScalarSize
}

// SplitToMatrix calculate row and column length for encoded blob, try to split it into rows x cols matrix
func SplitToMatrix(blobLength uint, targetRowNum uint) (uint, uint) {
	expectedLength := uint(NextPowerOf2(uint64(blobLength * 2)))
	var rows, cols uint
	if targetRowNum == 0 {
		// split into maximum rows
		rows = min(expectedLength, MaxRows)
		cols = expectedLength / rows
	} else {
		// try to split into target rows
		targetRowNum = min(MaxRows, uint(NextPowerOf2(uint64(targetRowNum))))
		rows = min(expectedLength, targetRowNum)
		cols = expectedLength / rows
		if cols > MaxCols {
			// split into maximum rows
			rows = min(expectedLength, MaxRows)
			cols = expectedLength / rows
		}
	}
	return rows, cols
}

// GetBlobSize converts from blob length in symbols to blob size in bytes. This is not an exact conversion.
func GetBlobSize(blobLength uint) uint {
	return blobLength * ScalarSize
}
