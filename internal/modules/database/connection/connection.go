package connection

import (
	"context"
)

type DatabaseDriver string

const (
	SQLiteDriver     DatabaseDriver = "sqlite3"
	MySQLDriver      DatabaseDriver = "mysql"
	PostgreSQLDriver DatabaseDriver = "postgres"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Schema   string
}

type DatabaseConnection interface {
	Prepared(ctx context.Context, query string, binding ...any) error
	Unprepared(ctx context.Context, query string) error
	Select(ctx context.Context, dest any, query string, args ...any) error
	SelectOne(ctx context.Context, dest any, query string, args ...any) error
	Scan(ctx context.Context, query string, args []any, dest ...any) error
	Cursor(ctx context.Context, query string, args []any) (Cursor, error)
	GetDriverName() DatabaseDriver
	From(table string) QueryBuilder
	Close() error
}
