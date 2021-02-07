package db

import (
	"database/sql"
	"fmt"

	"github.com/codefluence-x/monorepo/exception"
	"github.com/codefluence-x/monorepo/kontext"

	// Import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// FabricateMySQL will fabricate mysql connection and wrap it into SQL interfaces
func FabricateMySQL(config Config) (DB, exception.Exception) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&interpolateParams=true", config.Username, config.Password, config.Host, config.Port, config.Name))
	if err != nil {
		return nil, exception.Throw(err)
	}

	return Adapt(db), nil
}

// DB is database interface wrapper for *sql.DB
type DB interface {
	Transaction(ctx kontext.Context, transactionKey string, f func(tx TX) exception.Exception) exception.Exception
	ExecContext(ctx kontext.Context, queryKey, query string, args ...interface{}) (Result, exception.Exception)
	QueryContext(ctx kontext.Context, queryKey, query string, args ...interface{}) (Rows, exception.Exception)
	QueryRowContext(ctx kontext.Context, queryKey, query string, args ...interface{}) Row
}

// TX is database transaction
type TX interface {
	ExecContext(ctx kontext.Context, queryKey, query string, args ...interface{}) (Result, exception.Exception)
	QueryContext(ctx kontext.Context, queryKey, query string, args ...interface{}) (Rows, exception.Exception)
	QueryRowContext(ctx kontext.Context, queryKey, query string, args ...interface{}) Row
}

// Result summarizes an executed SQL command.
type Result interface {
	// LastInsertId returns the integer generated by the database
	// in response to a command. Typically this will be from an
	// "auto increment" column when inserting a new row. Not all
	// databases support this feature, and the syntax of such
	// statements varies.
	LastInsertId() (int64, exception.Exception)

	// RowsAffected returns the number of rows affected by an
	// update, insert, or delete. Not every database or database
	// driver may support this.
	RowsAffected() (int64, exception.Exception)
}

// Row single result of database query
type Row interface {
	Scan(dest ...interface{}) exception.Exception
}

// Rows multiple result of database query
type Rows interface {
	Close() exception.Exception
	Columns() ([]string, exception.Exception)
	Err() exception.Exception
	Next() bool
	NextResultSet() bool
	Scan(dest ...interface{}) exception.Exception
}
