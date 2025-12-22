package d1sql

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/d1"
)

// Conn implements driver.Conn, driver.ConnPrepareContext, driver.ExecerContext, and driver.QueryerContext
type Conn struct {
	client       *cloudflare.Client
	config       *Config
	ctx          context.Context
	inTx         bool
	txStatements []string
}

// Prepare creates a prepared statement for later queries or executions
func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	return c.PrepareContext(c.ctx, query)
}

// PrepareContext creates a prepared statement with context
func (c *Conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if c.isClosed() {
		return nil, driver.ErrBadConn
	}

	return &Stmt{
		conn:  c,
		query: query,
		ctx:   ctx,
	}, nil
}

// Close closes the connection
func (c *Conn) Close() error {
	if c.isClosed() {
		return nil
	}

	// If in transaction, rollback
	if c.inTx {
		tx := &Tx{conn: c}
		_ = tx.Rollback()
	}

	c.client = nil
	return nil
}

// Begin starts a transaction
func (c *Conn) Begin() (driver.Tx, error) {
	return c.BeginTx(c.ctx, driver.TxOptions{})
}

// BeginTx starts a transaction with context and options
func (c *Conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if c.isClosed() {
		return nil, driver.ErrBadConn
	}

	if c.inTx {
		return nil, errors.New("already in transaction")
	}

	// D1 doesn't support transaction isolation levels
	if opts.Isolation != driver.IsolationLevel(0) {
		return nil, errors.New("D1 does not support transaction isolation levels")
	}

	c.inTx = true
	c.txStatements = make([]string, 0)

	return &Tx{conn: c}, nil
}

// ExecContext executes a query without returning rows
func (c *Conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if c.isClosed() {
		return nil, driver.ErrBadConn
	}

	// Convert named values to regular values
	values, err := namedValuesToValues(args)
	if err != nil {
		return nil, err
	}

	// If in transaction, queue the statement
	if c.inTx {
		boundQuery, err := bindParameters(query, values)
		if err != nil {
			return nil, err
		}
		c.txStatements = append(c.txStatements, boundQuery)
		// Return a placeholder result for transaction
		return &Result{rowsAffected: 0}, nil
	}

	return c.execQuery(ctx, query, values)
}

// QueryContext executes a query that returns rows
func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if c.isClosed() {
		return nil, driver.ErrBadConn
	}

	if c.inTx {
		return nil, errors.New("cannot query within a transaction; use Commit first")
	}

	// Convert named values to regular values
	values, err := namedValuesToValues(args)
	if err != nil {
		return nil, err
	}

	return c.queryRows(ctx, query, values)
}

// execQuery executes a single query and returns the result
func (c *Conn) execQuery(ctx context.Context, query string, args []driver.Value) (driver.Result, error) {
	boundQuery, err := bindParameters(query, args)
	if err != nil {
		return nil, err
	}

	// Execute query using Cloudflare SDK
	params := d1.DatabaseQueryParams{
		AccountID: cloudflare.F(c.config.AccountID),
		Body: d1.DatabaseQueryParamsBodyD1SingleQuery{
			Sql: cloudflare.F(boundQuery),
		},
	}

	// Use QueryAutoPaging to get an autopager
	iter := c.client.D1.Database.QueryAutoPaging(ctx, c.config.DatabaseID, params)

	// Iterate through results to get metadata
	var rowsAffected int64
	var lastInsertID int64

	for iter.Next() {
		queryResult := iter.Current()
		if queryResult.Meta.RowsWritten > 0 {
			rowsAffected = int64(queryResult.Meta.RowsWritten)
		}
		if queryResult.Meta.LastRowID > 0 {
			lastInsertID = int64(queryResult.Meta.LastRowID)
		}
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error reading query results: %w", err)
	}

	return &Result{
		lastInsertID: lastInsertID,
		rowsAffected: rowsAffected,
	}, nil
}

// queryRows executes a query and returns rows
func (c *Conn) queryRows(ctx context.Context, query string, args []driver.Value) (driver.Rows, error) {
	boundQuery, err := bindParameters(query, args)
	if err != nil {
		return nil, err
	}

	// Execute query using Cloudflare SDK
	params := d1.DatabaseQueryParams{
		AccountID: cloudflare.F(c.config.AccountID),
		Body: d1.DatabaseQueryParamsBodyD1SingleQuery{
			Sql: cloudflare.F(boundQuery),
		},
	}

	// Use QueryAutoPaging to get an autopager
	iter := c.client.D1.Database.QueryAutoPaging(ctx, c.config.DatabaseID, params)

	// Collect all results
	var columns []string
	var rows [][]interface{}

	for iter.Next() {
		queryResult := iter.Current()

		// Get columns from metadata - D1.QueryResult has []interface{} for Results
		// Each result in Results is a map[string]interface{} representing a row

		// Append results
		for _, row := range queryResult.Results {
			// Convert row to []interface{}
			if rowMap, ok := row.(map[string]interface{}); ok {
				// If it's a map, extract values in order
				if len(columns) == 0 {
					// First row - extract column names in a consistent order
					for key := range rowMap {
						columns = append(columns, key)
					}
				}
				rowValues := make([]interface{}, len(columns))
				for i, col := range columns {
					rowValues[i] = rowMap[col]
				}
				rows = append(rows, rowValues)
			}
		}
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error reading query results: %w", err)
	}

	return &Rows{
		columns: columns,
		rows:    rows,
		index:   -1,
	}, nil
}

// isClosed checks if the connection is closed
func (c *Conn) isClosed() bool {
	return c.client == nil
}

// Tx implements driver.Tx
type Tx struct {
	conn *Conn
}

// Commit commits the transaction
func (t *Tx) Commit() error {
	if !t.conn.inTx {
		return errors.New("not in transaction")
	}

	defer func() {
		t.conn.inTx = false
		t.conn.txStatements = nil
	}()

	if len(t.conn.txStatements) == 0 {
		return nil
	}

	// Execute all statements as a batch - join with semicolon
	batchQuery := ""
	for i, stmt := range t.conn.txStatements {
		if i > 0 {
			batchQuery += "; "
		}
		batchQuery += stmt
	}

	// Execute batch query
	params := d1.DatabaseQueryParams{
		AccountID: cloudflare.F(t.conn.config.AccountID),
		Body: d1.DatabaseQueryParamsBodyD1SingleQuery{
			Sql: cloudflare.F(batchQuery),
		},
	}

	// Execute and consume the iterator
	iter := t.conn.client.D1.Database.QueryAutoPaging(t.conn.ctx, t.conn.config.DatabaseID, params)
	for iter.Next() {
		// Just consume the results
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

// Rollback rolls back the transaction
func (t *Tx) Rollback() error {
	if !t.conn.inTx {
		return errors.New("not in transaction")
	}

	t.conn.inTx = false
	t.conn.txStatements = nil
	return nil
}
