package builder

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/vayload/vayload/internal/services/database/connection"
)

type QueryBuilder struct {
	conn   connection.DatabaseConnection
	table  string
	driver connection.DatabaseDriver

	stmtType  StatementType
	columns   []string
	distinct  bool
	joins     []joinClause
	where     []whereClause
	orWhere   []whereClause
	groupBy   []string
	having    []havingClause
	orderBy   []orderClause
	limitVal  *int
	offsetVal *int
	args      []any

	// For mutations
	insertValues map[string]any
	updateValues map[string]any
}

type joinClause struct {
	joinType string
	table    string
	on       string
	args     []any
}

type whereClause struct {
	column   string
	operator string
	value    any
	isOr     bool
}

type havingClause struct {
	condition string
	args      []any
}

type orderClause struct {
	column    string
	direction string
}

type StatementType int

const (
	StmtSelect StatementType = iota
	StmtInsert
	StmtUpdate
	StmtDelete
)

func NewQueryBuilder(conn connection.DatabaseConnection, table string, driver connection.DatabaseDriver) connection.QueryBuilder {
	return &QueryBuilder{
		conn:     conn,
		table:    table,
		driver:   driver,
		stmtType: StmtSelect,
		columns:  []string{"*"},
	}
}

func (q *QueryBuilder) From(table string) connection.QueryBuilder {
	q.table = table
	return q
}

func (q *QueryBuilder) Select(columns ...string) connection.QueryBuilder {
	q.columns = columns
	return q
}

func (q *QueryBuilder) Distinct() connection.QueryBuilder {
	q.distinct = true
	return q
}

// Joins
func (q *QueryBuilder) Join(table string, on string, args ...any) connection.QueryBuilder {
	q.joins = append(q.joins, joinClause{
		joinType: "JOIN",
		table:    table,
		on:       on,
		args:     args,
	})
	return q
}

func (q *QueryBuilder) LeftJoin(table string, on string, args ...any) connection.QueryBuilder {
	q.joins = append(q.joins, joinClause{
		joinType: "LEFT JOIN",
		table:    table,
		on:       on,
		args:     args,
	})
	return q
}

func (q *QueryBuilder) RightJoin(table string, on string, args ...any) connection.QueryBuilder {
	q.joins = append(q.joins, joinClause{
		joinType: "RIGHT JOIN",
		table:    table,
		on:       on,
		args:     args,
	})
	return q
}

func (q *QueryBuilder) InnerJoin(table string, on string, args ...any) connection.QueryBuilder {
	q.joins = append(q.joins, joinClause{
		joinType: "INNER JOIN",
		table:    table,
		on:       on,
		args:     args,
	})
	return q
}

// Where clauses
func (q *QueryBuilder) Where(column string, operator string, value any) connection.QueryBuilder {
	q.where = append(q.where, whereClause{
		column:   column,
		operator: operator,
		value:    value,
		isOr:     false,
	})
	return q
}

func (q *QueryBuilder) OrWhere(column string, operator string, value any) connection.QueryBuilder {
	q.orWhere = append(q.orWhere, whereClause{
		column:   column,
		operator: operator,
		value:    value,
		isOr:     true,
	})
	return q
}

func (q *QueryBuilder) WhereIn(column string, values ...any) connection.QueryBuilder {
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = q.getPlaceholder(len(q.args) + i + 1)
		q.args = append(q.args, values[i])
	}

	q.where = append(q.where, whereClause{
		column:   column,
		operator: "IN",
		value:    fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")),
		isOr:     false,
	})
	return q
}

func (q *QueryBuilder) WhereNotIn(column string, values ...any) connection.QueryBuilder {
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = q.getPlaceholder(len(q.args) + i + 1)
		q.args = append(q.args, values[i])
	}

	q.where = append(q.where, whereClause{
		column:   column,
		operator: "NOT IN",
		value:    fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")),
		isOr:     false,
	})
	return q
}

func (q *QueryBuilder) WhereNull(column string) connection.QueryBuilder {
	q.where = append(q.where, whereClause{
		column:   column,
		operator: "IS NULL",
		value:    nil,
		isOr:     false,
	})
	return q
}

func (q *QueryBuilder) WhereNotNull(column string) connection.QueryBuilder {
	q.where = append(q.where, whereClause{
		column:   column,
		operator: "IS NOT NULL",
		value:    nil,
		isOr:     false,
	})
	return q
}

// Mutation methods
func (q *QueryBuilder) Insert(values map[string]any) connection.QueryBuilder {
	q.stmtType = StmtInsert
	q.insertValues = values
	return q
}

