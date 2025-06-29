# Data Model

## Disperser

### Request

```go
type DisperseBlobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The data to be dispersed.
	// The size of data must be <= 512KiB.
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	SecurityParams []*SecurityParams `protobuf:"bytes,2,rep,name=security_params,json=securityParams,proto3" json:"security_params,omitempty"`
	// The number of chunks that encoded blob split into.
	// The number will be aligned to the next power of 2 and be bounded by blob size.
	uint32 targetChunkNum = 3;
}

type RetrieveBlobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BatchHeaderHash []byte `protobuf:"bytes,1,opt,name=batch_header_hash,json=batchHeaderHash,proto3" json:"batch_header_hash,omitempty"`
	BlobIndex       uint32 `protobuf:"varint,2,opt,name=blob_index,json=blobIndex,proto3" json:"blob_index,omitempty"`
}
```

### Blob Key

```go
type BlobKey struct {
	BlobHash     BlobHash
	MetadataHash MetadataHash
}
```

### Blob Metadata

```go
type BlobMetadata struct {
	BlobHash     BlobHash     `json:"blob_hash"`
	MetadataHash MetadataHash `json:"metadata_hash"`
	BlobStatus   BlobStatus   `json:"blob_status"`
	// Expiry is unix epoch time in seconds at which the blob will expire
	Expiry uint64 `json:"expiry"`
	// NumRetries is the number of times the blob has been retried
	// After few failed attempts, the blob will be marked as failed
	NumRetries uint `json:"num_retries"`
	// RequestMetadata is the request metadata of the blob when it was requested
	// This field is omitted when marshalling to DynamoDB attributevalue as this field will be flattened
	RequestMetadata *RequestMetadata `json:"request_metadata" dynamodbav:"-"`
	// ConfirmationInfo is the confirmation metadata of the blob when it was confirmed
	// This field is nil if the blob has not been confirmed
	// This field is omitted when marshalling to DynamoDB attributevalue as this field will be flattened
	ConfirmationInfo *ConfirmationInfo `json:"blob_confirmation_info" dynamodbav:"-"`
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

## Batcher

### Batch Header

```go
type BatchHeader struct {
	// ReferenceBlockNumber is the block number at which all operator information (stakes, indexes, etc.) is taken from
	ReferenceBlockNumber uint
	// BatchRoot is the root of a Merkle tree whose leaves are the hashes of the blobs in the batch
	BatchRoot [32]byte
	// DataRoot is the root of a Merkle tree who's leaves are the merkle root encoded data blobs divided in zgs segment size
	DataRoot eth_common.Hash
}
```

### Encoded Blob

```go
type EncodedBlob = BlobMessage

type BlobMessage struct {
	BlobHeader *BlobHeader
	Bundles    Bundles
}

type Bundles = []*Chunk

type Chunk struct {
	// The Coeffs field contains the coefficients of the polynomial which interolates these evaluations. This is the same as the
	// interpolating polynomial, I(X), used in the KZG multi-reveal (https://dankradfeist.de/ethereum/2020/06/16/kate-polynomial-commitments.html#multiproofs)
	Coeffs []Symbol
	Proof  Proof
}
```

### Merkle Tree Proof

```go
// Proof is a proof of a Merkle tree.
type Proof struct {
	Hashes [][]byte
	Index  uint64
}
```

## Retriever

### Blob Request

```go
type BlobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The hash of the ReducedBatchHeader defined onchain, see:
	// https://github.com/0glab/0g-data-avail/blob/master/contracts/src/interfaces/IZGDAServiceManager.sol#L43
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
