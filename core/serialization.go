package core

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"regexp"

	"github.com/0glabs/0g-data-avail/pkg/kzg/bn254"
	bn "github.com/consensys/gnark-crypto/ecc/bn254"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"
	"golang.org/x/crypto/sha3"
)

var ErrInvalidCommitment = errors.New("invalid commitment")

func ComputeSignatoryRecordHash(referenceBlockNumber uint32, nonSignerKeys []*G1Point) [32]byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, referenceBlockNumber)
	for _, nonSignerKey := range nonSignerKeys {
		hash := nonSignerKey.GetOperatorID()
		buf = append(buf, hash[:]...)
	}

	var res [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(buf)
	copy(res[:], hasher.Sum(nil)[:32])

	return res
}

// SetBatchRoot sets the BatchRoot field of the BatchHeader to the Merkle root of the blob headers in the batch (i.e. the root of the Merkle tree whose leaves are the blob headers)
func (h *BatchHeader) SetBatchRoot(blobHeaders []*BlobHeader) (*merkletree.MerkleTree, error) {
	leafs := make([][]byte, len(blobHeaders))
	for i, header := range blobHeaders {
		leaf, err := header.GetBlobHeaderHash()
		if err != nil {
			return nil, fmt.Errorf("failed to compute blob header hash: %w", err)
		}
		leafs[i] = leaf[:]
	}

	tree, err := merkletree.NewTree(merkletree.WithData(leafs), merkletree.WithHashType(keccak256.New()))
	if err != nil {
		return nil, err
	}

	copy(h.BatchRoot[:], tree.Root())
	return tree, nil
}

func (h *BatchHeader) Encode() ([]byte, error) {
	// The order here has to match the field ordering of ReducedBatchHeader defined in IZGDAServiceManager.sol
	// ref: https://github.com/0glabs/0g-data-avail/blob/master/contracts/src/interfaces/IZGDAServiceManager.sol#L43
	batchHeaderType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "blobHeadersRoot",
			Type: "bytes32",
		},
		{
			Name: "referenceBlockNumber",
			Type: "uint32",
		},
	})
	if err != nil {
		return nil, err
	}

	arguments := abi.Arguments{
		{
			Type: batchHeaderType,
		},
	}

	s := struct {
		BlobHeadersRoot      [32]byte
		ReferenceBlockNumber uint32
	}{
		BlobHeadersRoot:      h.BatchRoot,
		ReferenceBlockNumber: uint32(h.ReferenceBlockNumber),
	}

	bytes, err := arguments.Pack(s)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// GetBatchHeaderHash returns the hash of the reduced BatchHeader that is used to sign the Batch
