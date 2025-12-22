package d1sql

import (
	"context"
	"database/sql/driver"
)

// Stmt implements driver.Stmt, driver.StmtExecContext, and driver.StmtQueryContext
type Stmt struct {
	conn  *Conn
	query string
	ctx   context.Context
}

// Close closes the statement
func (s *Stmt) Close() error {
	s.conn = nil
	return nil
}

// NumInput returns the number of placeholder parameters
// Returns -1 to indicate the driver doesn't know the number of parameters
func (s *Stmt) NumInput() int {
	return -1
}

// Exec executes a query that doesn't return rows
func (s *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.ExecContext(s.ctx, namedValuesFromValues(args))
}

// ExecContext executes a query with context
func (s *Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if s.conn == nil || s.conn.isClosed() {
		return nil, driver.ErrBadConn
	}

	return s.conn.ExecContext(ctx, s.query, args)
}

// Query executes a query that returns rows
func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.QueryContext(s.ctx, namedValuesFromValues(args))
}

// QueryContext executes a query with context
func (s *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if s.conn == nil || s.conn.isClosed() {
		return nil, driver.ErrBadConn
	}

	return s.conn.QueryContext(ctx, s.query, args)
}

