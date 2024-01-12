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
		batchRoot [32]byte) ([]byte, error)
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
	key, err := (&core.KVBlobInfoKey{
		BatchHeaderHash: batchHeaderHash,
		BlobIndex:       blobIndex,
	}).Serialize()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to serialize blob info key")
	}
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

func (r *retrievalClient) fetchProofAndValidate(
	dataRoot eth_common.Hash,
	encodingParams core.EncodingParams,
	blobHeader *core.BlobHeader,
	index core.ChunkNumber,
	proofOffset uint,
	coeffs []core.Symbol,
) (*core.Chunk, error) {
	clientIndex := proofOffset / core.SegmentSize % uint(len(r.Nodes))
	for j := 0; j < len(r.Nodes); j++ {
		data, err := r.Nodes[clientIndex].ZeroGStorage().DownloadSegment(
			dataRoot,
			uint64(proofOffset/core.EntrySize),
			uint64((proofOffset+core.ProofSize+core.EntrySize-1)/core.EntrySize),
		)
		if err != nil {
			r.logger.Info("Download segment error", "dataRoot", dataRoot, "clientIndex", clientIndex)
		}
		if uint(len(data)) == core.ProofSize {
			// read proof
			entryOffset := proofOffset % core.EntrySize
			var proof [64]byte
			copy(proof[:], data[entryOffset:entryOffset+core.ProofSize])
			// generate chunk
			chunk := &core.Chunk{
				Coeffs: coeffs,
				Proof:  core.BytesToProof(proof),
			}
			err = r.encoder.VerifyChunks([]*core.Chunk{chunk}, []core.ChunkNumber{index}, blobHeader.BlobCommitments, encodingParams)
			if err != nil {
				r.logger.Error("failed to verify chunk", "chunk index", index, "err", err)
			} else {
				// verified
				r.logger.Info("verified chunk", "chunk index", index)
				return chunk, nil
			}
		}
		clientIndex = (clientIndex + 1) % uint(len(r.Nodes))
	}
	return nil, errors.New("chunk verification failed")
}

func (r *retrievalClient) fetchChunksWithProof(dataRoot eth_common.Hash, blobInfo *core.KVBlobInfo, encodingParams core.EncodingParams) ([]*core.Chunk, []core.ChunkNumber, error) {
	coeffLength := encodingParams.ChunkLength * core.CoeffSize

	chunkOffset := blobInfo.ChunkOffset
	proofOffset := blobInfo.ProofOffset

	requiredNum := (encodingParams.NumChunks + 1) / 2
	chunks := make([]*core.Chunk, 0)
	indices := make([]core.ChunkNumber, 0)

	for i := uint(0); i < encodingParams.NumChunks; i++ {
		// fetch i-th chunk coeffs
		clientIndex := chunkOffset / core.SegmentSize % uint(len(r.Nodes))
		coeffs := make([]byte, 0)
		for j := 0; j < len(r.Nodes); j++ {
			data, err := r.Nodes[clientIndex].ZeroGStorage().DownloadSegment(
				dataRoot,
				uint64(chunkOffset/core.EntrySize),
				uint64((chunkOffset+coeffLength+core.EntrySize-1)/core.EntrySize),
			)
			if err != nil {
				r.logger.Info("Download segment error", "dataRoot", dataRoot, "clientIndex", clientIndex)
			}
			if uint(len(data)) == coeffLength {
				// read coeffs
				entryOffset := chunkOffset % core.EntrySize
				coeffs = append(coeffs, data[entryOffset:entryOffset+coeffLength]...)
				// fetch proof
				chunk, _ := r.fetchProofAndValidate(dataRoot, encodingParams, blobInfo.BlobHeader, i, proofOffset, core.BytesToCoeffs(coeffs))
				if chunk != nil {
					// verified successfully
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
		// move forward
		chunkOffset += coeffLength
		proofOffset += core.ProofSize
	}
	return nil, nil, errors.New("failed to download enough verified chunks")
}

func (r *retrievalClient) RetrieveBlob(
	ctx context.Context,
	batchHeaderHash [32]byte,
	batchDataRoot eth_common.Hash,
	blobIndex uint32,
	referenceBlockNumber uint,
	batchRoot [32]byte) ([]byte, error) {

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

	// Fetch chunks from all zgs nodes
	chunks, indices, err := r.fetchChunksWithProof(batchDataRoot, blobInfo, encodingParams)
	if err != nil {
		return nil, err
	}

	return r.encoder.Decode(chunks, indices, encodingParams, uint64(blobHeader.Length)*bn254.BYTES_PER_COEFFICIENT)
}
