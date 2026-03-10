package postgres

import (
	"fmt"
	"strings"

	"github.com/vayload/vayload/internal/modules/database/migrator"
)

type PostgresGrammar struct{}

func NewGrammar() *PostgresGrammar {
	return &PostgresGrammar{}
}

func (g *PostgresGrammar) CreateTable(bp migrator.Blueprint) (string, error) {
	var definitions []string

	for _, c := range bp.Columns {
		colType := c.Type
		if c.IsAutoInc {
			if c.Type == "BIGINT" {
				colType = "BIGSERIAL"
			} else {
				colType = "SERIAL"
			}
		} else if c.Type == "VARCHAR" {
			colType = fmt.Sprintf("VARCHAR(%d)", c.Length)
		} else if c.Type == "DECIMAL" {
			colType = fmt.Sprintf("DECIMAL(%d,%d)", c.Precision, c.Scale)
		} else if c.Type == "BOOLEAN" {
			colType = "BOOLEAN"
		}

		def := fmt.Sprintf("\"%s\" %s", c.Name, colType)
		if c.IsPrimary {
			def += " PRIMARY KEY"
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
		def := fmt.Sprintf("CONSTRAINT fk_%s FOREIGN KEY (\"%s\") REFERENCES \"%s\"(\"%s\")",
			fc.Name, fc.Name, fc.ForeignReference.Table, fc.ForeignReference.Column)
		if fc.OnDeleteAction != "" {
			def += fmt.Sprintf(" ON DELETE %s", strings.ToUpper(fc.OnDeleteAction))
		}
		definitions = append(definitions, def)
	}

	query := fmt.Sprintf("CREATE TABLE \"%s\" (%s);", bp.TableName, strings.Join(definitions, ", "))
	return query, nil
}

func (g *PostgresGrammar) DropTable(name string, ifExists bool) (string, error) {
	if ifExists {
		return fmt.Sprintf("DROP TABLE IF EXISTS \"%s\";", name), nil
	}
	return fmt.Sprintf("DROP TABLE \"%s\";", name), nil
}

func (g *PostgresGrammar) AddColumn(table string, column *migrator.Column) (string, error) {
	colType := column.Type
	if column.Type == "VARCHAR" {
		colType = fmt.Sprintf("VARCHAR(%d)", column.Length)
	}

	def := fmt.Sprintf("\"%s\" %s", column.Name, colType)
	if !column.IsNullable {
		def += " NOT NULL"
	}
	if column.DefaultVal != nil {
		def += fmt.Sprintf(" DEFAULT %v", g.formatDefault(column.DefaultVal))
	}

	return fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN %s;", table, def), nil
}

func (g *PostgresGrammar) DropColumn(table string, columnName string) (string, error) {
	return fmt.Sprintf("ALTER TABLE \"%s\" DROP COLUMN \"%s\";", table, columnName), nil
}

func (g *PostgresGrammar) RenameColumn(table string, from string, to string) (string, error) {
	return fmt.Sprintf("ALTER TABLE \"%s\" RENAME COLUMN \"%s\" TO \"%s\";", table, from, to), nil
}

func (g *PostgresGrammar) AddIndex(table string, columns []string, name string, unique bool) (string, error) {
	uniqueStr := ""
	if unique {
		uniqueStr = "UNIQUE "
	}
	quotedCols := make([]string, len(columns))
	for i, col := range columns {
		quotedCols[i] = fmt.Sprintf("\"%s\"", col)
	}
	colList := strings.Join(quotedCols, ", ")
	return fmt.Sprintf("CREATE %sINDEX \"%s\" ON \"%s\" (%s);", uniqueStr, name, table, colList), nil
}

func (g *PostgresGrammar) DropIndex(table string, name string) (string, error) {
	return fmt.Sprintf("DROP INDEX \"%s\";", name), nil
}

func (g *PostgresGrammar) HasTable(table string) (string, error) {
	return fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = '%s';", table), nil
}

func (g *PostgresGrammar) HasColumn(table string, column string) (string, error) {
	return fmt.Sprintf("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = 'public' AND table_name = '%s' AND column_name = '%s';", table, column), nil
}

func (g *PostgresGrammar) formatDefault(val any) string {
	if s, ok := val.(string); ok {
		if s == "now()" {
			return "CURRENT_TIMESTAMP"
		}
		return fmt.Sprintf("'%s'", s)
	}
	return fmt.Sprintf("'%v'", val)
}
