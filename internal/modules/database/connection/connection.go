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
	GetDriverName() DatabaseDriver
	From(table string) QueryBuilder
	Close() error
}

type StatementType int

const (
	StmtSelect StatementType = iota
	StmtInsert
	StmtUpdate
	StmtDelete
)

type QueryBuilder interface {
	// base
	From(table string) QueryBuilder
	Select(columns ...string) QueryBuilder
	Distinct() QueryBuilder

	// joins
	Join(table string, on string, args ...any) QueryBuilder
	LeftJoin(table string, on string, args ...any) QueryBuilder
	RightJoin(table string, on string, args ...any) QueryBuilder
	InnerJoin(table string, on string, args ...any) QueryBuilder

	// where
	Where(column string, operator string, value any) QueryBuilder
	Wheres(filters map[string]any) QueryBuilder
	OrWhere(column string, operator string, value any) QueryBuilder
	WhereIn(column string, values ...any) QueryBuilder
	WhereNotIn(column string, values ...any) QueryBuilder
	WhereNull(column string) QueryBuilder
	WhereNotNull(column string) QueryBuilder

	// mutation
	Insert(values map[string]any) QueryBuilder
	Update(values map[string]any) QueryBuilder
	Upsert(values map[string]any, onConflict []string) QueryBuilder
	Delete() QueryBuilder

	// group
	GroupBy(columns ...string) QueryBuilder
	Having(condition string, args ...any) QueryBuilder

	// order / paging
	OrderBy(column string, direction string) QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder

	// execution
	Get(dest any) error
	GetAll(dest any) error
	Scan(dest ...any) error
	Count() (int64, error)
	Exists() (bool, error)

	// Execute performs INSERT, UPDATE, or DELETE operations
	Exec() error

	// final sql
	ToSql() (string, []any)
}
