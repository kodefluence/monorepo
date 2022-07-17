package db

import (
	"database/sql"

	"github.com/kodefluence/monorepo/exception"
)

// RowsAdapter wrap default sql.Rows struct
type RowsAdapter struct {
	*sql.Rows
}

// AdaptRows adapting sql.Rows into adapter.Rows
func AdaptRows(rows *sql.Rows) Rows {
	return &RowsAdapter{Rows: rows}
}

// Close rows
func (r *RowsAdapter) Close() exception.Exception {
	if err := r.Rows.Close(); err != nil {
		return exception.Throw(err)
	}

	return nil
}

// Columns return rows column
func (r *RowsAdapter) Columns() ([]string, exception.Exception) {
	var columns []string
	var err error

	columns, err = r.Rows.Columns()
	if err != nil {
		return columns, exception.Throw(err)
	}

	return columns, nil
}

// Err return rows error
func (r *RowsAdapter) Err() exception.Exception {
	if err := r.Rows.Err(); err != nil {
		return exception.Throw(err)
	}

	return nil
}

// Scan row
func (r *RowsAdapter) Scan(dest ...interface{}) exception.Exception {
	if err := r.Rows.Scan(dest...); err != nil {
		return exception.Throw(err)
	}

	return nil
}
