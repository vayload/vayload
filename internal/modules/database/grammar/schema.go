package grammar

import "github.com/vayload/vayload/internal/modules/database/migrator"

// For manage database schema (drop, create tables, etc.)
type SchemaGrammar interface {
	CreateTable(bp migrator.Blueprint) (string, error)
	DropTable(name string, ifExists bool) (string, error)

	// Modificación de columnas
	AddColumn(table string, column *migrator.Column) (string, error)
	DropColumn(table string, columnName string) (string, error)
	RenameColumn(table string, from string, to string) (string, error)

	// Índices
	AddIndex(table string, columns []string, name string, unique bool) (string, error)
	DropIndex(table string, name string) (string, error)

	// Utilidades
	HasTable(table string) (string, error)
	HasColumn(table string, column string) (string, error)
}
