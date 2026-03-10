package sqlite3

import (
	"fmt"
	"strings"

	"github.com/vayload/vayload/internal/modules/database/migrator"
)

type SQLiteGrammar struct{}

func NewGrammar() *SQLiteGrammar {
	return &SQLiteGrammar{}
}

func (g *SQLiteGrammar) CreateTable(bp migrator.Blueprint) (string, error) {
	var definitions []string

	for _, c := range bp.Columns {
		def := fmt.Sprintf("%s %s", c.Name, c.Type)
		if c.IsPrimary {
			def += " PRIMARY KEY"
		}
		if c.IsAutoInc && c.Type == "INTEGER" {
			def += " AUTOINCREMENT"
		}
		if !c.IsNullable {
			def += " NOT NULL"
		}
		if c.IsUnique {
			def += " UNIQUE"
		}
		if c.DefaultVal != nil {
			def += fmt.Sprintf(" DEFAULT %v", g.formatDefault(c.DefaultVal))
		}
		definitions = append(definitions, def)
	}

	for _, fc := range bp.ForeignColumns {
		def := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)",
			fc.Name, fc.ForeignReference.Table, fc.ForeignReference.Column)
		if fc.OnDeleteAction != "" {
			def += fmt.Sprintf(" ON DELETE %s", strings.ToUpper(fc.OnDeleteAction))
		}
		definitions = append(definitions, def)
	}

	query := fmt.Sprintf("CREATE TABLE %s (%s);", bp.TableName, strings.Join(definitions, ", "))
	return query, nil
}

func (g *SQLiteGrammar) DropTable(name string, ifExists bool) (string, error) {
	if ifExists {
		return fmt.Sprintf("DROP TABLE IF EXISTS %s;", name), nil
	}
	return fmt.Sprintf("DROP TABLE %s;", name), nil
}

func (g *SQLiteGrammar) AddColumn(table string, column *migrator.Column) (string, error) {
	def := fmt.Sprintf("%s %s", column.Name, column.Type)
	if !column.IsNullable {
		def += " NOT NULL"
	}
	if column.DefaultVal != nil {
		def += fmt.Sprintf(" DEFAULT %v", g.formatDefault(column.DefaultVal))
	}
	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s;", table, def), nil
}

func (g *SQLiteGrammar) DropColumn(table string, columnName string) (string, error) {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", table, columnName), nil
}

func (g *SQLiteGrammar) RenameColumn(table string, from string, to string) (string, error) {
	return fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s;", table, from, to), nil
}

func (g *SQLiteGrammar) AddIndex(table string, columns []string, name string, unique bool) (string, error) {
	uniqueStr := ""
	if unique {
		uniqueStr = "UNIQUE "
	}
	colList := strings.Join(columns, ", ")
	return fmt.Sprintf("CREATE %sINDEX %s ON %s (%s);", uniqueStr, name, table, colList), nil
}

func (g *SQLiteGrammar) DropIndex(table string, name string) (string, error) {
	return fmt.Sprintf("DROP INDEX %s;", name), nil
}

func (g *SQLiteGrammar) HasTable(table string) (string, error) {
	return fmt.Sprintf("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='%s';", table), nil
}

func (g *SQLiteGrammar) HasColumn(table string, column string) (string, error) {
	return fmt.Sprintf("SELECT COUNT(*) FROM pragma_table_info('%s') WHERE name='%s';", table, column), nil
}

func (g *SQLiteGrammar) formatDefault(val any) string {
	if s, ok := val.(string); ok {
		if s == "now()" {
			return "CURRENT_TIMESTAMP"
		}
		return fmt.Sprintf("'%s'", s)
	}
	return fmt.Sprintf("%v", val)
}
