package drivers

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vayload/vayload/internal/modules/database/builder"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/modules/database/grammar"
)

type BaseConnection struct {
	db      *sqlx.DB
	ctx     context.Context
	cancel  context.CancelFunc
	driver  connection.DatabaseDriver
	grammar grammar.QueryGrammar
}

func NewBaseConnection(
	db *sqlx.DB,
	ctx context.Context,
	cancel context.CancelFunc,
	driver connection.DatabaseDriver,
	grammar grammar.QueryGrammar,
) BaseConnection {
	return BaseConnection{
		db:      db,
		ctx:     ctx,
		cancel:  cancel,
		driver:  driver,
		grammar: grammar,
	}
}

func (c *BaseConnection) GetDriverName() connection.DatabaseDriver {
	return c.driver
}

func (c *BaseConnection) From(table string) connection.QueryBuilder {
	return builder.NewQueryBuilder(c, c.grammar, table)
}

func (c *BaseConnection) Select(ctx context.Context, dest any, query string, args ...any) error {
	return c.db.SelectContext(ctx, dest, query, args...)
}

func (c *BaseConnection) SelectOne(ctx context.Context, dest any, query string, args ...any) error {
	return c.db.GetContext(ctx, dest, query, args...)
}

func (c *BaseConnection) Scan(ctx context.Context, query string, args []any, dest ...any) error {
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	return rows.Scan(dest...)
}

func (c *BaseConnection) Cursor(ctx context.Context, query string, args []any) (connection.Cursor, error) {
	rows, err := c.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return NewBaseCursor(rows), nil
}

func (c *BaseConnection) Prepared(ctx context.Context, query string, binding ...any) error {
	stmt, err := c.db.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, binding...)
	if err != nil {
		return fmt.Errorf("failed to execute prepared statement: %w", err)
	}

	return nil
}

func (c *BaseConnection) Unprepared(ctx context.Context, query string) error {
	_, err := c.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute unprepared statement: %w", err)
	}

	return nil
}

func (c *BaseConnection) Close() error {
	if c.db != nil {
		c.db.Close()
	}

	c.cancel()
	return nil
}

func (c *BaseConnection) Context() context.Context {
	return c.ctx
}

// Cursor implementation

type BaseCursor struct {
	rows *sqlx.Rows
}

func NewBaseCursor(rows *sqlx.Rows) *BaseCursor {
	return &BaseCursor{rows: rows}
}

func (c *BaseCursor) Next() bool {
	return c.rows.Next()
}

func (c *BaseCursor) Scan(dest ...any) error {
	return c.rows.Scan(dest...)
}

func (c *BaseCursor) Close() error {
	return c.rows.Close()
}
