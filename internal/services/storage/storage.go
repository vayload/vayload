package storage

import (
	"context"
	"io"

	"github.com/vayload/vayload/internal/kernel"
	"github.com/vayload/vayload/internal/vayload"
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
	Put(ctx context.Context, key string, reader io.Reader) (PutResult, error)

	Get(ctx context.Context, key string) (io.ReadCloser, GetInfo, error)

	Delete(ctx context.Context, key string) error
	Stat(ctx context.Context, key string) (ObjectInfo, error)
}

type StorageService struct {
	kernel.BaseService
	store ObjectStore
}

func NewStorageService(store ObjectStore) *StorageService {
	deps := []vayload.ServiceName{
		vayload.ServiceAuthName,
		vayload.ServiceDatabaseName,
	}

	return &StorageService{
		BaseService: kernel.NewBaseService(string(vayload.ServiceStorageName), false, deps...),
		store:       store,
	}
}

func (s *StorageService) Bootstrap(ctx context.Context, args map[string]any, reply *map[string]any) error {
	return nil
}

func (s *StorageService) Shutdown() {
	// free resources
}

func (s *StorageService) HttpRoutes() []vayload.HttpRoute {
	return []vayload.HttpRoute{}
}
