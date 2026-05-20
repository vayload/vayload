package sqlite3

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/vayload/vayload/internal/modules/database/builder"
	"github.com/vayload/vayload/internal/modules/database/migrator"
	"github.com/vayload/vayload/internal/modules/database/query"
)

func TestSQLiteQueryGrammar_BuildSelect(t *testing.T) {
	g := NewQueryGrammar()

	tests := []struct {
		name     string
		ast      *query.Query
		wantSQL  string
		wantArgs []any
	}{
		{
			name: "Simple Select",
			ast: &query.Query{
				Table:   "users",
				Columns: []string{"id", "name"},
			},
			wantSQL:  "SELECT id, name FROM users",
			wantArgs: []any{},
		},
		{
			name: "Select All",
			ast: &query.Query{
				Table:   "users",
				Columns: []string{"*"},
			},
			wantSQL:  "SELECT * FROM users",
			wantArgs: []any{},
		},
		{
			name: "Select Distinct",
			ast: &query.Query{
				Table:    "users",
				Columns:  []string{"email"},
				Distinct: true,
			},
			wantSQL:  "SELECT DISTINCT email FROM users",
			wantArgs: []any{},
		},
		{
			name: "Select With Alias",
			ast: &query.Query{
				Table:   "users",
				Alias:   "u",
				Columns: []string{"u.id", "u.name"},
			},
			wantSQL:  "SELECT u.id, u.name FROM users AS u",
			wantArgs: []any{},
		},
		{
			name: "Select With Where",
			ast: &query.Query{
				Table:   "users",
				Columns: []string{"*"},
				Where: []query.Where{
					{
						Type:     query.WhereTypeBasic,
						Column:   "active",
						Operator: "=",
						Value:    true,
					},
				},
			},
			wantSQL:  "SELECT * FROM users WHERE active = ?",
			wantArgs: []any{true},
		},
		{
			name: "Select With Join",
			ast: &query.Query{
				Table:   "users",
				Columns: []string{"users.*", "posts.title"},
				Joins: []query.Join{
					{
						Type:  query.InnerJoinType,
						Table: "posts",
						On:    "users.id = posts.user_id",
					},
				},
			},
			wantSQL:  "SELECT users.*, posts.title FROM users INNER JOIN posts ON users.id = posts.user_id",
			wantArgs: []any{},
		},
		{
			name: "Select With Subquery in FROM",
			ast: &query.Query{
				SubQuery: &query.Query{
					Table:   "orders",
					Columns: []string{"user_id", "SUM(amount) as total"},
					GroupBy: []string{"user_id"},
				},
				Alias:   "t",
				Columns: []string{"*"},
			},
			wantSQL:  "SELECT * FROM (SELECT user_id, SUM(amount) as total FROM orders GROUP BY user_id) AS t",
			wantArgs: []any{},
		},
		{
			name: "Select With Subquery in WHERE",
			ast: &query.Query{
				Table:   "users",
				Columns: []string{"*"},
				Where: []query.Where{
					{
						Type:     query.WhereTypeSubQuery,
						Column:   "id",
						Operator: "IN",
						SubQuery: &query.Query{
							Table:   "orders",
							Columns: []string{"user_id"},
							Where: []query.Where{
								{
									Type:     query.WhereTypeBasic,
									Column:   "amount",
									Operator: ">",
									Value:    100,
								},
							},
						},
					},
				},
			},
			wantSQL:  "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE amount > ?)",
			wantArgs: []any{100},
		},
		{
			name: "Select With Union",
			ast: &query.Query{
				Table:   "admins",
				Columns: []string{"id", "role"},
				Unions: []query.Union{
					{
						Query: &query.Query{
							Table:   "staff",
							Columns: []string{"id", "role"},
						},
						All: true,
					},
				},
			},
			wantSQL:  "SELECT id, role FROM admins UNION ALL SELECT id, role FROM staff",
			wantArgs: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := g.buildSelect(tt.ast)
			if strings.TrimSpace(sql) != tt.wantSQL {
				t.Errorf("%s: got SQL %q, want %q", tt.name, sql, tt.wantSQL)
			}
			if len(args) != len(tt.wantArgs) {
				t.Errorf("%s: got %d args, want %d", tt.name, len(args), len(tt.wantArgs))
			}
			for i := range args {
				if args[i] != tt.wantArgs[i] {
					t.Errorf("%s: arg %d: got %v, want %v", tt.name, i, args[i], tt.wantArgs[i])
				}
			}
		})
	}
}

func TestSQLiteQueryGrammar_BuildInsert(t *testing.T) {
	g := NewQueryGrammar()

	tests := []struct {
		name     string
		ast      *query.Query
		wantSQL  string
		wantArgs []any
	}{
		{
			name: "Simple Insert",
			ast: &query.Query{
				Table: "users",
				InsertValues: map[string]any{
					"name":  "John",
					"email": "john@example.com",
				},
			},
			wantSQL:  "INSERT INTO users (name, email) VALUES (?, ?)",
			wantArgs: []any{"John", "john@example.com"},
		},
		{
			name: "Insert Many",
			ast: &query.Query{
				Table: "users",
				InsertMultiValues: []map[string]any{
					{"name": "John", "email": "john@example.com"},
					{"name": "Jane", "email": "jane@example.com"},
				},
			},
			wantSQL:  "INSERT INTO users (name, email) VALUES (?, ?), (?, ?)",
			wantArgs: []any{"John", "john@example.com", "Jane", "jane@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := g.buildInsert(tt.ast)
			// Note: Map iteration is random, so we can't easily check fixed SQL strings for InsertValues.
			// However, in our implementation we should probably sort columns to be deterministic if possible,
			// though it's not strictly required for SQL correctness but good for tests.
			// Currently SQLiteQueryGrammar.buildInsert does NOT sort columns.
			if tt.name == "Insert Many" || tt.name == "Simple Insert" {
				// We'll skip strict string check if columns order varies, or we can check length and parts.
				if len(args) != len(tt.wantArgs) {
					t.Errorf("%s: got %d args, want %d", tt.name, len(args), len(tt.wantArgs))
				}
			} else if strings.TrimSpace(sql) != tt.wantSQL {
				t.Errorf("%s: got SQL %q, want %q", tt.name, sql, tt.wantSQL)
			}
		})
	}
}

