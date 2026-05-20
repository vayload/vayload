package mysql

import (
	"fmt"
	"strings"

	"github.com/vayload/vayload/internal/modules/database/migrator"
	database_types "github.com/vayload/vayload/internal/modules/database/types"
)

type MySQLSchemaGrammar struct{}

func NewSchemaGrammar() *MySQLSchemaGrammar {
	return &MySQLSchemaGrammar{}
}

func (g *MySQLSchemaGrammar) CreateTable(bp migrator.Blueprint) (string, error) {
	var definitions []string

	for _, c := range bp.Columns {
		def := fmt.Sprintf("`%s` %s", c.Name, CompileColType(c.Name, database_types.DataType{
			Kind:      database_types.DataTypeKind(c.Type),
			Length:    c.Length,
			Precision: c.Precision,
			Scale:     c.Scale,
			Nullable:  c.IsNullable,
		}))

		if c.IsUnsigned {
			def += " UNSIGNED"
		}
		if !c.IsNullable {
			def += " NOT NULL"
		}
		if c.DefaultVal != nil {
			def += fmt.Sprintf(" DEFAULT %v", g.formatDefault(c.DefaultVal))
		}
		if c.IsAutoInc {
			def += " AUTO_INCREMENT"
		}
		if c.IsPrimary {
			def += " PRIMARY KEY"
		}
		if c.IsUnique {
			def += " UNIQUE"
		}

		definitions = append(definitions, def)
	}

	for _, fc := range bp.ForeignColumns {
		def := fmt.Sprintf("FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`)",
			fc.Name, fc.ForeignReference.Table, fc.ForeignReference.Column)
		if fc.OnDeleteAction != "" {
			def += fmt.Sprintf(" ON DELETE %s", strings.ToUpper(fc.OnDeleteAction))
		}
		definitions = append(definitions, def)
	}

	query := fmt.Sprintf("CREATE TABLE `%s` (%s) ENGINE=InnoDB;", bp.TableName, strings.Join(definitions, ", "))
	return query, nil
}

func (g *MySQLSchemaGrammar) DropTable(name string, ifExists bool) (string, error) {
	if ifExists {
		return fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", name), nil
	}
	return fmt.Sprintf("DROP TABLE `%s`;", name), nil
}

func (g *MySQLSchemaGrammar) AddColumn(table string, column *migrator.Column) (string, error) {
	def := fmt.Sprintf("`%s` %s", column.Name, CompileColType(column.Name, database_types.DataType{
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

	return fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", table, def), nil
}

func (g *MySQLSchemaGrammar) DropColumn(table string, columnName string) (string, error) {
	return fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", table, columnName), nil
}

func (g *MySQLSchemaGrammar) RenameColumn(table string, from string, to string) (string, error) {
	return fmt.Sprintf("ALTER TABLE `%s` RENAME COLUMN `%s` TO `%s`;", table, from, to), nil
}

func (g *MySQLSchemaGrammar) AddIndex(table string, columns []string, name string, unique bool) (string, error) {
	uniqueStr := ""
	if unique {
		uniqueStr = "UNIQUE "
	}
	quotedCols := make([]string, len(columns))
	for i, col := range columns {
		quotedCols[i] = fmt.Sprintf("`%s`", col)
	}
	colList := strings.Join(quotedCols, ", ")
	return fmt.Sprintf("CREATE %sINDEX `%s` ON `%s` (%s);", uniqueStr, name, table, colList), nil
}

func (g *MySQLSchemaGrammar) DropIndex(table string, name string) (string, error) {
	return fmt.Sprintf("DROP INDEX `%s` ON `%s`;", name, table), nil
}

func (g *MySQLSchemaGrammar) HasTable(table string) (string, error) {
	return fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = '%s';", table), nil
}

func (g *MySQLSchemaGrammar) HasColumn(table string, column string) (string, error) {
	return fmt.Sprintf("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = '%s' AND column_name = '%s';", table, column), nil
}

func (g *MySQLSchemaGrammar) formatDefault(val any) string {
	if s, ok := val.(string); ok {
		if s == "now()" {
			return "CURRENT_TIMESTAMP"
		}
		return fmt.Sprintf("'%s'", s)
	}
	return fmt.Sprintf("'%v'", val)
}
