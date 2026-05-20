package connection

import "context"

type StatementType int

const (
	StmtSelect StatementType = iota
	StmtInsert
	StmtUpdate
	StmtDelete
)

type Cursor interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}

type QueryBuilder interface {
	// base
	From(table string) QueryBuilder
	FromSub(q QueryBuilder, alias string) QueryBuilder
	Select(columns ...string) QueryBuilder
	Distinct() QueryBuilder

	// joins
	Join(table string, on string, args ...any) QueryBuilder
	JoinSub(sub QueryBuilder, alias string, on string, args ...any) QueryBuilder
	LeftJoin(table string, on string, args ...any) QueryBuilder
	LeftJoinSub(sub QueryBuilder, alias string, on string, args ...any) QueryBuilder
	RightJoin(table string, on string, args ...any) QueryBuilder
	RightJoinSub(sub QueryBuilder, alias string, on string, args ...any) QueryBuilder
	InnerJoin(table string, on string, args ...any) QueryBuilder
	InnerJoinSub(sub QueryBuilder, alias string, on string, args ...any) QueryBuilder

	// where
	Where(column string, operator string, value any) QueryBuilder
	Wheres(filters map[string]any) QueryBuilder
	OrWhere(column string, operator string, value any) QueryBuilder
	WhereIn(column string, values ...any) QueryBuilder
	WhereColumn(col1 string, operator string, col2 string) QueryBuilder
	WhereBetween(column string, min any, max any) QueryBuilder
	WhereNotIn(column string, values ...any) QueryBuilder
	WhereNull(column string) QueryBuilder
	WhereNotNull(column string) QueryBuilder

	Cursor(ctx context.Context) (Cursor, error)

	// mutation
	InsertOne(values map[string]any) QueryBuilder
	InsertMany(values []map[string]any) QueryBuilder
	UpdateOne(values map[string]any) QueryBuilder
	UpsertOne(values map[string]any, onConflict []string) QueryBuilder
	Delete() QueryBuilder

	// group
	GroupBy(columns ...string) QueryBuilder
	Having(condition string, args ...any) QueryBuilder

	// order / paging
	OrderBy(column string, direction string) QueryBuilder
	Limit(limit int64) QueryBuilder
	Offset(offset int64) QueryBuilder
	Take(n int) QueryBuilder

	Seek(column, op string, value any) QueryBuilder

	// execution
	Get(ctx context.Context, dest any) error
	First(ctx context.Context, dest any) error
	Scan(ctx context.Context, dest ...any) error
	Count(ctx context.Context) (int64, error)
	Exists(ctx context.Context) (bool, error)

	// Find by id
	Find(ctx context.Context, id string, dest any) error

	Union(q QueryBuilder) QueryBuilder
	UnionAll(q QueryBuilder) QueryBuilder

	// locking
	ForUpdate() QueryBuilder
	ForShare() QueryBuilder

	// Execute performs INSERT, UPDATE, or DELETE operations
	Exec(ctx context.Context) error

	// final sql
	ToSql() (string, []any)
}
