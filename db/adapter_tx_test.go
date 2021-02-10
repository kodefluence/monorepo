package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/codefluence-x/monorepo/db"
	"github.com/codefluence-x/monorepo/exception"
	"github.com/codefluence-x/monorepo/kontext"
	"github.com/stretchr/testify/assert"
)

func TestTx(t *testing.T) {
	ktx := kontext.Fabricate()

	t.Run("QueryRowContext", func(t *testing.T) {
		t.Run("When querying done it will return db.Row and scan the value", func(t *testing.T) {
			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectBegin()
			mockDB.ExpectQuery(`select id from test_table where id = \?`).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mockDB.ExpectCommit()

			sql := db.Adapt(sqldb)
			err = sql.Transaction(ktx, "transaction-test", func(tx db.TX) exception.Exception {
				row := tx.QueryRowContext(ktx, "test-query-1", "select id from test_table where id = ?", 1)
				var x int
				return row.Scan(&x)
			})

			assert.Nil(t, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})

		t.Run("When begin tx failed then it will return error", func(t *testing.T) {
			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectBegin().WillReturnError(errors.New("unexpected error"))

			sql := db.Adapt(sqldb)
			err = sql.Transaction(ktx, "transaction-test", func(tx db.TX) exception.Exception {
				row := tx.QueryRowContext(ktx, "test-query-1", "select id from test_table where id = ?", 1)
				var x int
				return row.Scan(&x)
			})

			assert.NotNil(t, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})

		t.Run("When querying error then it will return exception", func(t *testing.T) {
			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectBegin()
			mockDB.ExpectQuery(`select id from test_table where id = \?`).WithArgs(1).WillReturnError(errors.New("unexpected"))
			mockDB.ExpectRollback()

			sql := db.Adapt(sqldb)
			err = sql.Transaction(ktx, "transaction-test", func(tx db.TX) exception.Exception {
				row := tx.QueryRowContext(ktx, "test-query-1", "select id from test_table where id = ?", 1)
				var x int
				return row.Scan(&x)
			})

			assert.NotNil(t, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})

	t.Run("QueryContext", func(t *testing.T) {
		t.Run("When querying done it will return db.Rows", func(t *testing.T) {
			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectBegin()
			mockDB.ExpectQuery(`select id from test_table`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(1))
			mockDB.ExpectCommit()

			sql := db.Adapt(sqldb)
			err = sql.Transaction(ktx, "transaction-test", func(tx db.TX) exception.Exception {
				rows, err := tx.QueryContext(ktx, "test-query-2", "select id from test_table")

				column, err := rows.Columns()
				assert.Nil(t, err)
				assert.Equal(t, []string{"id"}, column)

				for rows.Next() {
					var x int
					assert.Nil(t, rows.Scan(&x))
				}

				assert.Nil(t, rows.Err())
				assert.Nil(t, rows.Close())

				// Will return error
				_, exc := rows.Columns()
				assert.NotNil(t, exc)
				assert.Equal(t, exception.Unexpected, exc.Type())
				return err
			})

			assert.Nil(t, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})

		t.Run("When querying done it will return db.Rows but if the commit failed then it will return error", func(t *testing.T) {
			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectBegin()
			mockDB.ExpectQuery(`select id from test_table`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(1))
			mockDB.ExpectCommit().WillReturnError(errors.New("unexpected error"))

			sql := db.Adapt(sqldb)
			err = sql.Transaction(ktx, "transaction-test", func(tx db.TX) exception.Exception {
				rows, err := tx.QueryContext(ktx, "test-query-2", "select id from test_table")

				column, err := rows.Columns()
				assert.Nil(t, err)
				assert.Equal(t, []string{"id"}, column)

				for rows.Next() {
					var x int
					assert.Nil(t, rows.Scan(&x))
				}

				assert.Nil(t, rows.Err())
				assert.Nil(t, rows.Close())

				// Will return error
				_, exc := rows.Columns()
				assert.NotNil(t, exc)
				assert.Equal(t, exception.Unexpected, exc.Type())
				return err
			})

			assert.NotNil(t, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})

		t.Run("When querying error because the data is not found it will return error", func(t *testing.T) {
			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectBegin()
			mockDB.ExpectQuery(`select id from test_table`).WillReturnError(sql.ErrNoRows)
			mockDB.ExpectRollback()

			sql := db.Adapt(sqldb)
			err = sql.Transaction(ktx, "transaction-test", func(tx db.TX) exception.Exception {
				_, err := tx.QueryContext(ktx, "test-query-2", "select id from test_table")

				return err
			})

			assert.NotNil(t, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})

		t.Run("When querying error because unexpected error it will return error", func(t *testing.T) {
			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectBegin()
			mockDB.ExpectQuery(`select id from test_table`).WillReturnError(errors.New("unexpected error"))
			mockDB.ExpectRollback()

			sql := db.Adapt(sqldb)
			err = sql.Transaction(ktx, "transaction-test", func(tx db.TX) exception.Exception {
				_, err := tx.QueryContext(ktx, "test-query-2", "select id from test_table")

				return err
			})

			assert.NotNil(t, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})

	t.Run("ExecContext", func(t *testing.T) {
		t.Run("When query execution is done it will return results", func(t *testing.T) {
			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectBegin()
			mockDB.ExpectExec(`insert into users \(1\) values\(id\)`).WillReturnResult(sqlmock.NewResult(1, 1))
			mockDB.ExpectCommit()

			sql := db.Adapt(sqldb)
			err = sql.Transaction(ktx, "transaction-test", func(tx db.TX) exception.Exception {
				result, err := tx.ExecContext(ktx, "test-query-1", "insert into users (1) values(id)")
				lastID, err := result.LastInsertId()
				assert.Equal(t, int64(1), lastID)
				assert.Nil(t, err)

				rowsAffected, err := result.RowsAffected()
				assert.Equal(t, int64(1), rowsAffected)
				assert.Nil(t, err)

				return err
			})

			assert.Nil(t, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})

		t.Run("When query execution error then it will return error", func(t *testing.T) {
			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectBegin()
			mockDB.ExpectExec(`insert into users \(1\) values\(id\)`).WillReturnError(errors.New("unexpected error"))
			mockDB.ExpectRollback()

			sql := db.Adapt(sqldb)
			err = sql.Transaction(ktx, "transaction-test", func(tx db.TX) exception.Exception {
				_, err := tx.ExecContext(ktx, "test-query-1", "insert into users (1) values(id)")
				return err
			})

			assert.NotNil(t, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})
	})
}
