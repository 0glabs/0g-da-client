package memorydb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/0glabs/0g-data-avail/common"
	"github.com/0glabs/0g-data-avail/core"
	"github.com/0glabs/0g-data-avail/disperser"
)

// SharedBlobStore is an in-memory implementation of the SharedBlobStore interface
type SharedBlobStore struct {
	mu        sync.RWMutex
	Blobs     map[string]*BlobHolder
	Metadata  map[disperser.BlobKey]*disperser.BlobMetadata
	sizeLimit uint64
	size      uint64

	logger common.Logger
}

// BlobHolder stores the blob along with its status and any other metadata
type BlobHolder struct {
	Data []byte
}

var _ disperser.BlobStore = (*SharedBlobStore)(nil)

// NewBlobStore creates an empty BlobStore
func NewBlobStore(sizeLimit uint64, logger common.Logger) disperser.BlobStore {
	return &SharedBlobStore{
		Blobs:     make(map[string]*BlobHolder),
		Metadata:  make(map[disperser.BlobKey]*disperser.BlobMetadata),
		sizeLimit: sizeLimit,
		logger:    logger,
	}
}

func sizeOf(metadata *disperser.BlobMetadata) uint64 {
	var size uint64
	size += 16 + uint64(len(metadata.BlobHash))
	size += 16 + uint64(len(metadata.MetadataHash))
	size += 40 // other fields
	if metadata.RequestMetadata != nil {
		// AccountID
		size += 16 + uint64(len(metadata.RequestMetadata.AccountID))
		// blob commitments: 64+64+8
		// security params: 24
		// TargetChunkNum: 8
		// BlobSize: 8
		// RequestedAt: 8
		size += 184
	}
	if metadata.ConfirmationInfo != nil {
		// BatchRoot
		size += 24 + uint64(len(metadata.ConfirmationInfo.BatchRoot))
		// Inclusion Proof
		size += 24 + uint64(len(metadata.ConfirmationInfo.BlobInclusionProof))
		// BlobCommitments
		if metadata.ConfirmationInfo.BlobCommitment != nil {
			size += 136
		}
		// Fee
		size += 24 + uint64(len(metadata.ConfirmationInfo.Fee))
		// BatchHeaderHash: 32
		// BlobIndex: 8
		// BlobCount: 8
		// SignatoryRecordHash: 32
		// ReferenceBlockNumber: 8
		// BlobCommitment: 8
		// BatchID: 4
		// ConfirmationTxnHash: 32
		// ConfirmationBlockNumber: 4
		// QuorumResults: 8
		// BlobQuorumInfos: 24
		size += 168
	}
	return size
}

func (q *SharedBlobStore) MetadataHashAsBlobKey() bool {
	return true
}

func (q *SharedBlobStore) RemoveBlob(ctx context.Context, metadata *disperser.BlobMetadata) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if holder, ok := q.Blobs[metadata.MetadataHash]; ok {
		q.size -= uint64(len(holder.Data))
		delete(q.Blobs, metadata.MetadataHash)
	}
	if existing, ok := q.Metadata[metadata.GetBlobKey()]; ok {
		q.size -= sizeOf(existing)
		delete(q.Metadata, metadata.GetBlobKey())
	}
	q.logger.Info("[memdb] blob removed", "mem db used", q.size, "limit", q.sizeLimit)
	return nil
}

func (q *SharedBlobStore) StoreBlob(ctx context.Context, blob *core.Blob, requestedAt uint64) (disperser.BlobKey, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	blobKey := disperser.BlobKey{}
	// Generate the blob key
	blobHash := getBlobHash(blob)
	blobKey.BlobHash = blobHash
	blobKey.MetadataHash = getMetadataHash(requestedAt)

	if _, ok := q.Blobs[blobKey.MetadataHash]; !ok {
		q.size += uint64(len(blob.Data))
		if q.size > q.sizeLimit {
			return blobKey, disperser.ErrMemoryDbIsFull
		}
		// Add the blob to the queue
		q.Blobs[blobKey.MetadataHash] = &BlobHolder{
			Data: blob.Data,
		}
	}

	if _, ok := q.Metadata[blobKey]; !ok {
		metadata := &disperser.BlobMetadata{
			BlobHash:     blobHash,
			MetadataHash: blobKey.MetadataHash,
			BlobStatus:   disperser.Processing,
			NumRetries:   0,
			RequestMetadata: &disperser.RequestMetadata{
				BlobRequestHeader: blob.RequestHeader,
				BlobSize:          uint(len(blob.Data)),
				RequestedAt:       requestedAt,
			},
		}
		q.size += sizeOf(metadata)
		if q.size > q.sizeLimit {
			return blobKey, disperser.ErrMemoryDbIsFull
		}
		q.Metadata[blobKey] = metadata
	}
	q.logger.Info("[memdb] blob stored", "mem db used", q.size, "limit", q.sizeLimit)
	return blobKey, nil
}

func (q *SharedBlobStore) GetBlobContent(ctx context.Context, metadata *disperser.BlobMetadata) ([]byte, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if holder, ok := q.Blobs[metadata.MetadataHash]; ok {
		return holder.Data, nil
	} else {
		return nil, disperser.ErrBlobNotFound
	}
}

