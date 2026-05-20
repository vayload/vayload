package postgres

import (
	"fmt"
	"strings"

	"github.com/vayload/vayload/internal/modules/database/migrator"
	database_types "github.com/vayload/vayload/internal/modules/database/types"
)

type PostgresSchemaGrammar struct{}

func NewSchemaGrammar() *PostgresSchemaGrammar {
	return &PostgresSchemaGrammar{}
}

func (g *PostgresSchemaGrammar) CreateTable(bp migrator.Blueprint) (string, error) {
	var definitions []string

	for _, c := range bp.Columns {
		colType := c.Type
		if c.IsAutoInc {
			if c.Type == database_types.TypeBigInt {
				colType = database_types.TypeBigSerial
			} else {
				colType = database_types.TypeSerial
			}
		}

		def := fmt.Sprintf("\"%s\" %s", c.Name, CompileColType(c.Name, database_types.DataType{
			Kind:      database_types.DataTypeKind(colType),
			Length:    c.Length,
			Precision: c.Precision,
			Scale:     c.Scale,
			Nullable:  c.IsNullable,
		}))

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

func (g *PostgresSchemaGrammar) DropTable(name string, ifExists bool) (string, error) {
	if ifExists {
		return fmt.Sprintf("DROP TABLE IF EXISTS \"%s\";", name), nil
	}
	return fmt.Sprintf("DROP TABLE \"%s\";", name), nil
}

func (g *PostgresSchemaGrammar) AddColumn(table string, column *migrator.Column) (string, error) {
	def := fmt.Sprintf("\"%s\" %s", column.Name, CompileColType(column.Name, database_types.DataType{
		Kind:      database_types.DataTypeKind(column.Type),
		Length:    column.Length,
		Precision: column.Precision,
		Scale:     column.Scale,
		Nullable:  column.IsNullable,
	}))

	if !column.IsNullable {
		def += " NOT NULL"
	}
	if column.DefaultVal != nil {
		def += fmt.Sprintf(" DEFAULT %v", g.formatDefault(column.DefaultVal))
	}

	return fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN %s;", table, def), nil
}

func (g *PostgresSchemaGrammar) DropColumn(table string, columnName string) (string, error) {
	return fmt.Sprintf("ALTER TABLE \"%s\" DROP COLUMN \"%s\";", table, columnName), nil
}

func (g *PostgresSchemaGrammar) RenameColumn(table string, from string, to string) (string, error) {
	return fmt.Sprintf("ALTER TABLE \"%s\" RENAME COLUMN \"%s\" TO \"%s\";", table, from, to), nil
}

func (g *PostgresSchemaGrammar) AddIndex(table string, columns []string, name string, unique bool) (string, error) {
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

func (g *PostgresSchemaGrammar) DropIndex(table string, name string) (string, error) {
	return fmt.Sprintf("DROP INDEX \"%s\";", name), nil
}

func (g *PostgresSchemaGrammar) HasTable(table string) (string, error) {
	return fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = '%s';", table), nil
}

func (g *PostgresSchemaGrammar) HasColumn(table string, column string) (string, error) {
	return fmt.Sprintf("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = 'public' AND table_name = '%s' AND column_name = '%s';", table, column), nil
}

func (g *PostgresSchemaGrammar) formatDefault(val any) string {
	if s, ok := val.(string); ok {
		if s == "now()" {
			return "CURRENT_TIMESTAMP"
		}
		return fmt.Sprintf("'%s'", s)
	}
	return fmt.Sprintf("'%v'", val)
}
