package disperser

import "errors"

var (
	ErrBlobNotFound   = errors.New("blob not found")
	ErrMemoryDbIsFull = errors.New("memory db is full")
	ErrKeyNotFound    = errors.New("key not found in db")
)