func (q *SharedBlobStore) MarkBlobConfirmed(ctx context.Context, existingMetadata *disperser.BlobMetadata, confirmationInfo *disperser.ConfirmationInfo) (*disperser.BlobMetadata, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	// TODO (ian-shim): remove this check once we are sure that the metadata is never overwritten
	refreshedMetadata, err := q.GetBlobMetadata(ctx, existingMetadata.GetBlobKey())
	if err != nil {
		return nil, err
	}
	alreadyConfirmed, _ := refreshedMetadata.IsConfirmed()
	if alreadyConfirmed {
		return refreshedMetadata, nil
	}
	blobKey := existingMetadata.GetBlobKey()
	if _, ok := q.Metadata[blobKey]; !ok {
		return nil, disperser.ErrBlobNotFound
	}
	newMetadata := *existingMetadata
	newMetadata.BlobStatus = disperser.Confirmed
	newMetadata.ConfirmationInfo = confirmationInfo
	// update size
	if existing, ok := q.Metadata[blobKey]; ok {
		q.size -= sizeOf(existing)
	}
	q.size += sizeOf(&newMetadata)
	q.logger.Info("[memdb] blob confirmed", "mem db used", q.size, "limit", q.sizeLimit)
	// don't throw error here
	q.Metadata[blobKey] = &newMetadata
	return &newMetadata, nil
}

func (q *SharedBlobStore) MarkBlobFinalized(ctx context.Context, blobKey disperser.BlobKey) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.Metadata[blobKey]; !ok {
		return disperser.ErrBlobNotFound
	}

	q.Metadata[blobKey].BlobStatus = disperser.Finalized
	return nil
}

func (q *SharedBlobStore) MarkBlobProcessing(ctx context.Context, blobKey disperser.BlobKey) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.Metadata[blobKey]; !ok {
		return disperser.ErrBlobNotFound
	}

	q.Metadata[blobKey].BlobStatus = disperser.Processing
	return nil
}

func (q *SharedBlobStore) MarkBlobFailed(ctx context.Context, blobKey disperser.BlobKey) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.Metadata[blobKey]; !ok {
		return disperser.ErrBlobNotFound
	}

	q.Metadata[blobKey].BlobStatus = disperser.Failed
	return nil
}

func (q *SharedBlobStore) IncrementBlobRetryCount(ctx context.Context, existingMetadata *disperser.BlobMetadata) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.Metadata[existingMetadata.GetBlobKey()]; !ok {
		return disperser.ErrBlobNotFound
	}

	q.Metadata[existingMetadata.GetBlobKey()].NumRetries++
	return nil
}

func (q *SharedBlobStore) GetBlobsByMetadata(ctx context.Context, metadata []*disperser.BlobMetadata) (map[disperser.BlobKey]*core.Blob, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	blobs := make(map[disperser.BlobKey]*core.Blob)
	for _, meta := range metadata {
		if holder, ok := q.Blobs[meta.MetadataHash]; ok {
			blobs[meta.GetBlobKey()] = &core.Blob{
				RequestHeader: meta.RequestMetadata.BlobRequestHeader,
				Data:          holder.Data,
			}
		} else {
			return nil, disperser.ErrBlobNotFound
		}
	}
	return blobs, nil
}

func (q *SharedBlobStore) GetBlobMetadataByStatus(ctx context.Context, status disperser.BlobStatus) ([]*disperser.BlobMetadata, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	metas := make([]*disperser.BlobMetadata, 0)
	for _, meta := range q.Metadata {
		if meta.BlobStatus == status {
			metas = append(metas, meta)
		}
	}
	return metas, nil
}

func (q *SharedBlobStore) GetMetadataInBatch(ctx context.Context, batchHeaderHash [32]byte, blobIndex uint32) (*disperser.BlobMetadata, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	for _, meta := range q.Metadata {
		if meta.ConfirmationInfo != nil && meta.ConfirmationInfo.BatchHeaderHash == batchHeaderHash && meta.ConfirmationInfo.BlobIndex == blobIndex {
			return meta, nil
		}
	}

	return nil, disperser.ErrBlobNotFound
}

func (q *SharedBlobStore) GetAllBlobMetadataByBatch(ctx context.Context, batchHeaderHash [32]byte) ([]*disperser.BlobMetadata, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	metas := make([]*disperser.BlobMetadata, 0)
	for _, meta := range q.Metadata {
		if meta.ConfirmationInfo != nil && meta.ConfirmationInfo.BatchHeaderHash == batchHeaderHash {
			metas = append(metas, meta)
		}
	}
	return metas, nil
}

func (q *SharedBlobStore) GetBlobMetadata(ctx context.Context, blobKey disperser.BlobKey) (*disperser.BlobMetadata, error) {
	if meta, ok := q.Metadata[blobKey]; ok {
		return meta, nil
	}
	return nil, disperser.ErrBlobNotFound
}

func (q *SharedBlobStore) HandleBlobFailure(ctx context.Context, metadata *disperser.BlobMetadata, maxRetry uint) error {
	if metadata.NumRetries < maxRetry {
		return q.IncrementBlobRetryCount(ctx, metadata)
	} else {
		return q.MarkBlobFailed(ctx, metadata.GetBlobKey())
	}
}

func getBlobHash(blob *core.Blob) disperser.BlobHash {
	hasher := sha256.New()
	hasher.Write(blob.Data)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func getMetadataHash(requestedAt uint64) string {
	str := fmt.Sprintf("%d/", requestedAt)
	bytes := []byte(str)
	return hex.EncodeToString(sha256.New().Sum(bytes))
}
