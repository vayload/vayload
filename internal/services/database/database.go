package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/kernel"
	"github.com/vayload/vayload/internal/services/database/connection"
	"github.com/vayload/vayload/internal/vayload"
)

const (
	DATABASE_NAME       = "database.name"
	DATABASE_CONNECTION = "database.connection"
	SCHEMA_BUILDER      = "database.schema.builder"
)

type DatabaseService struct {
	kernel.BaseService

	config *config.Config
}

func NewDatabaseService(config *config.Config) *DatabaseService {
	return &DatabaseService{
		BaseService: kernel.NewBaseService(string(vayload.ServiceDatabaseName), true),
		config:      config,
	}
}

func (s *DatabaseService) Bootstrap(ctx context.Context, args map[string]any, reply *map[string]any) error {
	driverName := connection.DatabaseDriver(s.config.Database.Driver)
	conn, err := Factory.CreateConnection(ctx, driverName, connection.Config{
		User:     s.config.Database.User,
		Password: s.config.Database.Password,
		Host:     s.config.Database.Host,
		Port:     strconv.Itoa(s.config.Database.Port),
		Schema:   s.config.Database.Schema,
	})

	if err != nil {
		return err
	}

	container := s.Container()
	if container == nil {
		return fmt.Errorf("container not provided for database service")
	}

	// register database connection as singleton
	container.Singleton(DATABASE_CONNECTION, conn)

	// register schema builder as singleton deferred
	container.Deffered(SCHEMA_BUILDER, func(container vayload.Container) (any, error) {
		return Factory.CreateSchemaBuilder(ctx, driverName, conn)
	}, true)

	return nil
}

func (s *DatabaseService) Shutdown() {
	container := s.Container()
	if container == nil {
		return
	}

	// free resources close database connection only if exists
	if container.Has(DATABASE_CONNECTION) {
		var instance connection.DatabaseConnection
		if err := container.ResolveInto(DATABASE_CONNECTION, &instance); err != nil {
			return
		}

		instance.Close()
	}
}

var _ vayload.Service = (*DatabaseService)(nil)
