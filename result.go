package d1sql

// Result implements driver.Result
type Result struct {
	lastInsertID int64
	rowsAffected int64
}

// LastInsertId returns the database's auto-generated ID after an INSERT
func (r *Result) LastInsertId() (int64, error) {
	return r.lastInsertID, nil
}

// RowsAffected returns the number of rows affected by the query
func (r *Result) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