func (q *QueryBuilder) Update(values map[string]any) connection.QueryBuilder {
	q.stmtType = StmtUpdate
	q.updateValues = values
	return q
}

func (q *QueryBuilder) Delete() connection.QueryBuilder {
	q.stmtType = StmtDelete
	return q
}

// Group and Having
func (q *QueryBuilder) GroupBy(columns ...string) connection.QueryBuilder {
	q.groupBy = columns
	return q
}

func (q *QueryBuilder) Having(condition string, args ...any) connection.QueryBuilder {
	q.having = append(q.having, havingClause{
		condition: condition,
		args:      args,
	})
	return q
}

// Order and Pagination
func (q *QueryBuilder) OrderBy(column string, direction string) connection.QueryBuilder {
	q.orderBy = append(q.orderBy, orderClause{
		column:    column,
		direction: strings.ToUpper(direction),
	})
	return q
}

func (q *QueryBuilder) Limit(limit int) connection.QueryBuilder {
	q.limitVal = &limit
	return q
}

func (q *QueryBuilder) Offset(offset int) connection.QueryBuilder {
	q.offsetVal = &offset
	return q
}

// Execution methods
func (q *QueryBuilder) Get(dest any) error {
	query, args := q.ToSql()

	return q.conn.SelectOne(context.Background(), dest, query, args...)
}

func (q *QueryBuilder) GetAll(dest any) error {
	query, args := q.ToSql()
	return q.conn.Select(context.Background(), dest, query, args...)
}

func (q *QueryBuilder) Scan(dest ...any) error {
	query, args := q.ToSql()
	return q.conn.SelectOne(context.Background(), dest, query, args...)
}

func (q *QueryBuilder) Count() (int64, error) {
	// Save original columns and statement type
	originalCols := q.columns
	originalStmt := q.stmtType
	q.columns = []string{"COUNT(*) as count"}
	q.stmtType = StmtSelect

	var count int64
	query, args := q.ToSql()

	// Restore original values
	q.columns = originalCols
	q.stmtType = originalStmt

	if err := q.conn.SelectOne(context.Background(), &count, query, args...); err != nil {
		return 0, err
	}

	return count, nil
}

func (q *QueryBuilder) Exists() (bool, error) {
	count, err := q.Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Execute performs INSERT, UPDATE, or DELETE operations
func (q *QueryBuilder) Exec() error {
	query, args := q.ToSql()
	err := q.conn.Prepared(context.Background(), query, args...)
	return err
}

// ToSql generates the final SQL query
func (q *QueryBuilder) ToSql() (string, []any) {
	switch q.stmtType {
	case StmtInsert:
		return q.buildInsert()
	case StmtUpdate:
		return q.buildUpdate()
	case StmtDelete:
		return q.buildDelete()
	default:
		return q.buildSelect()
	}
}

// buildSelect generates SELECT queries
func (q *QueryBuilder) buildSelect() (string, []any) {
	var query strings.Builder
	args := make([]any, 0)

	// SELECT clause
	query.WriteString("SELECT ")
	if q.distinct {
		query.WriteString("DISTINCT ")
	}
	query.WriteString(strings.Join(q.columns, ", "))

	// FROM clause
	query.WriteString(" FROM ")
	query.WriteString(q.table)

	// JOIN clauses
	for _, join := range q.joins {
		query.WriteString(fmt.Sprintf(" %s %s ON %s", join.joinType, join.table, join.on))
		args = append(args, join.args...)
	}

	// WHERE clause
	whereSQL, whereArgs := q.buildWhere()
	if whereSQL != "" {
		query.WriteString(" WHERE ")
		query.WriteString(whereSQL)
		args = append(args, whereArgs...)
	}

	// GROUP BY clause
	if len(q.groupBy) > 0 {
		query.WriteString(" GROUP BY ")
		query.WriteString(strings.Join(q.groupBy, ", "))
	}

	// HAVING clause
	if len(q.having) > 0 {
		query.WriteString(" HAVING ")
		havingParts := make([]string, len(q.having))
		for i, h := range q.having {
			havingParts[i] = h.condition
			args = append(args, h.args...)
		}
		query.WriteString(strings.Join(havingParts, " AND "))
	}

	// ORDER BY clause
	if len(q.orderBy) > 0 {
		query.WriteString(" ORDER BY ")
		orderParts := make([]string, len(q.orderBy))
		for i, o := range q.orderBy {
			orderParts[i] = fmt.Sprintf("%s %s", o.column, o.direction)
		}
		query.WriteString(strings.Join(orderParts, ", "))
	}

	// LIMIT clause
	if q.limitVal != nil && *q.limitVal > 0 && *q.limitVal < math.MaxInt64 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", *q.limitVal))
	}

	// OFFSET clause
	if q.offsetVal != nil && *q.offsetVal >= 0 && *q.offsetVal < math.MaxInt64 {
		query.WriteString(fmt.Sprintf(" OFFSET %d", *q.offsetVal))
	}

	// Combine with stored args
	allArgs := append(q.args, args...)

	return query.String(), allArgs
}

