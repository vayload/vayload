package drivers

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vayload/vayload/internal/services/database/builder"
	"github.com/vayload/vayload/internal/services/database/connection"
)

type BaseConnection struct {
	db     *sqlx.DB
	ctx    context.Context
	cancel context.CancelFunc
	driver connection.DatabaseDriver
}

func NewBaseConnection(db *sqlx.DB, ctx context.Context, cancel context.CancelFunc, driver connection.DatabaseDriver) BaseConnection {
	return BaseConnection{
		db:     db,
		ctx:    ctx,
		cancel: cancel,
		driver: driver,
	}
}

func (c *BaseConnection) GetDriverName() connection.DatabaseDriver {
	return c.driver
}

func (c *BaseConnection) From(table string) connection.QueryBuilder {
	return builder.NewQueryBuilder(c, table, c.driver)
}

func (c *BaseConnection) Select(ctx context.Context, dest any, query string, args ...any) error {
	return c.db.SelectContext(ctx, dest, query, args...)
}

func (c *BaseConnection) SelectOne(ctx context.Context, dest any, query string, args ...any) error {
	return c.db.GetContext(ctx, dest, query, args...)
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
