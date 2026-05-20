package postgres

import (
	"fmt"
	"math"
	"strings"

	"github.com/vayload/vayload/internal/modules/database/grammar"
	"github.com/vayload/vayload/internal/modules/database/query"
)

type PostgresQueryGrammar struct{}

func NewQueryGrammar() *PostgresQueryGrammar {
	return &PostgresQueryGrammar{}
}

func (g *PostgresQueryGrammar) Compile(ast *query.Query) (string, []any) {
	switch ast.StmtType {
	case query.StmtSelect:
		return g.buildSelect(ast)
	case query.StmtDelete:
		return g.buildDelete(ast)
	case query.StmtInsert:
		return g.buildInsert(ast)
	case query.StmtUpdate:
		return g.buildUpdate(ast)
	case query.StmtUpsert:
		return g.buildUpsert(ast)
	default:
		return g.buildSelect(ast)
	}
}

// buildSelect generates SELECT queries
func (g *PostgresQueryGrammar) buildSelect(ast *query.Query) (string, []any) {
	var query strings.Builder
	args := make([]any, 0)

	// SELECT clause
	query.WriteString("SELECT ")
	if ast.Distinct {
		query.WriteString("DISTINCT ")
	}
	query.WriteString(strings.Join(ast.Columns, ", "))

	// FROM clause
	query.WriteString(" FROM ")
	query.WriteString(ast.Table)

	// JOIN clauses
	for _, join := range ast.Joins {
		query.WriteString(fmt.Sprintf(" %s %s ON %s", join.Type, join.Table, join.On))
		args = append(args, join.Args...)
	}

	// WHERE clause
	whereSQL, whereArgs := g.buildWhere(ast)
	if whereSQL != "" {
		query.WriteString(" WHERE ")
		query.WriteString(whereSQL)
		args = append(args, whereArgs...)
	}

	// GROUP BY clause
	if len(ast.GroupBy) > 0 {
		query.WriteString(" GROUP BY ")
		query.WriteString(strings.Join(ast.GroupBy, ", "))
	}

	// HAVING clause
	if len(ast.Having) > 0 {
		query.WriteString(" HAVING ")
		havingParts := make([]string, len(ast.Having))
		for i, h := range ast.Having {
			havingParts[i] = h.Condition
			args = append(args, h.Args...)
		}
		query.WriteString(strings.Join(havingParts, " AND "))
	}

	// ORDER BY clause
	if len(ast.OrderBy) > 0 {
		query.WriteString(" ORDER BY ")
		orderParts := make([]string, len(ast.OrderBy))
		for i, o := range ast.OrderBy {
			orderParts[i] = fmt.Sprintf("%s %s", o.Column, o.Direction)
		}

		query.WriteString(strings.Join(orderParts, ", "))
	}

	// LIMIT clause
	if ast.Limit != nil && *ast.Limit > 0 && *ast.Limit < math.MaxInt64 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", *ast.Limit))
	}

	// OFFSET clause
	if ast.Offset != nil && *ast.Offset >= 0 && *ast.Offset < math.MaxInt64 {
		query.WriteString(fmt.Sprintf(" OFFSET %d", *ast.Offset))
	}

	allArgs := append(ast.Args, args...)

	return query.String(), allArgs
}

// buildInsert generates INSERT queries
func (g *PostgresQueryGrammar) buildInsert(ast *query.Query) (string, []any) {
	var query strings.Builder
	args := make([]any, 0)

	if len(ast.InsertValues) == 0 {
		return "", nil
	}

	columns := make([]string, 0, len(ast.InsertValues))
	placeholders := make([]string, 0, len(ast.InsertValues))

	// Sort keys for consistent ordering
	i := 1
	for column, value := range ast.InsertValues {
		columns = append(columns, column)
		placeholders = append(placeholders, g.Placeholder(i))
		args = append(args, value)
		i++
	}

	query.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		ast.Table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", ")))

	return query.String(), args
}

