package disperser

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"time"

	"github.com/0glabs/0g-da-client/common"
	"github.com/0glabs/0g-da-client/disperser/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

const (
	// How many blobs to delete in one atomic operation during the expiration
	// garbage collection.
	numBlobsToDeleteAtomically = 8
)

// Store is a key-value database to store blob data (blob header, blob chunks etc).
type Store struct {
	db           DB
	timeToExpire uint
	logger       common.Logger
}

func NewLevelDBStore(path string, timeToExpire uint, logger common.Logger) (*Store, error) {
	// Create the db at the path. This is currently hardcoded to use
	// levelDB.
	db, err := leveldb.NewLevelDBStore(path)
	if err != nil {
		logger.Error("Could not create leveldb database", "err", err)
		return nil, err
	}

	return &Store{
		db:           db,
		timeToExpire: timeToExpire,
		logger:       logger,
	}, nil
}

func (s *Store) DeleteExpiredEntries(currentTimeUnixSec int64, timeLimitSec uint64) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeLimitSec)*time.Second)
	defer cancel()

	numBlobsDeleted := 0
	for {
		select {
		case <-ctx.Done():
			return numBlobsDeleted, ctx.Err()
		default:
			numDeleted, err := s.deleteNBlobs(currentTimeUnixSec, numBlobsToDeleteAtomically)
			if err != nil {
				return numBlobsDeleted, err
			}
			// When there is no error and we didn't delete any batch, it means we have
			// no obsolete batches to delete, so we can return.
			if numDeleted == 0 {
				return numBlobsDeleted, nil
			}
			numBlobsDeleted += numDeleted
		}
	}
}

func (s *Store) deleteNBlobs(currentTimeUnixSec int64, numBatches int) (int, error) {
	// Scan for expired batches.
	iter := s.db.NewIterator(EncodeBatchExpirationKeyPrefix())
	expiredKeys := make([][]byte, 0)
	expiredBlobs := make([][]byte, 0)
	for iter.Next() {
		ts, err := DecodeBatchExpirationKey(iter.Key())
		if err != nil {
			s.logger.Error("Could not decode the expiration key", "key:", iter.Key(), "error:", err)
			continue
		}
		// No more rows expired up to current time.
		if currentTimeUnixSec < ts {
			break
		}
		expiredKeys = append(expiredKeys, copyBytes(iter.Key()))
		expiredBlobs = append(expiredBlobs, copyBytes(iter.Value()))
		if len(expiredKeys) == numBatches {
			break
		}
	}
	iter.Release()

	// No expired batch found.
	if len(expiredKeys) == 0 {
		return 0, nil
	}

	for _, blobKey := range expiredBlobs {
		buf := bytes.NewReader(blobKey)

		for {
			var length uint64
			err := binary.Read(buf, binary.LittleEndian, &length)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return -1, err
			}

			chunk := make([]byte, length)
			_, err = buf.Read(chunk)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return -1, err
			}

			blobHeaderKey, err := EncodeBlobHeaderKey(chunk)
			if err != nil {
				s.logger.Error("Cannot generate the key for query blob:", "err", err)
				return -1, err
			}

			expiredKeys = append(expiredKeys, blobHeaderKey)
		}
	}

	// Perform the removal.
	err := s.db.DeleteBatch(expiredKeys)
	if err != nil {
		s.logger.Error("Failed to delete the expired keys in batch", "keys:", expiredKeys, "error:", err)
		return -1, err
	}

	return len(expiredBlobs), nil
}

func (s *Store) StoreMetadata(ctx context.Context, key []byte, value []byte) error {
	blobHeaderKey, err := EncodeBlobHeaderKey(key)
	if err != nil {
		s.logger.Error("Cannot generate the key for storing blob:", "err", err)
		return err
	}

	err = s.db.Put(blobHeaderKey, value)
	if err != nil {
		s.logger.Error("Failed to write the batch into local database:", "err", err)
		return err
	}
	s.logger.Debug("StoreMetadata succeeded")
	return nil
}

