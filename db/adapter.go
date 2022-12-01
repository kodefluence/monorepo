package db

import (
	"database/sql"

	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
)

// An Adapter for golang sql
type Adapter struct {
	db *sql.DB
}

// Adapt adapting golang sql.DB
func Adapt(db *sql.DB) DB {
	return &Adapter{db: db}
}

// Ping wrap sql Ping function
func (a *Adapter) Ping(ktx kontext.Context) exception.Exception {
	if err := a.db.Ping(); err != nil {
		return exception.Throw(err)
	}

	return nil
}

// Transaction wrap mysql transaction into a bit of simpler way
func (a *Adapter) Transaction(ktx kontext.Context, transactionKey string, f func(tx TX) exception.Exception) exception.Exception {
	return runWithSQLAnalyzer(ktx, "db", "Transaction", func() exception.Exception {
		tx, err := a.db.BeginTx(ktx.Ctx(), &sql.TxOptions{})
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
func (a *Adapter) ExecContext(ktx kontext.Context, queryKey, query string, args ...interface{}) (Result, exception.Exception) {
	var result sql.Result
	var err error
	var exc exception.Exception

	exc = runWithSQLAnalyzer(ktx, "db", "ExecContext", func() exception.Exception {
		result, err = a.db.ExecContext(ktx.Ctx(), query, args...)
		if err != nil {
			return exception.Throw(err)
		}

		return nil
	})

	return AdaptResult(result), exc
}

// QueryContext wrap sql QueryContext function
func (a *Adapter) QueryContext(ktx kontext.Context, queryKey, query string, args ...interface{}) (Rows, exception.Exception) {
	var rows *sql.Rows
	var err error
	var exc exception.Exception

	exc = runWithSQLAnalyzer(ktx, "db", "QueryContext", func() exception.Exception {
		rows, err = a.db.QueryContext(ktx.Ctx(), query, args...)
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
func (a *Adapter) QueryRowContext(ktx kontext.Context, queryKey, query string, args ...interface{}) Row {
	var row *sql.Row

	_ = runWithSQLAnalyzer(ktx, "db", "QueryRowContext", func() exception.Exception {
		row = a.db.QueryRowContext(ktx.Ctx(), query, args...)
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

// Eject sql.DB out of db adapter
func (a *Adapter) Eject() *sql.DB {
	return a.db
}
