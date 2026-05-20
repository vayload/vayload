package sqlite3

import (
	"bytes"
	"math"
	"slices"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/vayload/vayload/internal/modules/database/grammar"
	"github.com/vayload/vayload/internal/modules/database/query"
)

type SQLiteQueryGrammar struct{}

func NewQueryGrammar() *SQLiteQueryGrammar {
	return &SQLiteQueryGrammar{}
}

func (g *SQLiteQueryGrammar) Compile(ast *query.Query) (string, []any) {
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

var strPool = sync.Pool{
	New: func() any {
		b := &bytes.Buffer{}
		b.Grow(256)
		return b
	},
}

func newStrBuilder(size int) *bytes.Buffer {
	b := strPool.Get().(*bytes.Buffer)
	b.Reset()

	if size > 256 {
		b.Grow(size)
	}

	return b
}

func putStrBuilder(b *bytes.Buffer) {
	if b.Cap() <= 4096 {
		strPool.Put(b)
	}
}

var argsPool = sync.Pool{
	New: func() any {
		s := make([]any, 0, 16)
		return &s
	},
}

func getArgs() *[]any {
	args := argsPool.Get().(*[]any)
	*args = (*args)[:0]
	return args
}

func putArgs(args *[]any) {
	argsPool.Put(args)
}

// buildSelect generates SELECT queries
func (g *SQLiteQueryGrammar) buildSelect(ast *query.Query) (string, []any) {
	buf := newStrBuilder(128 + len(ast.Columns)*24 + len(ast.Where)*32 + len(ast.Joins)*48)
	bindings := getArgs()

	g.buildSelectToBuf(buf, bindings, ast)
	sql := bytesToString(buf.Bytes())

	putStrBuilder(buf)
	putArgs(bindings)
	return sql, *bindings
}

func (g *SQLiteQueryGrammar) buildSelectToBuf(buf *bytes.Buffer, bindings *[]any, ast *query.Query) {
	buf.WriteString("SELECT ")
	if ast.Distinct {
		buf.WriteString("DISTINCT ")
	}

	if len(ast.Columns) == 0 || (len(ast.Columns) == 1 && ast.Columns[0] == "*") {
		buf.WriteByte('*')
	} else {
		for i := range ast.Columns {
			if i > 0 {
				buf.WriteByte(',')
				buf.WriteByte(' ')
			}
			buf.WriteString(ast.Columns[i])
		}
	}

	buf.WriteString(" FROM ")
	if ast.SubQuery != nil {
		buf.WriteString("(")
		// Recursively build the subquery
		g.buildSelectToBuf(buf, bindings, ast.SubQuery)
		buf.WriteString(") ")
		if ast.Alias != "" {
			buf.WriteString("AS ")
			buf.WriteString(ast.Alias)
		}
	} else {
		buf.WriteString(ast.Table)
		if ast.Alias != "" {
			buf.WriteString(" AS ")
			buf.WriteString(ast.Alias)
		}
	}

	for _, join := range ast.Joins {
		buf.WriteByte(' ')
		buf.WriteString(string(join.Type))
		buf.WriteByte(' ')
		if join.SubQuery != nil {
			buf.WriteByte('(')
			// Recursively build the subquery (it will append to bindings)
			g.buildSelectToBuf(buf, bindings, join.SubQuery)
			buf.WriteByte(')')
			if join.Alias != "" {
				buf.WriteString("AS ")
				buf.WriteString(join.Alias)
			}
		} else {
			buf.WriteString(join.Table)
			if join.Alias != "" {
				buf.WriteString(" AS ")
				buf.WriteString(join.Alias)
			}
		}
		buf.WriteString(" ON ")
		buf.WriteString(join.On)

		*bindings = append(*bindings, join.Args...)
	}

	// WHERE clause write directly to buffer
	g.appendWheres(buf, bindings, ast)

	if len(ast.GroupBy) > 0 {
		buf.WriteString(" GROUP BY ")
		for i := range ast.GroupBy {
			if i > 0 {
				buf.WriteByte(',')
				buf.WriteByte(' ')
			}
			buf.WriteString(ast.GroupBy[i])
		}
	}

	if len(ast.Having) > 0 {
		buf.WriteString(" HAVING ")
		for i := range ast.Having {
			h := &ast.Having[i]
			if i > 0 {
				buf.WriteByte(' ')
				buf.WriteString("AND ")
			}
			buf.WriteString(h.Condition)
			*bindings = append(*bindings, h.Args...)
		}
	}

	if len(ast.OrderBy) > 0 {
		buf.WriteString(" ORDER BY ")
		for i := range ast.OrderBy {
			o := &ast.OrderBy[i]
			if i > 0 {
				buf.WriteByte(',')
				buf.WriteByte(' ')
			}

			buf.WriteString(o.Column)
			buf.WriteByte(' ')
			buf.WriteString(o.Direction)
		}
	}

	if ast.Limit != nil && *ast.Limit > 0 && *ast.Limit < math.MaxInt64 {
		buf.WriteString(" LIMIT ")
		var tmp [20]byte
		buf.Write(strconv.AppendInt(tmp[:0], *ast.Limit, 10))
	}

	if ast.Offset != nil && *ast.Offset >= 0 && *ast.Offset < math.MaxInt64 {
		buf.WriteString(" OFFSET ")
		var tmp [20]byte
		buf.Write(strconv.AppendInt(tmp[:0], *ast.Offset, 10))
	}

	// SQLite does not support lock modes
	if ast.LockMode != "" {
		// buf.WriteString(" ")
		// buf.WriteString(string(ast.LockMode))
	}

	*bindings = append(ast.Args, *bindings...)

	if len(ast.Unions) > 0 {
		for _, u := range ast.Unions {
			if u.All {
				buf.WriteString(" UNION ALL ")
			} else {
				buf.WriteString(" UNION ")
			}

			// Recursively build the union query
			g.buildSelectToBuf(buf, bindings, u.Query)
		}
	}
}

// buildInsert generates INSERT queries
func (g *SQLiteQueryGrammar) buildInsert(ast *query.Query) (string, []any) {
	if len(ast.InsertValues) == 0 && len(ast.InsertMultiValues) == 0 {
		return "", nil
	}

	buf := newStrBuilder(256)

	var columns []string
	var args []any

	if n := len(ast.InsertMultiValues); n > 0 {
		first := ast.InsertMultiValues[0]

		columns = make([]string, 0, len(first))
		for col := range first {
			columns = append(columns, col)
		}

		args = make([]any, 0, len(columns)*n)

	} else {
		columns = make([]string, 0, len(ast.InsertValues))
		for col := range ast.InsertValues {
			columns = append(columns, col)
		}

		args = make([]any, 0, len(columns))
	}

	buf.WriteString("INSERT INTO ")
	buf.WriteString(ast.Table)
	buf.WriteString(" (")

	for i, c := range columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(c)
	}

	buf.WriteString(") VALUES ")

	if rows := ast.InsertMultiValues; len(rows) > 0 {
		for i := range rows {
			if i > 0 {
				buf.WriteString(", ")
			}

			buf.WriteByte('(')
			row := rows[i]
			for j, col := range columns {
				if j > 0 {
					buf.WriteString(", ")
				}
				buf.WriteByte('?')
				args = append(args, row[col])
			}

			buf.WriteByte(')')
		}

	} else {
		buf.WriteByte('(')
		row := ast.InsertValues
		for i, col := range columns {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteByte('?')
			args = append(args, row[col])
		}
		buf.WriteByte(')')
	}

	sql := bytesToString(buf.Bytes())

	putStrBuilder(buf)

	return sql, args
}