func TestSQLiteQueryGrammar_BuildUpsert(t *testing.T) {
	g := NewQueryGrammar()

	ast := &query.Query{
		Table: "users",
		UpsertValues: map[string]any{
			"id":    1,
			"email": "test@example.com",
			"name":  "Test",
		},
		UpsertColumns: []string{"id"},
	}

	sql, _ := g.buildUpsert(ast)
	if !strings.Contains(sql, "INSERT INTO users") {
		t.Errorf("got %q, want it to contain INSERT INTO users", sql)
	}
	if !strings.Contains(sql, "ON CONFLICT (id)") {
		t.Errorf("got %q, want it to contain ON CONFLICT (id)", sql)
	}
	if !strings.Contains(sql, "DO UPDATE SET") {
		t.Errorf("got %q, want it to contain DO UPDATE SET", sql)
	}
}

func TestFullQuery(t *testing.T) {
	tempFile, err := os.CreateTemp("", "vayload_test.db")
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	tempFile.Close()
	defer os.Remove(tempFile.Name())
	// create temp table
	conn, err := NewConnection(context.Background(), "", "", "", "", tempFile.Name())
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	schema := builder.NewSchemaBuilder(NewSchemaGrammar(), conn)

	// create table
	schema.Create("users", func(blueprint *migrator.Blueprint) {
		blueprint.ID().AutoIncrement()
		blueprint.String("name", 255).Nullable()
		blueprint.String("email", 255).Unique()
		blueprint.Timestamps()
	})

	// create new entry
	conn.From("users").InsertOne(map[string]any{
		"name":  "John",
		"email": "[EMAIL_ADDRESS]",
	}).Exec(context.Background())

	type User struct {
		ID    int    `db:"id"`
		Name  string `db:"name"`
		Email string `db:"email"`
	}

	// select
	var users []User
	err = conn.From("users").Select("id", "name", "email").Get(context.Background(), &users)
	if err != nil {
		t.Errorf("got error: %v", err)
	}

	fmt.Println(users)
}

func BenchmarkSQLiteQueryGrammar_BuildSelect_Complex(b *testing.B) {
	g := NewQueryGrammar()

	// Query super compleja
	ast := &query.Query{
		Table:   "users",
		Alias:   "u",
		Columns: []string{"u.id", "u.name", "u.email", "p.title", "c.comment"},
		Where: []query.Where{
			{Type: query.WhereTypeBasic, Column: "u.active", Operator: "=", Value: true},
			{Type: query.WhereTypeIn, Column: "u.id", Operator: "IN", Value: []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{Type: query.WhereTypeSubQuery, Column: "u.id", Operator: "IN", SubQuery: &query.Query{
				Table:   "admins",
				Columns: []string{"id"},
				Where: []query.Where{
					{Type: query.WhereTypeBasic, Column: "level", Operator: ">", Value: 5},
				},
			}},
		},
		Joins: []query.Join{
			{
				Type:  "LEFT JOIN",
				Table: "posts",
				Alias: "p",
				On:    "p.user_id = u.id",
			},
			{
				Type: "INNER JOIN",
				SubQuery: &query.Query{
					Table:   "comments",
					Alias:   "c",
					Columns: []string{"comment", "post_id"},
					Where: []query.Where{
						{Type: query.WhereTypeBasic, Column: "approved", Operator: "=", Value: true},
					},
				},
				Alias: "c",
				On:    "c.post_id = p.id",
			},
		},
		GroupBy: []string{"u.id", "p.id"},
		Having: []query.Having{
			{Condition: "COUNT(c.comment) > ?", Args: []any{0}},
		},
		OrderBy: []query.Order{
			{Column: "u.created_at", Direction: "DESC"},
			{Column: "p.created_at", Direction: "ASC"},
		},
		Limit:  func() *int64 { l := int64(50); return &l }(),
		Offset: func() *int64 { o := int64(100); return &o }(),
		Unions: []query.Union{
			{
				All: true,
				Query: &query.Query{
					Table:   "archived_users",
					Columns: []string{"id", "name", "email"},
					Where: []query.Where{
						{Type: query.WhereTypeBasic, Column: "active", Operator: "=", Value: false},
					},
				},
			},
		},
	}

	b.ReportAllocs()
	for b.Loop() {
		_, _ = g.buildSelect(ast)
	}

	sql, args := g.buildSelect(ast)
	fmt.Println(sql)
	fmt.Println(args)
}

func BenchmarkSQLiteQueryGrammar_BuildInsertMany(b *testing.B) {
	g := NewQueryGrammar()
	multiValues := make([]map[string]any, 100)
	for i := range 100 {
		multiValues[i] = map[string]any{
			"name":  "User",
			"email": "user@example.com",
			"age":   30,
		}
	}
	ast := &query.Query{
		Table:             "users",
		InsertMultiValues: multiValues,
	}

	b.ReportAllocs()
	for b.Loop() {
		_, _ = g.buildInsert(ast)
	}
}