// buildInsert generates INSERT queries
func (q *QueryBuilder) buildInsert() (string, []any) {
	var query strings.Builder
	args := make([]any, 0)

	if len(q.insertValues) == 0 {
		return "", nil
	}

	columns := make([]string, 0, len(q.insertValues))
	placeholders := make([]string, 0, len(q.insertValues))

	// Sort keys for consistent ordering
	i := 1
	for column, value := range q.insertValues {
		columns = append(columns, column)
		placeholders = append(placeholders, q.getPlaceholder(i))
		args = append(args, value)
		i++
	}

	query.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		q.table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", ")))

	return query.String(), args
}

// buildUpdate generates UPDATE queries
func (q *QueryBuilder) buildUpdate() (string, []any) {
	var query strings.Builder
	args := make([]any, 0)

	if len(q.updateValues) == 0 {
		return "", nil
	}

	query.WriteString(fmt.Sprintf("UPDATE %s SET ", q.table))

	setParts := make([]string, 0, len(q.updateValues))
	i := 1

	for column, value := range q.updateValues {
		setParts = append(setParts, fmt.Sprintf("%s = %s", column, q.getPlaceholder(i)))
		args = append(args, value)
		i++
	}

	query.WriteString(strings.Join(setParts, ", "))

	// WHERE clause
	whereSQL, whereArgs := q.buildWhere()
	if whereSQL != "" {
		query.WriteString(" WHERE ")

		// Adjust placeholders for WHERE clause
		adjustedWhere := whereSQL
		for j := 0; j < len(whereArgs); j++ {
			oldPlaceholder := q.getPlaceholder(j + 1)
			newPlaceholder := q.getPlaceholder(i + j)
			adjustedWhere = strings.Replace(adjustedWhere, oldPlaceholder, newPlaceholder, 1)
		}

		query.WriteString(adjustedWhere)
		args = append(args, whereArgs...)
	}

	return query.String(), args
}

// buildDelete generates DELETE queries
func (q *QueryBuilder) buildDelete() (string, []any) {
	var query strings.Builder
	args := make([]any, 0)

	query.WriteString(fmt.Sprintf("DELETE FROM %s", q.table))

	// WHERE clause
	whereSQL, whereArgs := q.buildWhere()
	if whereSQL != "" {
		query.WriteString(" WHERE ")
		query.WriteString(whereSQL)
		args = append(args, whereArgs...)
	}

	return query.String(), args
}

// buildWhere generates WHERE clause
func (q *QueryBuilder) buildWhere() (string, []any) {
	if len(q.where) == 0 && len(q.orWhere) == 0 {
		return "", nil
	}

	var parts []string
	args := make([]any, 0)
	position := 1

	// AND WHERE clauses
	for _, w := range q.where {
		switch w.operator {
		case "IS NULL", "IS NOT NULL":
			parts = append(parts, fmt.Sprintf("%s %s", w.column, w.operator))
		case "IN", "NOT IN":
			parts = append(parts, fmt.Sprintf("%s %s %v", w.column, w.operator, w.value))
		default:
			placeholder := q.getPlaceholder(position)
			parts = append(parts, fmt.Sprintf("%s %s %s", w.column, w.operator, placeholder))
			args = append(args, w.value)
			position++
		}
	}

	whereSQL := strings.Join(parts, " AND ")

	// OR WHERE clauses
	for _, w := range q.orWhere {
		switch w.operator {
		case "IS NULL", "IS NOT NULL":
			whereSQL += fmt.Sprintf(" OR %s %s", w.column, w.operator)
		case "IN", "NOT IN":
			whereSQL += fmt.Sprintf(" OR %s %s %v", w.column, w.operator, w.value)
		default:
			placeholder := q.getPlaceholder(position)
			whereSQL += fmt.Sprintf(" OR %s %s %s", w.column, w.operator, placeholder)
			args = append(args, w.value)
			position++
		}
	}

	return whereSQL, args
}

// getPlaceholder returns the appropriate placeholder for the database driver
func (q *QueryBuilder) getPlaceholder(position int) string {
	switch q.driver {
	case connection.PostgreSQLDriver:
		return fmt.Sprintf("$%d", position)
	case connection.MySQLDriver:
		return "?"
	case connection.SQLiteDriver:
		return "?"
	default:
		return "?"
	}
}