// buildUpdate generates UPDATE queries
func (g *SQLiteQueryGrammar) buildUpdate(ast *query.Query) (string, []any) {
	if len(ast.UpdateValues) == 0 {
		return "", nil
	}

	buf := newStrBuilder(256)
	defer putStrBuilder(buf)

	args := make([]any, 0, len(ast.UpdateValues)+len(ast.Where))

	buf.WriteString("UPDATE ")
	buf.WriteString(ast.Table)
	buf.WriteString(" SET ")

	setParts := make([]string, 0, len(ast.UpdateValues))
	i := 1

	for column, value := range ast.UpdateValues {
		setParts = append(setParts, column+" = "+g.Placeholder(i))
		args = append(args, value)
		i++
	}

	for i, part := range setParts {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(part)
	}

	// WHERE clause
	g.appendWheres(buf, &args, ast)

	return buf.String(), args
}

func (g *SQLiteQueryGrammar) buildUpsert(ast *query.Query) (string, []any) {
	buf := newStrBuilder(256)
	defer putStrBuilder(buf)

	args := make([]any, 0, len(ast.UpsertValues)+len(ast.Where))

	if len(ast.UpsertValues) == 0 {
		return "", nil
	}

	var columns []string
	for col := range ast.UpsertValues {
		columns = append(columns, col)
	}

	buf.WriteString("INSERT INTO ")
	buf.WriteString(ast.Table)
	buf.WriteString(" (")
	for i, c := range columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(c)
	}
	buf.WriteString(") VALUES (")
	for i, col := range columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(g.Placeholder(len(args) + 1))
		args = append(args, ast.UpsertValues[col])
	}
	buf.WriteString(")")

	if len(ast.UpsertColumns) > 0 {
		buf.WriteString(" ON CONFLICT (")
		for i, c := range ast.UpsertColumns {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(c)
		}
		buf.WriteString(") DO UPDATE SET ")

		updated := 0
		for _, col := range columns {
			// Don't update conflict columns typically
			isConflict := slices.Contains(ast.UpsertColumns, col)
			if isConflict {
				continue
			}

			if updated > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(col)
			buf.WriteString(" = EXCLUDED.")
			buf.WriteString(col)
			updated++
		}
	}

	return buf.String(), args
}

