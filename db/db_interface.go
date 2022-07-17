package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/kodefluence/monorepo/exception"

	// Import mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var instanceList = &sync.Map{}

// FabricateMySQL will fabricate mysql connection and wrap it into SQL interfaces
func FabricateMySQL(instanceName string, config Config, opts ...Option) (DB, exception.Exception) {
	if val, ok := instanceList.Load(fmt.Sprintf("mysql-%s", instanceName)); ok {
		return Adapt(val.(*sql.DB)), nil
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&interpolateParams=true", config.Username, config.Password, config.Host, config.Port, config.Name))
	if err != nil {
		return nil, exception.Throw(err)
	}

	// Default value
	config.maxIdleConn = 2
	config.maxOpenConn = 0

	for _, opt := range opts {
		opt(&config)
	}

	db.SetConnMaxLifetime(config.connMaxLifetime)
	db.SetMaxIdleConns(config.maxIdleConn)
	db.SetMaxOpenConns(config.maxOpenConn)

	instanceList.Store(fmt.Sprintf("mysql-%s", instanceName), db)

	return Adapt(db), nil
}

// GetInstance that already fabricated before as an sql.DB
func GetInstance(instanceName string) (*sql.DB, exception.Exception) {
	if val, ok := instanceList.Load(fmt.Sprintf("mysql-%s", instanceName)); ok {
		return val.(*sql.DB), nil
	}

	return nil, exception.Throw(errors.New("unexpected error"), exception.WithType(exception.NotFound))
}

// CloseAll initiated mysql connection
func CloseAll() []exception.Exception {
	var excs []exception.Exception

	instanceList.Range(func(key, value interface{}) bool {
		if err := value.(*sql.DB).Close(); err != nil {
			excs = append(excs, exception.Throw(err, exception.WithTitle("error closing mysql connection"), exception.WithDetail(fmt.Sprintf("instance name: %s", key.(string)))))
		}

		return true
	})

	return excs
}
