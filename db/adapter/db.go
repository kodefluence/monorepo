package adapter

import (
	"database/sql"

	"github.com/codefluence-x/monorepo/exception"
	"github.com/codefluence-x/monorepo/kontext"
)

// A DBAdapter for golang sql
type DBAdapter struct {
	db *sql.DB
}

// AdaptDB adapting golang sql.DB
func AdaptDB(db *sql.DB) DB {
	return &DBAdapter{db: db}
}

// Transaction wrap mysql transaction into a bit of simpler way
func (d *DBAdapter) Transaction(ctx kontext.Context, transactionKey string, f func(tx TX) exception.Exception) exception.Exception {
	return runWithSQLAnalyzer(ctx, "db", "Transaction", func() exception.Exception {
		tx, err := d.db.BeginTx(ctx.Ctx(), &sql.TxOptions{})
		if err != nil {
			return exception.Throw(err)
		}

		adaptedTx := &TXAdapter{tx: tx}
		if err := f(adaptedTx); err != nil {
			_ = tx.Rollback()
			return exception.Throw(err)
		}

		if err := tx.Commit(); err != nil {
			_ = tx.Rollback()
			return exception.Throw(err)
		}

		return nil
	})
}

// ExecContext wrap sql ExecContext function
func (d *DBAdapter) ExecContext(ctx kontext.Context, queryKey, query string, args ...interface{}) (Result, exception.Exception) {
	var result sql.Result
	var err error
	var exc exception.Exception

	exc = runWithSQLAnalyzer(ctx, "db", "ExecContext", func() exception.Exception {
		result, err = d.db.ExecContext(ctx.Ctx(), query, args...)
		if err != nil {
			return exception.Throw(err)
		}

		return nil
	})

	return AdaptResult(result), exc
}

// QueryContext wrap sql QueryContext function
func (d *DBAdapter) QueryContext(ctx kontext.Context, queryKey, query string, args ...interface{}) (Rows, exception.Exception) {
	var rows *sql.Rows
	var err error
	var exc exception.Exception

	exc = runWithSQLAnalyzer(ctx, "db", "QueryContext", func() exception.Exception {
		rows, err = d.db.QueryContext(ctx.Ctx(), query, args...)
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
func (d *DBAdapter) QueryRowContext(ctx kontext.Context, queryKey, query string, args ...interface{}) Row {
	var row *sql.Row

	_ = runWithSQLAnalyzer(ctx, "db", "QueryRowContext", func() exception.Exception {
		row = d.db.QueryRowContext(ctx.Ctx(), query, args...)
		return nil
	})

	return AdaptRow(row)
}

func runWithSQLAnalyzer(ctx kontext.Context, executionLevel, function string, f func() exception.Exception) exception.Exception {
	if err := f(); err != nil {
		return err
	}
	return nil
}
