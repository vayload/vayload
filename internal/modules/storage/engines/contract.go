package engines

import (
	"context"
	"io"
)

type PutResult struct {
	Size   int64
	SHA256 [32]byte
}

type ObjectInfo struct {
	Size   int64
	SHA256 [32]byte
}

type GetInfo = ObjectInfo

type ObjectStore interface {
	Put(ctx context.Context, bucket string, key string, reader io.Reader) (PutResult, error)

	Get(ctx context.Context, bucket string, key string) (io.ReadCloser, GetInfo, error)

	Delete(ctx context.Context, bucket string, key string) error
	Stat(ctx context.Context, bucket string, key string) (ObjectInfo, error)
}
