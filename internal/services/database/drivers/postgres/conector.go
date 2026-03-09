package postgres

import (
	"context"
	"fmt"

	"github.com/vayload/vayload/internal/services/database/connection"
	"github.com/vayload/vayload/internal/services/database/drivers"

	// postgres driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, schema)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	return &databaseConnection{
		BaseConnection: drivers.NewBaseConnection(db, ctx, cancel, connection.PostgreSQLDriver),
		User:           user,
		Password:       password,
		Host:           host,
		Port:           port,
		Schema:         schema,
	}, nil
}

var _ connection.DatabaseConnection = (*databaseConnection)(nil)
