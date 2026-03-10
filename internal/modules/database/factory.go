package database

import (
	"context"
	"fmt"

	"github.com/vayload/vayload/internal/modules/database/builder"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/modules/database/drivers/mysql"
	"github.com/vayload/vayload/internal/modules/database/drivers/postgres"
	"github.com/vayload/vayload/internal/modules/database/drivers/sqlite3"
)

type ConnectionFactory interface {
	CreateConnection(ctx context.Context, driver connection.DatabaseDriver, config connection.Config) (connection.DatabaseConnection, error)
	CreateSchemaBuilder(ctx context.Context, driver connection.DatabaseDriver, conn connection.DatabaseConnection) (*builder.SchemaBuilder, error)
}

type connectionFactory struct {
}

func (c *connectionFactory) CreateConnection(ctx context.Context, driver connection.DatabaseDriver, config connection.Config) (connection.DatabaseConnection, error) {
	switch driver {
	case connection.PostgreSQLDriver:
		return postgres.NewConnection(ctx, config.User, config.Password, config.Host, config.Port, config.Schema)
	case connection.MySQLDriver:
		return mysql.NewConnection(ctx, config.User, config.Password, config.Host, config.Port, config.Schema)
	case connection.SQLiteDriver:
		return sqlite3.NewConnection(ctx, config.User, config.Password, config.Host, config.Port, config.Schema)
	default:
		return nil, fmt.Errorf("driver %s currently not supported", driver)
	}
}

func (c *connectionFactory) CreateSchemaBuilder(ctx context.Context, driver connection.DatabaseDriver, conn connection.DatabaseConnection) (*builder.SchemaBuilder, error) {
	switch driver {
	case connection.PostgreSQLDriver:
		return builder.NewSchemaBuilder(postgres.NewGrammar(), conn), nil
	case connection.MySQLDriver:
		return builder.NewSchemaBuilder(mysql.NewGrammar(), conn), nil
	case connection.SQLiteDriver:
		return builder.NewSchemaBuilder(sqlite3.NewGrammar(), conn), nil
	default:
		return nil, fmt.Errorf("driver %s currently not supported", driver)
	}
}

var Factory ConnectionFactory

func init() {
	Factory = &connectionFactory{}
}
