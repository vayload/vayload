package builder

import (
	"context"
	"fmt"
	"strings"

	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/modules/database/grammar"
	"github.com/vayload/vayload/internal/modules/database/query"
)

type QueryBuilder struct {
	conn    connection.DatabaseConnection
	driver  connection.DatabaseDriver
	grammar grammar.QueryGrammar

	query *query.Query
}

func NewQueryBuilder(conn connection.DatabaseConnection, grammar grammar.QueryGrammar, table string) connection.QueryBuilder {
	return &QueryBuilder{
		conn:    conn,
		grammar: grammar,
		driver:  conn.GetDriverName(),
		query: &query.Query{
			StmtType: query.StmtSelect,
			Columns:  []string{"*"},
			Table:    table,
		},
	}
}

func (q *QueryBuilder) From(table string) connection.QueryBuilder {
	q.query.Table = table
	return q
}

func (q *QueryBuilder) FromSub(sub connection.QueryBuilder, alias string) connection.QueryBuilder {
	q.query.Table = ""
	q.query.SubQuery = q.extractQuery(sub)
	q.query.Alias = alias
	return q
}

func (q *QueryBuilder) Select(columns ...string) connection.QueryBuilder {
	q.query.Columns = columns
	return q
}

func (q *QueryBuilder) Distinct() connection.QueryBuilder {
	q.query.Distinct = true
	return q
}

// Joins
func (q *QueryBuilder) Join(table string, on string, args ...any) connection.QueryBuilder {
	q.query.Joins = append(q.query.Joins, query.Join{
		Type:  query.InnerJoinType,
		Table: table,
		On:    on,
		Args:  args,
	})

	return q
}

func (q *QueryBuilder) JoinSub(sub connection.QueryBuilder, alias string, on string, args ...any) connection.QueryBuilder {
	q.query.Joins = append(q.query.Joins, query.Join{
		Type:     query.InnerJoinType,
		SubQuery: q.extractQuery(sub),
		Alias:    alias,
		On:       on,
		Args:     args,
	})
	return q
}

func (q *QueryBuilder) LeftJoin(table string, on string, args ...any) connection.QueryBuilder {
	q.query.Joins = append(q.query.Joins, query.Join{
		Type:  query.LeftJoinType,
		Table: table,
		On:    on,
		Args:  args,
	})

	return q
}

func (q *QueryBuilder) LeftJoinSub(sub connection.QueryBuilder, alias string, on string, args ...any) connection.QueryBuilder {
	q.query.Joins = append(q.query.Joins, query.Join{
		Type:     query.LeftJoinType,
		SubQuery: q.extractQuery(sub),
		Alias:    alias,
		On:       on,
		Args:     args,
	})
	return q
}

func (q *QueryBuilder) RightJoin(table string, on string, args ...any) connection.QueryBuilder {
	q.query.Joins = append(q.query.Joins, query.Join{
		Type:  query.RightJoinType,
		Table: table,
		On:    on,
		Args:  args,
	})

	return q
}

func (q *QueryBuilder) RightJoinSub(sub connection.QueryBuilder, alias string, on string, args ...any) connection.QueryBuilder {
	q.query.Joins = append(q.query.Joins, query.Join{
		Type:     query.RightJoinType,
		SubQuery: q.extractQuery(sub),
		Alias:    alias,
		On:       on,
		Args:     args,
	})
	return q
}

func (q *QueryBuilder) InnerJoin(table string, on string, args ...any) connection.QueryBuilder {
	q.query.Joins = append(q.query.Joins, query.Join{
		Type:  query.InnerJoinType,
		Table: table,
		On:    on,
		Args:  args,
	})

	return q
}

func (q *QueryBuilder) InnerJoinSub(sub connection.QueryBuilder, alias string, on string, args ...any) connection.QueryBuilder {
	q.query.Joins = append(q.query.Joins, query.Join{
		Type:     query.InnerJoinType,
		SubQuery: q.extractQuery(sub),
		Alias:    alias,
		On:       on,
		Args:     args,
	})
	return q
}

