package db_test

import (
	"testing"
	"time"

	"github.com/kodefluence/monorepo/db"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationDB(t *testing.T) {
	t.Run("MySQL", func(t *testing.T) {
		t.Run("Complete open config to mysql", func(t *testing.T) {
			config := db.Config{
				Username: "root",
				Password: "rootpw",
				Host:     "localhost",
				Port:     "3306",
				Name:     "test_database",
			}

			sqldb, err := db.FabricateMySQL("main_db", config, db.WithConnMaxLifetime(time.Second), db.WithMaxIdleConn(100), db.WithMaxOpenConn(100))
			assert.NotNil(t, sqldb)
			assert.Nil(t, err)

			sqldb, err = db.FabricateMySQL("main_db", config, db.WithConnMaxLifetime(time.Second), db.WithMaxIdleConn(100), db.WithMaxOpenConn(100))
			assert.NotNil(t, sqldb)
			assert.Nil(t, err)

			sqldbvalue, err := db.GetInstance("main_db")
			assert.NotNil(t, sqldbvalue)
			assert.Nil(t, err)

			sqldbvalue, err = db.GetInstance("no_db")
			assert.Nil(t, sqldbvalue)
			assert.NotNil(t, err)

			sqldb, err = db.FabricateMySQL("secondary_db", config, db.WithConnMaxLifetime(time.Second), db.WithMaxIdleConn(100), db.WithMaxOpenConn(100))
			assert.NotNil(t, sqldb)
			assert.Nil(t, err)

			assert.Equal(t, 0, len(db.CloseAll()))
		})
	})
}
