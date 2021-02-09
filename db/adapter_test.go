package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/codefluence-x/monorepo/db"
	"github.com/codefluence-x/monorepo/exception"
	"github.com/codefluence-x/monorepo/kontext"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAdapter(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ktx := kontext.Fabricate()

	t.Run("QueryRowContext", func(t *testing.T) {
		t.Run("When querying done it will return db.Row and scan the value", func(t *testing.T) {

			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectQuery(`select id from test_table where id = \?`).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

			sql := db.Adapt(sqldb)
			row := sql.QueryRowContext(ktx, "test-query-1", "select id from test_table where id = ?", 1)

			var x int

			assert.Nil(t, row.Scan(&x))
			assert.Equal(t, 1, x)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})

		t.Run("When querying done it will return db.Row and if the scan failed then it will return exception", func(t *testing.T) {

			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectQuery(`select id from test_table where id = \?`).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).RowError(0, errors.New("unexpected error")))

			sql := db.Adapt(sqldb)
			row := sql.QueryRowContext(ktx, "test-query-1", "select id from test_table where id = ?", 1)

			var x int
			exc := row.Scan(&x)
			assert.NotNil(t, exc)
			assert.Nil(t, mockDB.ExpectationsWereMet())
			assert.Equal(t, exception.Unexpected, exc.Type())
		})

		t.Run("When querying done it will return db.Row and if the scan failed because rows is not found then it will return not found exception", func(t *testing.T) {

			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectQuery(`select id from test_table where id = \?`).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id"}).RowError(0, sql.ErrNoRows))

			sql := db.Adapt(sqldb)
			row := sql.QueryRowContext(ktx, "test-query-1", "select id from test_table where id = ?", 1)

			var x int

			exc := row.Scan(&x)
			assert.NotNil(t, exc)
			assert.Nil(t, mockDB.ExpectationsWereMet())
			assert.Equal(t, exception.NotFound, exc.Type())
		})
	})

	t.Run("QueryContext", func(t *testing.T) {
		t.Run("When querying done it will return db.Rows", func(t *testing.T) {
			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectQuery(`select id from test_table`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(1))

			sql := db.Adapt(sqldb)
			rows, err := sql.QueryContext(ktx, "test-query-2", "select id from test_table")

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

			assert.Nil(t, err)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})

		t.Run("When querying got unexpected error then it will return exception", func(t *testing.T) {
			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectQuery(`select id from test_table`).WillReturnError(errors.New("unexpected error"))

			sql := db.Adapt(sqldb)
			_, exc := sql.QueryContext(ktx, "test-query-2", "select id from test_table")

			assert.Equal(t, exception.Unexpected, exc.Type())
			assert.NotNil(t, exc)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})

		t.Run("When querying got no rows then it will return not found exception", func(t *testing.T) {
			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectQuery(`select id from test_table`).WillReturnError(sql.ErrNoRows)

			sql := db.Adapt(sqldb)
			_, exc := sql.QueryContext(ktx, "test-query-2", "select id from test_table")

			assert.Equal(t, exception.NotFound, exc.Type())
			assert.NotNil(t, exc)
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})

		t.Run("When querying done it will return db.Rows, but if there is failed in row", func(t *testing.T) {
			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectQuery(`select id from test_table`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(1).RowError(1, errors.New("unexpected")))

			sql := db.Adapt(sqldb)
			rows, err := sql.QueryContext(ktx, "test-query-2", "select id from test_table")
			assert.Nil(t, err)

			i := 0
			for rows.Next() {
				i++
			}

			assert.Equal(t, 1, i)
			assert.NotNil(t, rows.Err())
			assert.Nil(t, mockDB.ExpectationsWereMet())
		})

		t.Run("When querying done but there is an error when scanning rows it will return error", func(t *testing.T) {

			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectQuery(`select id from test_table`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(1))

			sql := db.Adapt(sqldb)
			rows, err := sql.QueryContext(ktx, "test-query-2", "select id from test_table")

			var x string
			assert.NotNil(t, rows.Scan(&x))

			assert.Nil(t, mockDB.ExpectationsWereMet())
		})

		t.Run("When querying done but there is an error when closing rows it will return error", func(t *testing.T) {

			sqldb, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer sqldb.Close()

			mockDB.ExpectQuery(`select id from test_table`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).CloseError(errors.New("unexpected")))

			sql := db.Adapt(sqldb)
			rows, err := sql.QueryContext(ktx, "test-query-2", "select id from test_table")

			var x int
			assert.NotNil(t, rows.Scan(&x))
			assert.NotNil(t, rows.Close())

			assert.Nil(t, mockDB.ExpectationsWereMet())
		})

	})
}
