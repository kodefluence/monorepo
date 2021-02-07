package adapter

import (
	"database/sql"

	"github.com/codefluence-x/monorepo/exception"
)

// RowAdapter wrap single sql row
type RowAdapter struct {
	sqlrow *sql.Row
}

// AdaptRow wrap provider row
func AdaptRow(sqlrow *sql.Row) Row {
	return &RowAdapter{sqlrow: sqlrow}
}

// Scan warp default row scan function
func (r *RowAdapter) Scan(dest ...interface{}) exception.Exception {
	err := r.sqlrow.Scan(dest...)
	if err == sql.ErrNoRows {
		return exception.Throw(err, exception.WithType(exception.NotFound))
	} else if err != nil {
		return exception.Throw(err)
	}

	return nil
}
