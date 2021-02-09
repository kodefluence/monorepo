package db

import (
	"database/sql"

	"github.com/codefluence-x/monorepo/exception"
	"github.com/codefluence-x/monorepo/kontext"
)

// An Adapter for golang sql
type Adapter struct {
	db *sql.DB
}

// Adapt adapting golang sql.DB
func Adapt(db *sql.DB) DB {
	return &Adapter{db: db}
}

// Transaction wrap mysql transaction into a bit of simpler way
func (d *Adapter) Transaction(ktx kontext.Context, transactionKey string, f func(tx TX) exception.Exception) exception.Exception {
	return runWithSQLAnalyzer(ktx, "db", "Transaction", func() exception.Exception {
		tx, err := d.db.BeginTx(ktx.Ctx(), &sql.TxOptions{})
		if err != nil {
			return exception.Throw(err)
		}

		adaptedTx := AdaptTXAdapter(tx)
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
func (d *Adapter) ExecContext(ktx kontext.Context, queryKey, query string, args ...interface{}) (Result, exception.Exception) {
	var result sql.Result
	var err error
	var exc exception.Exception

	exc = runWithSQLAnalyzer(ktx, "db", "ExecContext", func() exception.Exception {
		result, err = d.db.ExecContext(ktx.Ctx(), query, args...)
		if err != nil {
			return exception.Throw(err)
		}

		return nil
	})

	return AdaptResult(result), exc
}

// QueryContext wrap sql QueryContext function
func (d *Adapter) QueryContext(ktx kontext.Context, queryKey, query string, args ...interface{}) (Rows, exception.Exception) {
	var rows *sql.Rows
	var err error
	var exc exception.Exception

	exc = runWithSQLAnalyzer(ktx, "db", "QueryContext", func() exception.Exception {
		rows, err = d.db.QueryContext(ktx.Ctx(), query, args...)
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
func (d *Adapter) QueryRowContext(ktx kontext.Context, queryKey, query string, args ...interface{}) Row {
	var row *sql.Row

	_ = runWithSQLAnalyzer(ktx, "db", "QueryRowContext", func() exception.Exception {
		row = d.db.QueryRowContext(ktx.Ctx(), query, args...)
		return nil
	})

	return AdaptRow(row)
}

func runWithSQLAnalyzer(ktx kontext.Context, executionLevel, function string, f func() exception.Exception) exception.Exception {
	if err := f(); err != nil {
		return err
	}
	return nil
}