// Where clauses
func (q *QueryBuilder) Where(column string, operator string, value any) connection.QueryBuilder {
	w := query.Where{
		Type:     query.WhereTypeBasic,
		Column:   column,
		Operator: operator,
		Value:    value,
		IsOr:     false,
	}

	if sub, ok := value.(connection.QueryBuilder); ok {
		w.Type = query.WhereTypeSubQuery
		w.SubQuery = q.extractQuery(sub)
		w.Value = nil
	}

	q.query.Where = append(q.query.Where, w)
	return q
}

func (q *QueryBuilder) Wheres(filters map[string]any) connection.QueryBuilder {
	for column, value := range filters {
		q.Where(column, "=", value)
	}

	return q
}

func (q *QueryBuilder) OrWhere(column string, operator string, value any) connection.QueryBuilder {
	w := query.Where{
		Type:     query.WhereTypeBasic,
		Column:   column,
		Operator: operator,
		Value:    value,
		IsOr:     true,
	}

	if sub, ok := value.(connection.QueryBuilder); ok {
		w.Type = query.WhereTypeSubQuery
		w.SubQuery = q.extractQuery(sub)
		w.Value = nil
	}

	q.query.Where = append(q.query.Where, w)
	return q
}

func (q *QueryBuilder) WhereIn(column string, values ...any) connection.QueryBuilder {
	q.query.Where = append(q.query.Where, query.Where{
		Type:     query.WhereTypeIn,
		Column:   column,
		Operator: "IN",
		Value:    values,
		IsOr:     false,
	})

	return q
}

func (q *QueryBuilder) WhereNotIn(column string, values ...any) connection.QueryBuilder {
	q.query.Where = append(q.query.Where, query.Where{
		Type:     query.WhereTypeNotIn,
		Column:   column,
		Operator: "NOT IN",
		Value:    values,
		IsOr:     false,
	})

	return q
}

func (q *QueryBuilder) WhereColumn(col1 string, operator string, col2 string) connection.QueryBuilder {
	q.query.Where = append(q.query.Where, query.Where{
		Type:     query.WhereTypeColumn,
		Column:   col1,
		Operator: operator,
		Value:    col2,
		IsOr:     false,
	})
	return q
}

func (q *QueryBuilder) WhereBetween(column string, min any, max any) connection.QueryBuilder {
	q.query.Where = append(q.query.Where, query.Where{
		Type:     query.WhereTypeBetween,
		Column:   column,
		Operator: "BETWEEN",
		Value:    min,
		Value2:   max,
		IsOr:     false,
	})
	return q
}

func (q *QueryBuilder) WhereNull(column string) connection.QueryBuilder {
	q.query.Where = append(q.query.Where, query.Where{
		Type:     query.WhereTypeNull,
		Column:   column,
		Operator: "IS NULL",
		IsOr:     false,
	})

	return q
}

func (q *QueryBuilder) WhereNotNull(column string) connection.QueryBuilder {
	q.query.Where = append(q.query.Where, query.Where{
		Type:     query.WhereTypeNotNull,
		Column:   column,
		Operator: "IS NOT NULL",
		IsOr:     false,
	})

	return q
}

func (q *QueryBuilder) Cursor(ctx context.Context) (connection.Cursor, error) {
	query, args := q.ToSql()

	return q.conn.Cursor(ctx, query, args)
}

// Mutation methods
func (q *QueryBuilder) InsertOne(values map[string]any) connection.QueryBuilder {
	q.query.StmtType = query.StmtInsert
	q.query.InsertValues = values
	return q
}

func (q *QueryBuilder) InsertMany(values []map[string]any) connection.QueryBuilder {
	q.query.StmtType = query.StmtInsert
	q.query.InsertMultiValues = values
	return q
}

func (q *QueryBuilder) UpdateOne(values map[string]any) connection.QueryBuilder {
	q.query.StmtType = query.StmtUpdate
	q.query.UpdateValues = values
	return q
}

func (q *QueryBuilder) UpsertOne(values map[string]any, onConflict []string) connection.QueryBuilder {
	q.query.StmtType = query.StmtUpsert
	q.query.UpsertValues = values
	q.query.UpsertColumns = onConflict
	return q
}

func (q *QueryBuilder) Delete() connection.QueryBuilder {
	q.query.StmtType = query.StmtDelete
	return q
}

// Group and Having
func (q *QueryBuilder) GroupBy(columns ...string) connection.QueryBuilder {
	q.query.GroupBy = columns
	return q
}

