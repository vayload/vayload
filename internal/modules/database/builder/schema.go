package builder

import (
	"context"
	"fmt"

	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/modules/database/grammar"
	"github.com/vayload/vayload/internal/modules/database/migrator"
)

type SchemaBuilder struct {
	grammar    grammar.SchemaGrammar
	connection connection.DatabaseConnection
}

func NewSchemaBuilder(grammar grammar.SchemaGrammar, connection connection.DatabaseConnection) *SchemaBuilder {
	return &SchemaBuilder{grammar: grammar, connection: connection}
}

func (s *SchemaBuilder) Create(table string, callback func(*migrator.Blueprint)) error {
	bp := migrator.NewBlueprint(table)
	callback(bp)

	query, err := s.grammar.CreateTable(*bp)
	if err != nil {
		return err
	}

	return s.connection.Unprepared(context.Background(), query)
}

func (s *SchemaBuilder) Table(table string, callback func(*migrator.Blueprint)) error {
	bp := migrator.NewBlueprint(table)
	bp.IsAltering = true
	callback(bp)

	for _, cmd := range bp.Commands {
		var query string
		var err error

		switch cmd.Type {
		case migrator.AddColumnCommand:
			query, err = s.grammar.AddColumn(table, cmd.Column)
		case migrator.DropColumnCommand:
			query, err = s.grammar.DropColumn(table, cmd.Name)
		case migrator.RenameColumnCommand:
			query, err = s.grammar.RenameColumn(table, cmd.From, cmd.To)
		case migrator.AddIndexCommand:
			query, err = s.grammar.AddIndex(table, cmd.Columns, cmd.Name, cmd.IsUnique)
		case migrator.DropIndexCommand:
			query, err = s.grammar.DropIndex(table, cmd.Name)
		default:
			continue
		}

		if err != nil {
			return err
		}

		if err := s.connection.Unprepared(context.Background(), query); err != nil {
			return err
		}
	}

	return nil
}

func (s *SchemaBuilder) Drop(table string) error {
	query, err := s.grammar.DropTable(table, false)
	if err != nil {
		return err
	}

	return s.connection.Unprepared(context.Background(), query)
}

func (s *SchemaBuilder) DropIfExists(table string) error {
	query, err := s.grammar.DropTable(table, true)
	if err != nil {
		return err
	}

	return s.connection.Unprepared(context.Background(), query)
}

func (s *SchemaBuilder) AddIndex(table string, columns []string, name string, unique bool) error {
	if name == "" {
		name = fmt.Sprintf("idx_%s_%s", table, columns[0])
	}

	query, err := s.grammar.AddIndex(table, columns, name, unique)
	if err != nil {
		return err
	}

	return s.connection.Unprepared(context.Background(), query)
}

func (s *SchemaBuilder) DropIndex(table string, name string) error {
	query, err := s.grammar.DropIndex(table, name)
	if err != nil {
		return err
	}

	return s.connection.Unprepared(context.Background(), query)
}

func (s *SchemaBuilder) AddColumn(table string, name string, colType string) error {
	col := &migrator.Column{
		Name: name,
		Type: colType,
	}

	query, err := s.grammar.AddColumn(table, col)
	if err != nil {
		return err
	}

	return s.connection.Unprepared(context.Background(), query)
}

func (s *SchemaBuilder) DropColumn(table string, name string) error {
	query, err := s.grammar.DropColumn(table, name)
	if err != nil {
		return err
	}

	return s.connection.Unprepared(context.Background(), query)
}

func (s *SchemaBuilder) RenameColumn(table string, from string, to string) error {
	query, err := s.grammar.RenameColumn(table, from, to)
	if err != nil {
		return err
	}

	return s.connection.Unprepared(context.Background(), query)
}

func (s *SchemaBuilder) HasTable(table string) (bool, error) {
	query, err := s.grammar.HasTable(table)
	if err != nil {
		return false, err
	}

	return s.executeExists(query)
}

func (s *SchemaBuilder) HasColumn(table string, column string) (bool, error) {
	query, err := s.grammar.HasColumn(table, column)
	if err != nil {
		return false, err
	}

	return s.executeExists(query)
}

func (s *SchemaBuilder) executeExists(_ string) (bool, error) {
	return false, fmt.Errorf("not implemented: needs DB access")
}