// buildUpdate generates UPDATE queries
func (g *PostgresQueryGrammar) buildUpdate(ast *query.Query) (string, []any) {
	var query strings.Builder
	args := make([]any, 0)

	if len(ast.UpdateValues) == 0 {
		return "", nil
	}

	query.WriteString(fmt.Sprintf("UPDATE %s SET ", ast.Table))

	setParts := make([]string, 0, len(ast.UpdateValues))
	i := 1

	for column, value := range ast.UpdateValues {
		setParts = append(setParts, fmt.Sprintf("%s = %s", column, g.Placeholder(i)))
		args = append(args, value)
		i++
	}

	query.WriteString(strings.Join(setParts, ", "))

	// WHERE clause
	whereSQL, whereArgs := g.buildWhere(ast)
	if whereSQL != "" {
		query.WriteString(" WHERE ")

		// Adjust placeholders for WHERE clause
		adjustedWhere := whereSQL
		for j := 0; j < len(whereArgs); j++ {
			oldPlaceholder := g.Placeholder(j + 1)
			newPlaceholder := g.Placeholder(i + j)
			adjustedWhere = strings.Replace(adjustedWhere, oldPlaceholder, newPlaceholder, 1)
		}

		query.WriteString(adjustedWhere)
		args = append(args, whereArgs...)
	}

	return query.String(), args
}

func (g *PostgresQueryGrammar) buildUpsert(ast *query.Query) (string, []any) {
	var query strings.Builder
	args := make([]any, 0)

	if len(ast.UpsertValues) == 0 {
		return "", nil
	}

	columns := make([]string, 0, len(ast.UpsertValues))
	placeholders := make([]string, 0, len(ast.UpsertValues))

	i := 1
	for column, value := range ast.UpsertValues {
		columns = append(columns, column)
		placeholders = append(placeholders, g.Placeholder(i))
		args = append(args, value)
		i++
	}

	query.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		ast.Table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", ")))

	return query.String(), args
}

// buildDelete generates DELETE queries
func (g *PostgresQueryGrammar) buildDelete(ast *query.Query) (string, []any) {
	var query strings.Builder
	args := make([]any, 0)

	query.WriteString(fmt.Sprintf("DELETE FROM %s", ast.Table))

	// WHERE clause
	whereSQL, whereArgs := g.buildWhere(ast)
	if whereSQL != "" {
		query.WriteString(" WHERE ")
		query.WriteString(whereSQL)
		args = append(args, whereArgs...)
	}

	return query.String(), args
}

// buildWhere generates WHERE clause
func (g *PostgresQueryGrammar) buildWhere(ast *query.Query) (string, []any) {
	if len(ast.Where) == 0 && len(ast.OrWhere) == 0 {
		return "", nil
	}

	var parts []string
	args := make([]any, 0)
	position := 1

	// AND WHERE clauses
	for _, w := range ast.Where {
		switch w.Operator {
		case "IS NULL", "IS NOT NULL":
			parts = append(parts, fmt.Sprintf("%s %s", w.Column, w.Operator))
		case "IN", "NOT IN":
			parts = append(parts, fmt.Sprintf("%s %s %v", w.Column, w.Operator, w.Value))
		default:
			placeholder := g.Placeholder(position)
			parts = append(parts, fmt.Sprintf("%s %s %s", w.Column, w.Operator, placeholder))
			args = append(args, w.Value)
			position++
		}
	}

	var whereSQL strings.Builder
	whereSQL.WriteString(strings.Join(parts, " AND "))

	// OR WHERE clauses
	for _, w := range ast.OrWhere {
		switch w.Operator {
		case "IS NULL", "IS NOT NULL":
			whereSQL.WriteString(fmt.Sprintf(" OR %s %s", w.Column, w.Operator))
		case "IN", "NOT IN":
			whereSQL.WriteString(fmt.Sprintf(" OR %s %s %v", w.Column, w.Operator, w.Value))
		default:
			placeholder := g.Placeholder(position)
			whereSQL.WriteString(fmt.Sprintf(" OR %s %s %s", w.Column, w.Operator, placeholder))
			args = append(args, w.Value)
			position++
		}
	}

	return whereSQL.String(), args
}

// getPlaceholder returns the appropriate placeholder for the database driver
func (g *PostgresQueryGrammar) Placeholder(position int) string {
	return "?"
}

func (g *PostgresQueryGrammar) Wrap(value string) string {
	return ""
}

var _ grammar.QueryGrammar = (*PostgresQueryGrammar)(nil)
