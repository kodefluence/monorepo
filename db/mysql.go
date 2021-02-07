package db

import (
	"database/sql"
	"fmt"

	"github.com/codefluence-x/monorepo/exception"
)

// FabricateMySQL will fabricate mysql connection and wrap it into SQL interfaces
func FabricateMySQL(config Config) (DB, exception.Exception) {
	_, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&interpolateParams=true", config.Username, config.Password, config.Host, config.Port, config.Name))
	if err != nil {
		return nil, exception.Throw(err)
	}

	return nil, nil
}
