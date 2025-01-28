package core

import (
	"math/big"

	eth_common "github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"

	"github.com/0glabs/0g-da-client/common"
)

type AccountID = string

// Security and Quorum Parameters

// QuorumID is a unique identifier for a quorum; initially ZGDA wil support upt to 256 quorums
type QuorumID = uint8

// SecurityParam contains the quorum ID and the adversary threshold for the quorum;
type SecurityParam struct {
	QuorumID QuorumID `json:"quorum_id"`
	// AdversaryThreshold is the maximum amount of stake that can be controlled by an adversary in the quorum as a percentage of the total stake in the quorum
	AdversaryThreshold uint8 `json:"adversary_threshold"`
	// QuorumThreshold is the amount of stake that must sign a message for it to be considered valid as a percentage of the total stake in the quorum
	QuorumThreshold uint8 `json:"quorum_threshold"`
	// Rate Limit. This is a temporary measure until the node can derive rates on its own using rollup authentication. This is used
	// for restricting the rate at which retrievers are able to download data from the DA node to a multiple of the rate at which the
	// data was posted to the DA node.
	QuorumRate common.RateParam `json:"quorum_rate"`
}

// QuorumResult contains the quorum ID and the amount signed for the quorum
type QuorumResult struct {
	QuorumID QuorumID
	// PercentSigned is percentage of the total stake for the quorum that signed for a particular batch.
	PercentSigned uint8
}

// Blob stores the data and header of a single data blob. Blobs are the fundamental unit of data posted to ZGDA by users.
type Blob struct {
	RequestHeader BlobRequestHeader
	Data          []byte
}

// BlobRequestHeader contains the original data size of a blob and the security required
type BlobRequestHeader struct {
	// For a blob to be accepted by ZGDA, it satisfy the AdversaryThreshold of each quorum contained in SecurityParams
	SecurityParams []*SecurityParam `json:"security_params"`
	// AccountID is the account that is paying for the blob to be stored
	AccountID AccountID `json:"account_id"`
}

// BlobQuorumInfo contains the quorum IDs and parameters for a blob specific to a given quorum
type BlobQuorumInfo struct {
	SecurityParam
	// ChunkLength is the number of symbols in a chunk
	ChunkLength uint
}

// BlobHeader contains all metadata related to a blob including commitments and parameters for encoding
type BlobHeader struct {
	// CommitmentRoot the root of merkle tree of kzg commitments
	CommitmentRoot []byte `json:"commitment_root"`
	Length         uint   `json:"length"`
}

type Coeff = [32]byte

type EncodedRow = []Coeff

type Commitment = [48]byte

type CommitRootSubmission struct {
	DataRoot          [32]byte
	Epoch             *big.Int
	QuorumId          *big.Int
	ErasureCommitment *G1Point
	QuorumBitmap      []byte
	AggPkG2           *G2Point
	AggSigs           *Signature
}

type BlobCommitments struct {
	ErasureCommitment *G1Point
	StorageRoot       []byte
	EncodedData       []byte
	EncodedSlice      [][]byte
}

func (b *BlobCommitments) GetHash() [32]byte {
	var message [32]byte
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(b.ErasureCommitment.Serialize())
	hasher.Write(b.StorageRoot)
	copy(message[:], hasher.Sum(nil)[:32])
	return message
}

type ExtendedMatrix struct {
	Length      uint
	Rows        []EncodedRow
	Commitments []Commitment
}

func (m *ExtendedMatrix) GetRows() int {
	return len(m.Rows)
}

func (m *ExtendedMatrix) GetCols() int {
	return len(m.Rows[0])
}

func (m *ExtendedMatrix) GetRowInBytes(idx int) []byte {
	result := make([]byte, 0)
	for _, chunk := range m.Rows[idx] {
		result = append(result, chunk[:]...)
	}
	return result
}

// Batch
// A batch is a collection of blobs. DA nodes receive and attest to the blobs in a batch together to amortize signature verification costs

// BatchHeader contains the metadata associated with a Batch for which DA nodes must attest; DA nodes sign on the hash of the batch header
type BatchHeader struct {
	// BatchRoot is the root of a Merkle tree whose leaves are the hashes of the blobs in the batch
	BatchRoot [32]byte `json:"batch_root"`
	// DataRoot is the root of a Merkle tree whos leaves are the merkle root encoded data blobs divided in zgs segment size
	DataRoot eth_common.Hash `json:"data_root"`
}

type KVBlobInfoKey struct {
	BatchHeaderHash [32]byte `json:"batch_header_hash"`
	BlobIndex       uint32   `json:"blob_index"`
}

type MerkleProof struct {
	Hashes [][]byte `json:"hashes"`
	Index  uint64   `json:"index"`
}

// BlobInfo to write to KV stream, the key is (BatchHeaderHash, BlobIndex)
type KVBlobInfo struct {
	BlobHeader  *BlobHeader  `json:"blob_header"`
	MerkleProof *MerkleProof `json:"merkle_proof"` // to prove blob exists in batch
}

type BlobDisperseInfo struct {
	BlobLength uint `json:"blob_length"`
	Rows       uint `json:"rows"`
	Cols       uint `json:"cols"`
}

// KVBatchInfo to write to KV stream, blob lengths are used to recompute the chunks allocations
type KVBatchInfo struct {
	BatchHeader       *BatchHeader       `json:"batch_header"`
	BlobDisperseInfos []BlobDisperseInfo `json:"blob_disperse_infos"`
}

// BlobLocation the segment index and offset of rows in encoded blob
type BlobLocation struct {
	Rows           uint
	Cols           uint
	SegmentIndexes []uint
	Offsets        []uint
}
