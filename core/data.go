package core

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	eth_common "github.com/ethereum/go-ethereum/common"

	"github.com/wealdtech/go-merkletree"
	"github.com/zero-gravity-labs/zerog-data-avail/common"
	"github.com/zero-gravity-labs/zerog-data-avail/pkg/kzg/bn254"
)

type AccountID = string

// Security and Quorum Paramaters

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

// BlobRequestHeader contains the orignal data size of a blob and the security required
type BlobRequestHeader struct {
	// Commitments
	BlobCommitments `json:"commitments"`
	// For a blob to be accepted by ZGDA, it satisfy the AdversaryThreshold of each quorum contained in SecurityParams
	SecurityParams []*SecurityParam `json:"security_params"`
	// AccountID is the account that is paying for the blob to be stored
	AccountID AccountID `json:"account_id"`
}

func (h *BlobRequestHeader) Validate() error {
	return nil
}

// BlobQuorumInfo contains the quorum IDs and parameters for a blob specific to a given quorum
type BlobQuorumInfo struct {
	SecurityParam
	// ChunkLength is the number of symbols in a chunk
	ChunkLength uint
}

// BlobHeader contains all metadata related to a blob including commitments and parameters for encoding
type BlobHeader struct {
	BlobCommitments
	// QuorumInfos contains the quorum specific parameters for the blob
	QuorumInfos []*BlobQuorumInfo

	// AccountID is the account that is paying for the blob to be stored
	AccountID AccountID `json:"account_id"`
}

func (b *BlobHeader) GetQuorumInfo(quorum QuorumID) *BlobQuorumInfo {
	for _, quorumInfo := range b.QuorumInfos {
		if quorumInfo.QuorumID == quorum {
			return quorumInfo
		}
	}
	return nil
}

// Returns the total encoded size in bytes of the blob across all quorums.
func (b *BlobHeader) EncodedSizeAllQuorums() int64 {
	size := int64(0)
	for _, quorum := range b.QuorumInfos {

		size += int64(roundUpDivide(b.Length*percentMultiplier*bn254.BYTES_PER_COEFFICIENT, uint(quorum.QuorumThreshold-quorum.AdversaryThreshold)))
	}
	return size
}

// BlomCommitments contains the blob's commitment, degree proof, and the actual degree.
type BlobCommitments struct {
	Commitment  *Commitment `json:"commitment"`
	LengthProof *Commitment `json:"length_proof"`
	Length      uint        `json:"length"`
}

// Batch
// A batch is a collection of blobs. DA nodes receive and attest to the blobs in a batch together to amortize signature verification costs

// BatchHeader contains the metadata associated with a Batch for which DA nodes must attest; DA nodes sign on the hash of the batch header
type BatchHeader struct {
	// ReferenceBlockNumber is the block number at which all operator information (stakes, indexes, etc.) is taken from
	ReferenceBlockNumber uint
	// BatchRoot is the root of a Merkle tree whose leaves are the hashes of the blobs in the batch
	BatchRoot [32]byte
	// DataRoot is the root of a Merkle tree whos leaves are the merkle root encoded data blobs divided in zgs segment size
	DataRoot eth_common.Hash
}

// EncodedBlob contains the messages to be sent to a group of DA nodes corresponding to a single blob
type EncodedBlob = BlobMessage

// Chunks

// Chunk is the smallest unit that is distributed to DA nodes, including both data and the associated polynomial opening proofs.
// A chunk corresponds to a set of evaluations of the global polynomial whose coefficients are used to construct the blob Commitment.
type Chunk struct {
	// The Coeffs field contains the coefficients of the polynomial which interolates these evaluations. This is the same as the
	// interpolating polynomial, I(X), used in the KZG multi-reveal (https://dankradfeist.de/ethereum/2020/06/16/kate-polynomial-commitments.html#multiproofs)
	Coeffs []Symbol
	Proof  Proof
}

// Returns 32 * len(Coeffs) bytes
func (c *Chunk) CoeffsToBytes() []byte {
	res := make([]byte, 0)
	for _, coeff := range c.Coeffs {
		b := bn254.FrToBytes(&coeff)
		res = append(res, b[:]...)
	}
	return res
}

func BytesToCoeffs(bytes []byte) []Symbol {
	if len(bytes)%32 != 0 {
		panic("invalid bytes coeffs conversion")
	}
	n := len(bytes) / 32
	res := make([]Symbol, n)
	for i := 0; i < n; i++ {
		var tmp [32]byte
		copy(tmp[:], bytes[i*32:i*32+32])
		bn254.FrFrom32(&res[i], tmp)
	}
	return res
}

// Returns 64 bytes
func (c *Chunk) ProofToBytes() [64]byte {
	x := c.Proof.X.Bytes()
	y := c.Proof.Y.Bytes()
	var res [64]byte
	copy(res[:], append(x[:], y[:]...))
	return res
}

func BytesToProof(bytes [64]byte) Proof {
	var x, y fp.Element
	x.SetBytes(bytes[:32])
	y.SetBytes(bytes[32:])
	return Proof{
		X: x,
		Y: y,
	}
}

func (c *Chunk) Length() int {
	return len(c.Coeffs)
}

// Returns the size of chunk in bytes.
func (c *Chunk) Size() int {
	return c.Length() * bn254.BYTES_PER_COEFFICIENT
}

type Bundles = []*Chunk

// BlobMessage is the message that is sent to DA nodes. It contains the blob header and the associated chunk bundles.
type BlobMessage struct {
	BlobHeader *BlobHeader
	Bundles    Bundles
}

// Sample is a chunk with associated metadata used by the Universal Batch Verifier
type Sample struct {
	Commitment      *Commitment
	Chunk           *Chunk
	AssignmentIndex ChunkNumber
	BlobIndex       int
}

// SubBatch is a part of the whole Batch with identical Encoding Parameters, i.e. (ChunkLen, NumChunk)
// Blobs with the same encoding parameters are collected in a single subBatch
type SubBatch struct {
	Samples  []Sample
	NumBlobs int
}

type KVBlobInfoKey struct {
	BatchHeaderHash [32]byte
	BlobIndex       uint32
}

// BlobInfo to write to KV stream, the key is (BatchHeaderHash, BlobIndex)
type KVBlobInfo struct {
	BlobHeader  *BlobHeader
	MerkleProof *merkletree.Proof // to prove blob exists in batch
	ChunkOffset uint              // start index of first chunk in batch data
	ProofOffset uint              // start index of first proof in batch data
}
