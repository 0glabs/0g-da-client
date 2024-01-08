package batcher

import (
	"fmt"
	"sync"

	"github.com/zero-gravity-labs/zgda/common"
	"github.com/zero-gravity-labs/zgda/core"
	"github.com/zero-gravity-labs/zgda/disperser"
)

type requestID string

type encodedBlobStore struct {
	mu sync.RWMutex

	requested map[requestID]struct{}
	encoded   map[requestID]*EncodingResult
	// encodedResultSize is the total size of all the chunks in the encoded results in bytes
	encodedResultSize uint64

	logger common.Logger
}

// EncodingResult contains information about the encoding of a blob
type EncodingResult struct {
	BlobMetadata         *disperser.BlobMetadata
	ReferenceBlockNumber uint
	// BlobQuorumInfo       *core.BlobQuorumInfo
	Commitment *core.BlobCommitments
	Chunks     []*core.Chunk
	// Assignments          map[core.OperatorID]core.Assignment
}

// EncodingResultOrStatus is a wrapper for EncodingResult that also contains an error
type EncodingResultOrStatus struct {
	EncodingResult
	// Err is set if there was an error during encoding
	Err error
}

func newEncodedBlobStore(logger common.Logger) *encodedBlobStore {
	return &encodedBlobStore{
		requested:         make(map[requestID]struct{}),
		encoded:           make(map[requestID]*EncodingResult),
		encodedResultSize: 0,
		logger:            logger,
	}
}

func (e *encodedBlobStore) PutEncodingRequest(blobKey disperser.BlobKey) {
	e.mu.Lock()
	defer e.mu.Unlock()

	requestID := getRequestID(blobKey)
	e.requested[requestID] = struct{}{}
}

func (e *encodedBlobStore) HasEncodingRequested(blobKey disperser.BlobKey, quorumID core.QuorumID, referenceBlockNumber uint) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	requestID := getRequestID(blobKey)
	if _, ok := e.requested[requestID]; ok {
		return true
	}

	res, ok := e.encoded[requestID]
	if ok && res.ReferenceBlockNumber == referenceBlockNumber {
		return true
	}
	return false
}

func (e *encodedBlobStore) DeleteEncodingRequest(blobKey disperser.BlobKey) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	requestID := getRequestID(blobKey)
	if _, ok := e.requested[requestID]; !ok {
		return
	}

	delete(e.requested, requestID)
}

func (e *encodedBlobStore) PutEncodingResult(result *EncodingResult) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	blobKey := disperser.BlobKey{
		BlobHash:     result.BlobMetadata.BlobHash,
		MetadataHash: result.BlobMetadata.MetadataHash,
	}
	requestID := getRequestID(blobKey)
	if _, ok := e.requested[requestID]; !ok {
		return fmt.Errorf("PutEncodedBlob: no such key (%s) in requested set", requestID)
	}

	if _, ok := e.encoded[requestID]; !ok {
		e.encodedResultSize += getChunksSize(result)
	}
	e.encoded[requestID] = result
	delete(e.requested, requestID)

	return nil
}

func (e *encodedBlobStore) GetEncodingResult(blobKey disperser.BlobKey, quorumID core.QuorumID) (*EncodingResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	requestID := getRequestID(blobKey)
	if _, ok := e.encoded[requestID]; !ok {
		return nil, fmt.Errorf("GetEncodedBlob: no such key (%s) in encoded set", requestID)
	}

	return e.encoded[requestID], nil
}

func (e *encodedBlobStore) DeleteEncodingResult(blobKey disperser.BlobKey, quorumID core.QuorumID) {
	e.mu.Lock()
	defer e.mu.Unlock()

	requestID := getRequestID(blobKey)
	encodedResult, ok := e.encoded[requestID]
	if !ok {
		return
	}

	delete(e.encoded, requestID)
	e.encodedResultSize -= getChunksSize(encodedResult)
}

// GetNewEncodingResults returns all the fresh encoded results
func (e *encodedBlobStore) GetNewEncodingResults() []*EncodingResult {
	e.mu.RLock()
	defer e.mu.RUnlock()
	fetched := make([]*EncodingResult, 0)
	for _, encodedResult := range e.encoded {
		fetched = append(fetched, encodedResult)
	}
	e.logger.Trace("consumed encoded results", "fetched", len(fetched), "encodedSize", e.encodedResultSize)
	return fetched
}

// DeleteStaleEncodingResults deletes all the stale results
func (e *encodedBlobStore) DeleteStaleEncodingResults() {
	e.mu.Lock()
	defer e.mu.Unlock()
	for ri := range e.encoded {
		delete(e.encoded, ri)
	}
}

// GetEncodedResultSize returns the total size of all the chunks in the encoded results in bytes
func (e *encodedBlobStore) GetEncodedResultSize() (int, uint64) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return len(e.encoded), e.encodedResultSize
}

func getRequestID(key disperser.BlobKey) requestID {
	return requestID(fmt.Sprintf("%s", key.String()))
}

func getChunksSize(result *EncodingResult) uint64 {
	var size uint64

	for _, chunk := range result.Chunks {
		size += uint64(len(chunk.Coeffs) * 256) // 256 bytes per symbol
	}
	return size + 256*2 // + 256 * 2 bytes for proof
}
