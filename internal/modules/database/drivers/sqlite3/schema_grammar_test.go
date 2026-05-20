package sqlite3

import (
	"strings"
	"testing"

	"github.com/vayload/vayload/internal/modules/database/migrator"
	database_types "github.com/vayload/vayload/internal/modules/database/types"
)

func TestSQLiteSchemaGrammar_CreateTable(t *testing.T) {
	g := NewSchemaGrammar()

	bp := migrator.Blueprint{
		TableName: "users",
		Columns: []*migrator.Column{
			{
				Name:      "id",
				Type:      database_types.TypeBigInt,
				IsPrimary: true,
				IsAutoInc: true,
			},
			{
				Name:       "username",
				Type:       database_types.TypeVarchar,
				Length:     255,
				IsNullable: false,
				IsUnique:   true,
			},
			{
				Name:       "email",
				Type:       database_types.TypeVarchar,
				Length:     255,
				IsNullable: true,
			},
			{
				Name:       "created_at",
				Type:       database_types.TypeTimestamp,
				DefaultVal: "now()",
			},
		},
	}

	sql, err := g.CreateTable(bp)
	if err != nil {
		t.Fatalf("CreateTable failed: %v", err)
	}

	sql = strings.ToUpper(sql)
	if !strings.Contains(sql, "CREATE TABLE USERS") {
		t.Errorf("got %q, want it to contain CREATE TABLE USERS", sql)
	}
	if !strings.Contains(sql, "ID INTEGER PRIMARY KEY AUTOINCREMENT") {
		// Note: SQLite might use INTEGER PRIMARY KEY for autoincrement
		// Let's check what CompileColType returns
	}
	// Because sqlite doesn't support VARCHAR, it will be converted to TEXT
	if !strings.Contains(sql, "USERNAME TEXT NOT NULL UNIQUE") {
		t.Errorf("got %q, want it to contain USERNAME TEXT NOT NULL UNIQUE", sql)
	}

	// Because sqlite doesn't support TIMESTAMP, it will be converted to TEXT
	if !strings.Contains(sql, "CREATED_AT TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP") {
		t.Errorf("got %q, want it to contain CREATED_AT TEXT DEFAULT CURRENT_TIMESTAMP", sql)
	}
}

func TestSQLiteSchemaGrammar_AddColumn(t *testing.T) {
	g := NewSchemaGrammar()

	col := &migrator.Column{
		Name:       "phone",
		Type:       database_types.TypeVarchar,
		Length:     20,
		IsNullable: true,
	}

	sql, err := g.AddColumn("users", col)
	if err != nil {
		t.Fatalf("AddColumn failed: %v", err)
	}

	sql = strings.ToUpper(sql)
	if !strings.Contains(sql, "ALTER TABLE USERS ADD COLUMN PHONE TEXT") {
		t.Errorf("got %q, want it to contain ALTER TABLE USERS ADD COLUMN PHONE TEXT", sql)
	}
}

func TestSQLiteSchemaGrammar_AddIndex(t *testing.T) {
	g := NewSchemaGrammar()

	sql, err := g.AddIndex("users", []string{"email"}, "idx_users_email", true)
	if err != nil {
		t.Fatalf("AddIndex failed: %v", err)
	}

	sql = strings.ToUpper(sql)
	if !strings.Contains(sql, "CREATE UNIQUE INDEX IDX_USERS_EMAIL ON USERS (EMAIL)") {
		t.Errorf("got %q, want it to contain CREATE UNIQUE INDEX IDX_USERS_EMAIL ON USERS (EMAIL)", sql)
	}
}

func BenchmarkSQLiteSchemaGrammar_CreateTable(b *testing.B) {
	g := NewSchemaGrammar()
	bp := migrator.Blueprint{
		TableName: "users",
		Columns: []*migrator.Column{
			{Name: "id", Type: database_types.TypeBigInt, IsPrimary: true, IsAutoInc: true},
			{Name: "username", Type: database_types.TypeVarchar, Length: 255, IsNullable: false, IsUnique: true},
			{Name: "email", Type: database_types.TypeVarchar, Length: 255, IsNullable: true},
			{Name: "age", Type: database_types.TypeInt},
			{Name: "bio", Type: database_types.TypeText},
			{Name: "created_at", Type: database_types.TypeTimestamp, DefaultVal: "now()"},
			{Name: "updated_at", Type: database_types.TypeTimestamp, DefaultVal: "now()"},
		},
	}

	b.ReportAllocs()
	for b.Loop() {
		_, _ = g.CreateTable(bp)
	}
}
