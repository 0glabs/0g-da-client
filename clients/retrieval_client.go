package clients

import (
	"context"

	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"
	"github.com/zero-gravity-labs/zerog-data-avail/common"
	"github.com/zero-gravity-labs/zerog-data-avail/common/storage_node"
	"github.com/zero-gravity-labs/zerog-data-avail/core"
	"github.com/zero-gravity-labs/zerog-data-avail/pkg/kzg/bn254"
	"github.com/zero-gravity-labs/zerog-storage-client/kv"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
)

type RetrievalClient interface {
	RetrieveBlob(
		ctx context.Context,
		batchHeaderHash [32]byte,
		batchDataRoot eth_common.Hash,
		blobIndex uint32,
		referenceBlockNumber uint,
		batchRoot [32]byte,
		blobLengths []uint) ([]byte, error)
}

type retrievalClient struct {
	logger                common.Logger
	indexedChainState     core.IndexedChainState
	assignmentCoordinator core.AssignmentCoordinator
	encoder               core.Encoder
	numConnections        int

	Nodes    []*node.Client
	KVNode   *kv.Client
	StreamId eth_common.Hash
}

var _ RetrievalClient = (*retrievalClient)(nil)

func NewRetrievalClient(
	logger common.Logger,
	encoder core.Encoder,
	numConnections int,
	storageNodeConfig storage_node.ClientConfig,
) (*retrievalClient, error) {

	return &retrievalClient{
		logger:         logger,
		encoder:        encoder,
		numConnections: numConnections,

		Nodes:    node.MustNewClients(storageNodeConfig.StorageNodeURLs),
		KVNode:   kv.NewClient(node.MustNewClient(storageNodeConfig.KVNodeURL), nil),
		StreamId: storageNodeConfig.KVStreamId,
	}, nil
}

func (r *retrievalClient) fetchBlobInfo(batchHeaderHash [32]byte, blobIndex uint32) (*core.KVBlobInfo, error) {
	key := (&core.KVBlobInfoKey{
		BatchHeaderHash: batchHeaderHash,
		BlobIndex:       blobIndex,
	}).Bytes()

	val, err := r.KVNode.GetValue(r.StreamId, key)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get blob info from kv node")
	}
	blobInfo, err := new(core.KVBlobInfo).Deserialize(val.Data)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to deserialize KV blob info")
	}
	return blobInfo, nil
}

func (r *retrievalClient) fetchChunksWithProof(dataRoot eth_common.Hash, blobInfo *core.KVBlobInfo, encodingParams core.EncodingParams, location *core.BlobLocation) ([]*core.Chunk, []core.ChunkNumber, error) {
	coeffLength := encodingParams.ChunkLength * core.CoeffSize

	requiredNum := (encodingParams.NumChunks + 1) / 2
	chunks := make([]*core.Chunk, 0)
	indices := make([]core.ChunkNumber, 0)

	for i := uint(0); i < encodingParams.NumChunks; i++ {
		// fetch i-th chunk coeffs
		offset := location.SegmentIndexes[i]*core.SegmentSize + location.Offsets[i]
		clientIndex := location.SegmentIndexes[i] % uint(len(r.Nodes))
		for j := 0; j < len(r.Nodes); j++ {
			data, err := r.Nodes[clientIndex].ZeroGStorage().DownloadSegment(
				dataRoot,
				uint64(offset/core.EntrySize),
				uint64((offset+coeffLength+core.ProofSize+core.EntrySize-1)/core.EntrySize),
			)
			if err != nil {
				r.logger.Debug("Download chunk error", "chunk_index", i, "dataRoot", dataRoot, "clientIndex", clientIndex)
				clientIndex = (clientIndex + 1) % uint(len(r.Nodes))
				continue
			}
			entryOffset := offset % core.EntrySize
			if len(data) >= int(entryOffset+coeffLength+core.ProofSize) {
				// read coeffs
				coeffs := core.BytesToCoeffs(data[entryOffset : entryOffset+coeffLength])
				// read proof
				var proof [64]byte
				copy(proof[:], data[entryOffset+coeffLength:entryOffset+coeffLength+core.ProofSize])
				// generate chunk
				chunk := &core.Chunk{
					Coeffs: coeffs,
					Proof:  core.BytesToProof(proof),
				}
				err = r.encoder.VerifyChunks([]*core.Chunk{chunk}, []core.ChunkNumber{i}, blobInfo.BlobHeader.BlobCommitments, encodingParams)
				if err != nil {
					r.logger.Debug("Validate chunk failed", "chunk_index", i, "clientIndex", clientIndex, "error", err)
				} else {
					// verified successfully
					r.logger.Info("verified chunk", "chunk index", i)
					chunks = append(chunks, chunk)
					indices = append(indices, i)
					break
				}
			}
			clientIndex = (clientIndex + 1) % uint(len(r.Nodes))
		}
		// check number of validated chunks
		if len(chunks) >= int(requiredNum) {
			return chunks, indices, nil
		}
	}
	return nil, nil, errors.New("failed to download enough verified chunks")
}

func (r *retrievalClient) RetrieveBlob(
	ctx context.Context,
	batchHeaderHash [32]byte,
	batchDataRoot eth_common.Hash,
	blobIndex uint32,
	referenceBlockNumber uint,
	batchRoot [32]byte,
	blobLengths []uint) ([]byte, error) {

	// Get blob header from any operator
	blobInfo, err := r.fetchBlobInfo(batchHeaderHash, blobIndex)
	if err != nil {
		return nil, err
	}

	// validate header merkle proof
	blobHeader := blobInfo.BlobHeader
	proof := blobInfo.MerkleProof
	blobHeaderHash, err := blobHeader.GetBlobHeaderHash()
	if err != nil {
		r.logger.Warn("got invalid blob header: ", err)
		return nil, err
	}
	proofVerified, err := merkletree.VerifyProofUsing(blobHeaderHash[:], false, proof, [][]byte{batchRoot[:]}, keccak256.New())
	if err != nil {
		r.logger.Warn("got invalid blob header proof: ", err)
		return nil, err
	}
	if !proofVerified {
		r.logger.Warn("failed to verify blob header against given proof")
		return nil, err
	}

	// Validate the blob length
	err = r.encoder.VerifyBlobLength(blobHeader.BlobCommitments)
	if err != nil {
		return nil, err
	}

	// encoding params
	chunkLength, chunkNum := core.SplitToChunks(blobHeader.Length)
	encodingParams, err := core.GetEncodingParams(chunkLength, chunkNum)
	if err != nil {
		return nil, err
	}

	r.logger.Debugf("encoding params: %v\n", encodingParams)

	// blob locations
	blobLocations := make([]*core.BlobLocation, len(blobLengths))
	for i, l := range blobLengths {
		chunkLength, chunkNum := core.SplitToChunks(l)
		blobLocations[i] = &core.BlobLocation{
			ChunkLength:    chunkLength,
			ChunkNum:       chunkNum,
			SegmentIndexes: make([]uint, chunkNum),
			Offsets:        make([]uint, chunkNum),
		}
	}
	core.AllocateChunks(blobLocations)

	// Fetch chunks from all zgs nodes
	chunks, indices, err := r.fetchChunksWithProof(batchDataRoot, blobInfo, encodingParams, blobLocations[blobIndex])
	if err != nil {
		return nil, err
	}

	return r.encoder.Decode(chunks, indices, encodingParams, uint64(blobHeader.Length)*bn254.BYTES_PER_COEFFICIENT)
}
