package disperser

import "errors"

var (
	ErrBlobNotFound   = errors.New("blob not found")
	ErrMemoryDbIsFull = errors.New("memory db is full")
)
