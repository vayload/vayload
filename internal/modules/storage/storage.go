package storage

import (
	"context"
	"path/filepath"

	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/kernel"
	"github.com/vayload/vayload/internal/modules/database"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/modules/storage/domain"
	"github.com/vayload/vayload/internal/modules/storage/engines"
	"github.com/vayload/vayload/internal/modules/storage/handlers"
	"github.com/vayload/vayload/internal/modules/storage/repositories"
	"github.com/vayload/vayload/internal/modules/storage/services"
	"github.com/vayload/vayload/internal/vayload"
	"github.com/vayload/vayload/pkg/crypto"
)

type StorageService struct {
	kernel.BaseService
	store      engines.ObjectStore
	db         connection.DatabaseConnection
	crypto     crypto.Encryption
	repository domain.StorageRepository
	service    *services.StorageService

	// Concrete handlers
	httpHandler *handlers.StorageHttpHandler
}

func NewStorageService(config *config.Config) *StorageService {
	deps := []vayload.ServiceName{
		vayload.ServiceAuthName,
		vayload.ServiceDatabaseName,
	}

	return &StorageService{
		BaseService: kernel.NewBaseService(vayload.ServiceStorageName, false, deps...),
		store:       engines.NewManager(filepath.Join(config.DataDir, "storage")),
	}
}

func (s *StorageService) Bootstrap(ctx context.Context, args map[string]any, reply *map[string]any) error {
	var err error
	s.db, err = kernel.MapTo[connection.DatabaseConnection](s.Container(), database.DATABASE_CONNECTION)
	if err != nil {
		return err
	}

	// Start initial services
	s.repository = repositories.NewSQLStorageRepository(s.db)
	s.service = services.NewStorageService(s.repository)
	s.crypto = crypto.NewEncryption("vayload-storage-secret")

	// Create handlers
	s.httpHandler = handlers.NewStorageHttpHandler(s.service)

	return nil
}

func (s *StorageService) Shutdown(ctx context.Context) error {
	// free resources
	return nil
}

func (s *StorageService) HttpRoutes() []vayload.HttpRoutesGroup {
	if s.httpHandler == nil {
		s.httpHandler = handlers.NewStorageHttpHandler(s.service)
	}

	return s.httpHandler.HttpRoutes()
}

var _ vayload.Service = (*StorageService)(nil)