func (q *QueryBuilder) Having(condition string, args ...any) connection.QueryBuilder {
	q.query.Having = append(q.query.Having, query.Having{
		Condition: condition,
		Args:      args,
	})

	return q
}

// Order and Pagination
func (q *QueryBuilder) OrderBy(column string, direction string) connection.QueryBuilder {
	q.query.OrderBy = append(q.query.OrderBy, query.Order{
		Column:    column,
		Direction: strings.ToUpper(direction),
	})

	return q
}

func (q *QueryBuilder) Limit(limit int64) connection.QueryBuilder {
	q.query.Limit = &limit
	return q
}

func (q *QueryBuilder) Offset(offset int64) connection.QueryBuilder {
	q.query.Offset = &offset
	return q
}

func (q *QueryBuilder) Take(n int) connection.QueryBuilder {
	limit := int64(n)
	q.query.Limit = &limit
	return q
}

func (q *QueryBuilder) Seek(column, op string, value any) connection.QueryBuilder {
	q.query.SeekColumn = column
	q.query.SeekOperator = op
	q.query.SeekValue = value
	return q
}

func (q *QueryBuilder) Union(sub connection.QueryBuilder) connection.QueryBuilder {
	q.query.Unions = append(q.query.Unions, query.Union{
		Query: q.extractQuery(sub),
		All:   false,
	})
	return q
}

func (q *QueryBuilder) UnionAll(sub connection.QueryBuilder) connection.QueryBuilder {
	q.query.Unions = append(q.query.Unions, query.Union{
		Query: q.extractQuery(sub),
		All:   true,
	})
	return q
}

func (q *QueryBuilder) ForUpdate() connection.QueryBuilder {
	q.query.LockMode = query.LockUpdate
	return q
}

func (q *QueryBuilder) ForShare() connection.QueryBuilder {
	q.query.LockMode = query.LockShare
	return q
}

// Execution methods
func (q *QueryBuilder) Get(ctx context.Context, dest any) error {
	if dest == nil {
		return fmt.Errorf("query builder: dest cannot be nil")
	}

	query, args := q.ToSql()

	return q.conn.Select(ctx, dest, query, args...)
}

func (q *QueryBuilder) First(ctx context.Context, dest any) error {
	if dest == nil {
		return fmt.Errorf("query builder: dest cannot be nil")
	}

	query, args := q.ToSql()

	return q.conn.SelectOne(ctx, dest, query, args...)
}

func (q *QueryBuilder) Find(ctx context.Context, id string, dest any) error {
	if dest == nil {
		return fmt.Errorf("query builder: dest cannot be nil")
	}

	q.Where("id", "=", id)
	query, args := q.ToSql()

	return q.conn.SelectOne(ctx, dest, query, args...)
}

func (q *QueryBuilder) Scan(ctx context.Context, dest ...any) error {
	if len(dest) == 0 {
		return fmt.Errorf("query builder: dest cannot be empty")
	}

	query, args := q.ToSql()
	return q.conn.Scan(ctx, query, args, dest...)
}

func (q *QueryBuilder) Count(ctx context.Context) (int64, error) {
	q.query.StmtType = query.StmtCount

	var count int64
	query, args := q.ToSql()

	if err := q.conn.SelectOne(ctx, &count, query, args...); err != nil {
		return 0, err
	}

	return count, nil
}

func (q *QueryBuilder) Exists(ctx context.Context) (bool, error) {
	q.query.StmtType = query.StmtExists

	var count bool
	query, args := q.ToSql()

	if err := q.conn.SelectOne(ctx, &count, query, args...); err != nil {
		return false, err
	}

	return count, nil
}

// Execute performs INSERT, UPDATE, or DELETE operations
func (q *QueryBuilder) Exec(ctx context.Context) error {
	query, args := q.ToSql()
	err := q.conn.Prepared(ctx, query, args...)
	return err
}

// ToSql generates the final SQL query
func (q *QueryBuilder) ToSql() (string, []any) {
	return q.grammar.Compile(q.query)
}

func (q *QueryBuilder) extractQuery(sub connection.QueryBuilder) *query.Query {
	if b, ok := sub.(*QueryBuilder); ok {
		return b.query
	}
	return nil
}