func (s *Store) StoreMetadataBatch(ctx context.Context, blobKeys [][]byte, metadatas [][]byte) (*[][]byte, error) {
	keys := make([][]byte, 0)
	values := make([][]byte, 0)

	buf := bytes.NewBuffer(make([]byte, 0))

	for idx, key := range blobKeys {
		blobHeaderKey, err := EncodeBlobHeaderKey(key)
		if err != nil {
			s.logger.Error("Cannot generate the key for storing blob:", "err", err)
			return nil, err
		}

		keys = append(keys, blobHeaderKey)
		values = append(values, metadatas[idx])

		if err := binary.Write(buf, binary.LittleEndian, uint64(len(key))); err != nil {
			return nil, err
		}
		if _, err := buf.Write(key); err != nil {
			return nil, err
		}
	}

	curr := time.Now().Unix()
	// timeToExpire := 60 * 24 * 60 * 60
	expirationTime := curr + int64(s.timeToExpire)
	expirationKey := EncodeBatchExpirationKey(expirationTime)
	keys = append(keys, expirationKey)
	values = append(values, buf.Bytes())

	start := time.Now()

	err := s.db.WriteBatch(keys, values)
	if err != nil {
		s.logger.Error("Failed to write the batch into local database:", "err", err)
		return nil, err
	}
	s.logger.Debug("StoreMetaBatch succeeded", "write batch duration", time.Since(start), "size", len(keys))

	return &keys, nil
}

func (s *Store) GetMetadata(ctx context.Context, key []byte) ([]byte, error) {
	blobHeaderKey, err := EncodeBlobHeaderKey(key)
	if err != nil {
		s.logger.Error("Cannot generate the key for get blob:", "err", err)
		return nil, err
	}

	data, err := s.db.Get(blobHeaderKey)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	return data, nil
}

func (s *Store) MetadataIterator(ctx context.Context) iterator.Iterator {
	return s.db.NewIterator(EncodeBlobHeaderKeyPrefix())
}

// HasKey returns if a given key has been stored.
func (s *Store) HasKey(ctx context.Context, key []byte) bool {
	_, err := s.db.Get(key)
	return err == nil
}

// DeleteKeys removes a list of keys from the store atomically.
//
// Note: caller should ensure these keys are exactly all the data items for a single batch
// to maintain the integrity of the store.
func (s *Store) DeleteKeys(ctx context.Context, keys *[][]byte) error {
	return s.db.DeleteBatch(*keys)
}

const (
	// Caution: the change to these prefixes needs to handle the backward compatibility,
	// making sure the new code work with old data in DA Node store.
	blobHeaderPrefix      = "_BLOB_HEADER_" // The prefix of the blob header key.
	batchExpirationPrefix = "_EXPIRATION_"  // The prefix of the batch expiration key.
)

func EncodeBatchExpirationKey(expirationTime int64) []byte {
	prefix := []byte(batchExpirationPrefix)
	ts := make([]byte, 8)
	binary.BigEndian.PutUint64(ts[0:8], uint64(expirationTime))
	buf := bytes.NewBuffer(append(prefix, ts[:]...))
	return buf.Bytes()
}

// Returns the encoded prefix for batch expiration key.
func EncodeBatchExpirationKeyPrefix() []byte {
	return []byte(batchExpirationPrefix)
}

// Returns the expiration timestamp encoded in the key.
func DecodeBatchExpirationKey(key []byte) (int64, error) {
	if len(key) != len(batchExpirationPrefix)+8 {
		return 0, errors.New("the expiration key is invalid")
	}
	ts := int64(binary.BigEndian.Uint64(key[len(key)-8:]))
	return ts, nil
}

// EncodeBlobHeaderKey returns an encoded key as blob header identification.
func EncodeBlobHeaderKey(key []byte) ([]byte, error) {
	prefix := []byte(blobHeaderPrefix)
	buf := bytes.NewBuffer(append(prefix, key[:]...))
	return buf.Bytes(), nil
}

// Returns an encoded prefix of blob header key.
func EncodeBlobHeaderKeyPrefix() []byte {
	return []byte(blobHeaderPrefix)
}

func copyBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