// buildDelete generates DELETE queries
func (g *SQLiteQueryGrammar) buildDelete(ast *query.Query) (string, []any) {
	buf := newStrBuilder(256)
	args := getArgs()

	buf.WriteString("DELETE FROM ")
	buf.WriteString(ast.Table)

	// WHERE clause
	g.appendWheres(buf, args, ast)
	sql := bytesToString(buf.Bytes())

	putStrBuilder(buf)
	return sql, *args
}

func (g *SQLiteQueryGrammar) appendWheres(buf *bytes.Buffer, args *[]any, ast *query.Query) {
	if len(ast.Where) == 0 {
		return
	}

	buf.WriteString(" WHERE ")
	for i := range ast.Where {
		w := &ast.Where[i]

		if i > 0 {
			if w.IsOr {
				buf.WriteString(" OR ")
			} else {
				buf.WriteString(" AND ")
			}
		}

		switch w.Type {

		case query.WhereTypeBasic:
			buf.WriteString(w.Column)
			buf.WriteByte(' ')
			buf.WriteString(w.Operator)
			buf.WriteByte(' ')
			buf.WriteByte('?')

			*args = append(*args, w.Value)

		case query.WhereTypeColumn:
			buf.WriteString(w.Column)
			buf.WriteByte(' ')
			buf.WriteString(w.Operator)
			buf.WriteByte(' ')
			buf.WriteString(w.Value.(string))

		case query.WhereTypeBetween:
			buf.WriteString(w.Column)
			buf.WriteString(" BETWEEN ? AND ?")

			*args = append(*args, w.Value, w.Value2)

		case query.WhereTypeNull, query.WhereTypeNotNull:
			buf.WriteString(w.Column)
			buf.WriteByte(' ')
			buf.WriteString(w.Operator)

		case query.WhereTypeIn, query.WhereTypeNotIn:
			values := w.Value.([]any)

			buf.WriteString(w.Column)
			buf.WriteByte(' ')
			buf.WriteString(w.Operator)
			buf.WriteString(" (")

			for i := range values {
				if i > 0 {
					buf.WriteByte(',')
					buf.WriteByte(' ')
				}
				buf.WriteByte('?')
				*args = append(*args, values[i])
			}

			buf.WriteByte(')')

		case query.WhereTypeSubQuery:
			buf.WriteString(w.Column)
			buf.WriteByte(' ')
			buf.WriteString(w.Operator)
			buf.WriteString(" (")

			// append subquery and push args (bind values)
			g.buildSelectToBuf(buf, args, w.SubQuery)
			buf.WriteByte(')')
		}
	}
}

// getPlaceholder returns the appropriate placeholder for the database driver
func (g *SQLiteQueryGrammar) Placeholder(position int) string {
	return "?"
}

func (g *SQLiteQueryGrammar) Wrap(ident string) string {
	if ident == "*" {
		return ident
	}

	q := byte('"')

	var b strings.Builder
	b.Grow(len(ident) + 8)

	i := 0
	start := 0
	n := len(ident)

	for i < n {
		switch ident[i] {

		case '.':
			if start < i {
				part := ident[start:i]
				if part != "*" {
					b.WriteByte(q)
					b.WriteString(part)
					b.WriteByte(q)
				} else {
					b.WriteByte('*')
				}
			}

			b.WriteByte('.')
			start = i + 1

		case ' ':
			part := ident[start:i]

			if part != "*" {
				b.WriteByte(q)
				b.WriteString(part)
				b.WriteByte(q)
			}

			rest := ident[i:]
			b.WriteString(rest)

			return b.String()
		}

		i++
	}

	if start < n {
		part := ident[start:n]

		if part != "*" {
			b.WriteByte(q)
			b.WriteString(part)
			b.WriteByte(q)
		} else {
			b.WriteByte('*')
		}
	}

	return b.String()
}

func (g *SQLiteQueryGrammar) Quote(ident string) string {
	if ident == "" {
		return ident
	}

	q := byte('"')

	var b strings.Builder
	b.Grow(len(ident) + 2)

	b.WriteByte(q)
	b.WriteString(ident)
	b.WriteByte(q)

	return b.String()
}

func (g *SQLiteQueryGrammar) WrapColumn(column string) string {
	if column == "*" {
		return column
	}

	return g.Wrap(column)
}

func (g *SQLiteQueryGrammar) WrapTable(table string) string {
	return g.Wrap(table)
}

var _ grammar.QueryGrammar = (*SQLiteQueryGrammar)(nil)

func bytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}

	return unsafe.String(unsafe.SliceData(b), len(b))
}
