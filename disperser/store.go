package disperser

import (
	"context"
	"errors"

	"github.com/0glabs/0g-data-avail/common"
	"github.com/0glabs/0g-data-avail/disperser/leveldb"
)

// Store is a key-value database to store blob data (blob header, blob chunks etc).
type Store struct {
	db     DB
	logger common.Logger
}

func NewLevelDBStore(path string, logger common.Logger) (*Store, error) {
	// Create the db at the path. This is currently hardcoded to use
	// levelDB.
	db, err := leveldb.NewLevelDBStore(path)
	if err != nil {
		logger.Error("Could not create leveldb database", "err", err)
		return nil, err
	}

	return &Store{
		db:     db,
		logger: logger,
	}, nil
}

func (s *Store) StoreMetadata(ctx context.Context, key []byte, value []byte) error {
	err := s.db.Put(key, value)
	if err != nil {
		s.logger.Error("Failed to write the batch into local database:", "err", err)
		return err
	}
	s.logger.Debug("StoreMetadata succeeded")
	return nil
}

func (s *Store) GetMetadata(ctx context.Context, key []byte) ([]byte, error) {
	data, err := s.db.Get(key)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	return data, nil
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
