# Data Model

###



### Blob Request

```go
type BlobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The hash of the ReducedBatchHeader defined onchain, see:
	// https://github.com/zero-gravity-labs/zerog-data-avail/blob/master/contracts/src/interfaces/IZGDAServiceManager.sol#L43
	// This identifies the batch that this blob belongs to.
	BatchHeaderHash []byte `protobuf:"bytes,1,opt,name=batch_header_hash,json=batchHeaderHash,proto3" json:"batch_header_hash,omitempty"`
	// Which blob in the batch this is requesting for (note: a batch is logically an
	// ordered list of blobs).
	BlobIndex uint32 `protobuf:"varint,2,opt,name=blob_index,json=blobIndex,proto3" json:"blob_index,omitempty"`
	// DEPRECATED
	ReferenceBlockNumber uint32 `protobuf:"varint,3,opt,name=reference_block_number,json=referenceBlockNumber,proto3" json:"reference_block_number,omitempty"`
	// DEPRECATED
	QuorumId uint32 `protobuf:"varint,4,opt,name=quorum_id,json=quorumId,proto3" json:"quorum_id,omitempty"`
}
```

### Blob Header

```go
type BlobHeader struct {
	// KZG commitments of the blob
	BlobCommitments
	// DEPRECATED
	QuorumInfos []*BlobQuorumInfo
	// DEPRECATED
	AccountID AccountID `json:"account_id"`
}

// BlomCommitments contains the blob's commitment, degree proof, and the actual degree.
type BlobCommitments struct {
	Commitment  *Commitment
	LengthProof *Commitment
	Length      uint
}
```

### Batch Headers

```go
// BatchHeader contains the metadata associated with a Batch for which DA nodes must attest; DA nodes sign on the hash of the batch header
type BatchHeader struct {
	// BlobHeaders contains the headers of the blobs in the batch
	BlobHeaders []*BlobHeader
	// QuorumResults contains the quorum parameters for each quorum that must sign the batch; all quorum parameters must be satisfied
	// for the batch to be considered valid
	QuorumResults []QuorumResult
	// ReferenceBlockNumber is the block number at which all operator information (stakes, indexes, etc.) is taken from
	ReferenceBlockNumber uint
	// BatchRoot is the root of a Merkle tree whose leaves are the hashes of the blobs in the batch
	BatchRoot [32]byte
}
```

### Encoded Data Products

```go
// EncodedBatch is a container for a batch of blobs. DA nodes receive and attest to the blobs in a batch together to amortize signature verification costs
type EncodedBatch struct {
	ChunkBatches map[OperatorID]ChunkBatch
}

// Chunks

// Chunk is the smallest unit that is distributed to DA nodes, including both data and the associated polynomial opening proofs.
// A chunk corresponds to a set of evaluations of the global polynomial whose coefficients are used to construct the blob Commitment.
type Chunk struct {
	// The Coeffs field contains the coefficients of the polynomial which interpolates these evaluations. This is the same as the
	// interpolating polynomial, I(X), used in the KZG multi-reveal (https://dankradfeist.de/ethereum/2020/06/16/kate-polynomial-commitments.html#multiproofs)
	Coeffs []Symbol
	Proof  Proof
}

func (c *Chunk) Length() int {
	return len(c.Coeffs)
}

// ChunkBatch is the collection of chunks associated with a single operator and a single batch.
type ChunkBatch struct {
	// Bundles contains the chunks associated with each blob in the batch; each bundle contains the chunks associated with a single blob
	// The number of bundles should be equal to the total number of blobs in the batch. The number of chunks per bundle will vary
	Bundles [][]*Chunk
}
```

