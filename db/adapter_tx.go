package db

import (
	"database/sql"

	"github.com/codefluence-x/monorepo/exception"
	"github.com/codefluence-x/monorepo/kontext"
)

// A TXAdapter adapater for golang sql
type TXAdapter struct {
	tx *sql.Tx
}

// AdaptTXAdapter do adapting mysql transaction
func AdaptTXAdapter(tx *sql.Tx) *TXAdapter {
	return &TXAdapter{tx: tx}
}

// ExecContext wrap sql ExecContext function
func (t *TXAdapter) ExecContext(ctx kontext.Context, queryKey, query string, args ...interface{}) (Result, exception.Exception) {
	var result sql.Result
	var err error
	var exc exception.Exception

	exc = runWithSQLAnalyzer(ctx, "tx", "ExecContext", func() exception.Exception {
		result, err = t.tx.ExecContext(ctx.Ctx(), query, args...)
		if err != nil {
			return exception.Throw(err)
		}

		return nil
	})

	return AdaptResult(result), exc
}

// QueryContext wrap sql QueryContext function
func (t *TXAdapter) QueryContext(ctx kontext.Context, queryKey, query string, args ...interface{}) (Rows, exception.Exception) {
	var rows *sql.Rows
	var err error
	var exc exception.Exception

	exc = runWithSQLAnalyzer(ctx, "tx", "QueryContext", func() exception.Exception {
		rows, err = t.tx.QueryContext(ctx.Ctx(), query, args...)
		if err == sql.ErrNoRows {
			return exception.Throw(err, exception.WithType(exception.NotFound))
		} else if err != nil {
			return exception.Throw(err)
		}

		return nil
	})

	return AdaptRows(rows), exc
}

// QueryRowContext wrap sql QueryRowContext function
func (t *TXAdapter) QueryRowContext(ctx kontext.Context, queryKey, query string, args ...interface{}) Row {
	var row *sql.Row

	_ = runWithSQLAnalyzer(ctx, "tx", "QueryRowContext", func() exception.Exception {
		row = t.tx.QueryRowContext(ctx.Ctx(), query, args...)
		return nil
	})

	return AdaptRow(row)
}
