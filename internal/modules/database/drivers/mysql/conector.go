package mysql

import (
	"context"
	"fmt"

	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/modules/database/drivers"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, schema)

	db, err := sqlx.ConnectContext(ctx, "mysql", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	return &databaseConnection{
		BaseConnection: drivers.NewBaseConnection(db, ctx, cancel, connection.MySQLDriver),
		User:           user,
		Password:       password,
		Host:           host,
		Port:           port,
		Schema:         schema,
	}, nil
}

var _ connection.DatabaseConnection = (*databaseConnection)(nil)
