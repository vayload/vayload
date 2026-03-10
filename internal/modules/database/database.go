package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/kernel"
	"github.com/vayload/vayload/internal/modules/database/connection"
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
		BaseService: kernel.NewBaseService(vayload.ServiceDatabaseName, true),
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
	container.SetInstance(DATABASE_CONNECTION, conn)

	// register schema builder as singleton deferred
	container.Deferred(SCHEMA_BUILDER, func(container vayload.Container) (any, error) {
		return Factory.CreateSchemaBuilder(ctx, driverName, conn)
	}, true)

	return nil
}

func (s *DatabaseService) Shutdown(ctx context.Context) error {
	container := s.Container()
	if container == nil {
		return fmt.Errorf("container not provided for database service")
	}

	// free resources close database connection only if exists
	if container.Has(DATABASE_CONNECTION) {
		var instance connection.DatabaseConnection
		if err := container.GetInto(DATABASE_CONNECTION, &instance); err != nil {
			return err
		}

		instance.Close()
	}

	return nil
}

func (s *DatabaseService) HttpRoutes() []vayload.HttpRoutesGroup {
	return []vayload.HttpRoutesGroup{
		{
			Prefix: "/database",
			Routes: []vayload.HttpRoute{
				{
					Path:   "/files/upload",
					Method: vayload.HttpPost,
					Handler: func(req vayload.HttpRequest, res vayload.HttpResponse) error {
						return res.Send([]byte("Hello "))
					},
				},
			},
		},
	}
}

var _ vayload.Service = (*DatabaseService)(nil)
