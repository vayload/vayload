package mysql

import (
	"fmt"
	"strings"

	"github.com/vayload/vayload/internal/modules/database/migrator"
)

type MySQLGrammar struct{}

func NewGrammar() *MySQLGrammar {
	return &MySQLGrammar{}
}

func (g *MySQLGrammar) CreateTable(bp migrator.Blueprint) (string, error) {
	var definitions []string

	for _, c := range bp.Columns {
		typeName := c.Type
		if c.Type == "VARCHAR" {
			typeName = fmt.Sprintf("VARCHAR(%d)", c.Length)
		} else if c.Type == "DECIMAL" {
			typeName = fmt.Sprintf("DECIMAL(%d,%d)", c.Precision, c.Scale)
		}

		def := fmt.Sprintf("`%s` %s", c.Name, typeName)
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

func (g *MySQLGrammar) DropTable(name string, ifExists bool) (string, error) {
	if ifExists {
		return fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", name), nil
	}
	return fmt.Sprintf("DROP TABLE `%s`;", name), nil
}

func (g *MySQLGrammar) AddColumn(table string, column *migrator.Column) (string, error) {
	typeName := column.Type
	if column.Type == "VARCHAR" {
		typeName = fmt.Sprintf("VARCHAR(%d)", column.Length)
	}

	def := fmt.Sprintf("`%s` %s", column.Name, typeName)
	if !column.IsNullable {
		def += " NOT NULL"
	}
	if column.DefaultVal != nil {
		def += fmt.Sprintf(" DEFAULT %v", g.formatDefault(column.DefaultVal))
	}

	return fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN %s;", table, def), nil
}

func (g *MySQLGrammar) DropColumn(table string, columnName string) (string, error) {
	return fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", table, columnName), nil
}

func (g *MySQLGrammar) RenameColumn(table string, from string, to string) (string, error) {
	return fmt.Sprintf("ALTER TABLE `%s` RENAME COLUMN `%s` TO `%s`;", table, from, to), nil
}

func (g *MySQLGrammar) AddIndex(table string, columns []string, name string, unique bool) (string, error) {
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

func (g *MySQLGrammar) DropIndex(table string, name string) (string, error) {
	return fmt.Sprintf("DROP INDEX `%s` ON `%s`;", name, table), nil
}

func (g *MySQLGrammar) HasTable(table string) (string, error) {
	return fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = '%s';", table), nil
}

func (g *MySQLGrammar) HasColumn(table string, column string) (string, error) {
	return fmt.Sprintf("SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = '%s' AND column_name = '%s';", table, column), nil
}

func (g *MySQLGrammar) formatDefault(val any) string {
	if s, ok := val.(string); ok {
		if s == "now()" {
			return "CURRENT_TIMESTAMP"
		}
		return fmt.Sprintf("'%s'", s)
	}
	return fmt.Sprintf("'%v'", val)
}