// ref: https://github.com/0glabs/0g-data-avail/blob/master/contracts/src/libraries/ZGDAHasher.sol#L65
func (h BatchHeader) GetBatchHeaderHash() ([32]byte, error) {
	headerByte, err := h.Encode()
	if err != nil {
		return [32]byte{}, err
	}

	var headerHash [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(headerByte)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return headerHash, nil
}

// GetBlobHeaderHash returns the hash of the BlobHeader that is used to sign the Blob
func (h BlobHeader) GetBlobHeaderHash() ([32]byte, error) {
	headerByte, err := h.Encode()
	if err != nil {
		return [32]byte{}, err
	}

	var headerHash [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(headerByte)
	copy(headerHash[:], hasher.Sum(nil)[:32])

	return headerHash, nil
}

func (h *BlobHeader) GetQuorumBlobParamsHash() ([32]byte, error) {
	quorumBlobParamsType, err := abi.NewType("tuple[]", "", []abi.ArgumentMarshaling{
		{
			Name: "quorumNumber",
			Type: "uint8",
		},
		{
			Name: "adversaryThresholdPercentage",
			Type: "uint8",
		},
		{
			Name: "quorumThresholdPercentage",
			Type: "uint8",
		},
		{
			Name: "quantizationParameter",
			Type: "uint8",
		},
	})

	if err != nil {
		return [32]byte{}, err
	}

	arguments := abi.Arguments{
		{
			Type: quorumBlobParamsType,
		},
	}

	type quorumBlobParams struct {
		QuorumNumber                 uint8
		AdversaryThresholdPercentage uint8
		QuorumThresholdPercentage    uint8
		QuantizationParameter        uint8
	}

	qbp := make([]quorumBlobParams, len(h.QuorumInfos))
	for i, q := range h.QuorumInfos {
		qbp[i] = quorumBlobParams{
			QuorumNumber:                 uint8(q.QuorumID),
			AdversaryThresholdPercentage: uint8(q.AdversaryThreshold),
			QuorumThresholdPercentage:    uint8(q.QuorumThreshold),
			QuantizationParameter:        0,
		}
	}

	bytes, err := arguments.Pack(qbp)
	if err != nil {
		return [32]byte{}, err
	}

	var res [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(bytes)
	copy(res[:], hasher.Sum(nil)[:32])

	return res, nil
}

func (h *BlobHeader) Encode() ([]byte, error) {
	if h.Commitment == nil || h.Commitment.G1Point == nil {
		return nil, ErrInvalidCommitment
	}

	// The order here has to match the field ordering of BlobHeader defined in IZGDAServiceManager.sol
	blobHeaderType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "commitment",
			Type: "tuple",
			Components: []abi.ArgumentMarshaling{
				{
					Name: "X",
					Type: "uint256",
				},
				{
					Name: "Y",
					Type: "uint256",
				},
			},
		},
		{
			Name: "dataLength",
			Type: "uint32",
		},
		{
			Name: "quorumBlobParams",
			Type: "tuple[]",
			Components: []abi.ArgumentMarshaling{
				{
					Name: "quorumNumber",
					Type: "uint8",
				},
				{
					Name: "adversaryThresholdPercentage",
					Type: "uint8",
				},
				{
					Name: "quorumThresholdPercentage",
					Type: "uint8",
				},
				{
					Name: "quantizationParameter",
					Type: "uint8",
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	arguments := abi.Arguments{
		{
			Type: blobHeaderType,
		},
	}

	type quorumBlobParams struct {
		QuorumNumber                 uint8
		AdversaryThresholdPercentage uint8
		QuorumThresholdPercentage    uint8
		QuantizationParameter        uint8
	}

	type commitment struct {
		X *big.Int
		Y *big.Int
	}

	qbp := make([]quorumBlobParams, len(h.QuorumInfos))
	for i, q := range h.QuorumInfos {
		qbp[i] = quorumBlobParams{
			QuorumNumber:                 uint8(q.QuorumID),
			AdversaryThresholdPercentage: uint8(q.AdversaryThreshold),
			QuorumThresholdPercentage:    uint8(q.QuorumThreshold),
			QuantizationParameter:        0,
		}
	}

	s := struct {
		Commitment       commitment
		DataLength       uint32
		QuorumBlobParams []quorumBlobParams
	}{
		Commitment: commitment{
			X: h.Commitment.X.BigInt(new(big.Int)),
			Y: h.Commitment.Y.BigInt(new(big.Int)),
		},
		DataLength:       uint32(h.Length),
		QuorumBlobParams: qbp,
	}

	bytes, err := arguments.Pack(s)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (h *BatchHeader) Serialize() ([]byte, error) {
	return Encode(h)
}

func (h *BatchHeader) Deserialize(data []byte) (*BatchHeader, error) {
	err := Decode(data, h)
	return h, err
}

func (h *BlobHeader) Serialize() ([]byte, error) {
	return Encode(h)
}

func (h *BlobHeader) Deserialize(data []byte) (*BlobHeader, error) {
	err := Decode(data, h)
	return h, err
}

func (c *Chunk) Serialize() ([]byte, error) {
	return Encode(c)
}

func (c *Chunk) Deserialize(data []byte) (*Chunk, error) {
	err := Decode(data, c)
	return c, err
}

func (c Commitment) Serialize() ([]byte, error) {
	return Encode(c)
}

func (c *Commitment) Deserialize(data []byte) (*Commitment, error) {
	err := Decode(data, c)
	return c, err
}

func (c *Commitment) UnmarshalJSON(data []byte) error {
	var g1Point bn.G1Affine
	err := json.Unmarshal(data, &g1Point)
	if err != nil {
		return err
	}
	c.G1Point = &bn254.G1Point{
		X: g1Point.X,
		Y: g1Point.Y,
	}

	return nil
}

func (h *KVBlobInfoKey) Bytes() []byte {
	b := make([]byte, 0)
	b = append(b, h.BatchHeaderHash[:]...)
	blobIndex := make([]byte, 4)
	binary.BigEndian.PutUint32(blobIndex, h.BlobIndex)
	b = append(b, blobIndex...)
	return b
}

func (h *KVBlobInfoKey) FromBytes(data []byte) (*KVBlobInfoKey, error) {
	if len(data) < 36 {
		return nil, errors.New("data length mismatch")
	}
	copy(h.BatchHeaderHash[:], data[:32])
	h.BlobIndex = binary.BigEndian.Uint32(data[32:36])
	return h, nil
}

func (h *KVBlobInfo) Serialize() ([]byte, error) {
	return Encode(h)
}

func (h *KVBlobInfo) Deserialize(data []byte) (*KVBlobInfo, error) {
	err := Decode(data, h)
	return h, err
}

func (h *KVBatchInfo) Serialize() ([]byte, error) {
	return Encode(h)
}

func (h *KVBatchInfo) Deserialize(data []byte) (*KVBatchInfo, error) {
	err := Decode(data, h)
	return h, err
}

func Encode(obj any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(data []byte, obj any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}

func (s OperatorSocket) GetDispersalSocket() string {
	ip, port1, _, err := extractIPAndPorts(string(s))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s:%s", ip, port1)
}

func (s OperatorSocket) GetRetrievalSocket() string {
	ip, _, port2, err := extractIPAndPorts(string(s))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s:%s", ip, port2)
}

func extractIPAndPorts(s string) (string, string, string, error) {
	regex := regexp.MustCompile(`^([^:]+):([^;]+);([^;]+)$`)
	matches := regex.FindStringSubmatch(s)

	if len(matches) != 4 {
		return "", "", "", errors.New("input string does not match expected format")
	}

	ip := matches[1]
	port1 := matches[2]
	port2 := matches[3]

	return ip, port1, port2, nil
}
