package sqlite3

import (
	"context"
	"fmt"

	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/modules/database/drivers"

	// sqlite3 driver
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type databaseConnection struct {
	drivers.BaseConnection

	User     string
	Password string
	Host     string
	Port     string
	Schema   string
}

func NewConnection(ctx context.Context, user, password, host, port, schema string) (*databaseConnection, error) {
	dsn := fmt.Sprintf("file:./data/%s.db?_journal_mode=WAL&_cache_size=10000&_foreign_keys=on", schema)

	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	return &databaseConnection{
		BaseConnection: drivers.NewBaseConnection(db, ctx, cancel, connection.SQLiteDriver),
		User:           user,
		Password:       password,
		Host:           host,
		Port:           port,
		Schema:         schema,
	}, nil
}

var _ connection.DatabaseConnection = (*databaseConnection)(nil)
