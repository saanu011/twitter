package api

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

var database *sqlx.DB

// warn user if the database is not set up after 10 seconds
func init() {
	go func() {
		<-time.After(10 * time.Second)
		if database == nil {
			logrus.Warn("DB variable was not properly set and is nil")
		}
	}()
}

// SetDBs sets the global db variable, making
//  a pool of database connections
// globally available to the package.  It is
// best practice initiating db connections from
// a main package and distribute them to libraries
// as needed, so that is what this function does.
func SetDBs(dbWriterParam *sqlx.DB) error {
	if err := dbWriterParam.Ping(); err != nil {
		fmt.Println("error: ", err)
		return err
	}

	database = dbWriterParam

	return nil
}

