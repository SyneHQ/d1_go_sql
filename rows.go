package d1sql

import (
	"database/sql/driver"
	"io"
)

// Rows implements driver.Rows
type Rows struct {
	columns []string
	rows    [][]interface{}
	index   int
}

// Columns returns the names of the columns
func (r *Rows) Columns() []string {
	return r.columns
}

// Close closes the rows iterator
func (r *Rows) Close() error {
	r.rows = nil
	return nil
}

// Next is called to populate the next row of data into the provided slice
func (r *Rows) Next(dest []driver.Value) error {
	r.index++
	
	if r.index >= len(r.rows) {
		return io.EOF
	}

	row := r.rows[r.index]
	
	if len(dest) != len(row) {
		return driver.ErrSkip
	}

	for i, val := range row {
		dest[i] = convertValue(val)
	}

	return nil
}

