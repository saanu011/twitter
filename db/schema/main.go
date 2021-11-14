package main

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"twitter/api"
	"twitter/utils"
)

var (
	DBServer = os.Getenv("DB_HOST")
	DBUser   = os.Getenv("DB_PORT")
	DBPass   = os.Getenv("DB_PASSWORD")
	DBPort   = os.Getenv("DB_PORT")
	DBName   = os.Getenv("DB_NAME")
)

func main() {

	// ------- SETUP --------

	// set format of logs to json
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// set output to log more data about callers
	logrus.SetReportCaller(true)

	// ------- Flags --------

	// parse cli flags
	flag.Parse()

	// ------- Connections --------

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// set up the mysql connection
	db, err := setupDB()
	if err != nil {
		log.Fatal(fmt.Errorf("could not connect to database"))
	}

	// ------- Create Schema --------

	err = runUpsertQueries(db, getCreateQueries())
	if err != nil {
		log.Fatal(fmt.Errorf("error running queries: %v", err))
	}

	log.Println("Created successfully")

}

func runUpsertQueries(db *sqlx.DB, queries []utils.UpsertQuery) (err error) {
	errTag := "runUpsertQueries"

	// guard against nothing to run and empty transactions
	if len(queries) == 0 {
		return nil
	}

	// open a new flourishDB transaction
	txo, err := db.Beginx()
	if err != nil {
		return fmt.Errorf(err.Error(), "failed to begin transaction runUpsetQueries")
	}

	// loop through each
	for _, uq := range queries {

		// exec & terminate if error
		_, err := txo.Exec(uq.Query, uq.Parameters...)
		if err != nil {
			logrus.WithField("query", uq).Error("query failed")
			if terminateErr := txo.Rollback(); terminateErr != nil {
				logrus.WithField("query", uq).Error("could not rollback txn")
				return fmt.Errorf(terminateErr.Error(), "could not rollback txn")
			}

			return fmt.Errorf(err.Error(), "query execution failed")
		}

		logrus.Info("----- successfully created a table -----")
	}

	// commit the transaction
	if txo.Commit() != nil {
		return fmt.Errorf(fmt.Sprintf("failed to commit transaction"), errTag)
	}

	return nil
}

// getCreateQueries gets all the queries to create Flourish schema
func getCreateQueries() []utils.UpsertQuery {
	// Build query to run in MySQL
	createQueries := []utils.UpsertQuery{

		// user
		{
			Query: fmt.Sprintf(`
				CREATE TABLE user (
					user_id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
					email VARCHAR(55) NOT NULL,
					name VARCHAR(55) NOT NULL,
					address VARCHAR(55) NOT NULL,
					last_updated_by VARCHAR(45) NOT NULL,
					last_updated_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
					created_by VARCHAR(45) NOT NULL,
					created_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);`,
			),
			Parameters: []interface{}{},
		},
		// auth_client
		{
			Query: fmt.Sprintf(`
				CREATE TABLE auth_client (
					id                     INT       NOT NULL PRIMARY KEY AUTO_INCREMENT,
					client_ref             TEXT      NOT NULL,
					client_secret          TEXT      NOT NULL,
					created_timestamp      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					last_updated_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
				);`,
			),
			Parameters: []interface{}{},
		},
		// auth_access_tokens
		{
			Query: fmt.Sprintf(`
				CREATE TABLE auth_access_token (
					id                              INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
					auth_client_id                  INTEGER NOT NULL,
					access_token                    TEXT NOT NULL,
					expires_in                      VARCHAR(32) NOT NULL,
					expires_on                      TIMESTAMP(0) NOT NULL,
					created_timestamp         		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					last_updated_timestamp          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
					FOREIGN KEY (auth_client_id)    REFERENCES auth_client (id)
				);`,
			),
			Parameters: []interface{}{},
		},
		{
			Query: fmt.Sprintf(`
				CREATE TABLE relationships (
					id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
					followed_id INT NOT NULL COMMENT 'user_id of the user which is being followed',
					follower_id INT NOT NULL COMMENT 'user_id of the user which is following',
					last_updated_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
					created_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					CONSTRAINT fk_followed_id FOREIGN KEY (followed_id) 
					REFERENCES user(user_id),
					CONSTRAINT fk_follower_id FOREIGN KEY (follower_id) 
					REFERENCES user(user_id)
				);`,
			),
			Parameters: []interface{}{},
		},
		{
			Query: fmt.Sprintf(`
				CREATE TABLE tweet (
					id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
					content TEXT,
					user_id INT NOT NULL COMMENT '',
					last_updated_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
					created_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					CONSTRAINT fk_user_id FOREIGN KEY (user_id) 
					REFERENCES user(user_id)
				);`,
			),
			Parameters: []interface{}{},
		},
	}
	return createQueries
}

func setupDB() (db *sqlx.DB, err error) {
	// var New func(proto, laddr, raddr, user, passwd string, db ...string) Conn
	// New("tcp", "", fmt.Sprintf("%s:%s", clientdb_server, clientdb_port), clientdb_user, clientdb_pass
	// DSN format: id:password@tcp(your-amazonaws-uri.com:3306)/dbname
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&clientFoundRows=true",
		DBUser,
		DBPass,
		DBServer,
		DBPort,
		DBName,
	)
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Println("error sqlx.Open: ", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		fmt.Println("error db.Ping: ", err)
		return nil, err
	}

	if err = api.SetDBs(db); err != nil {
		fmt.Println("error api.SetDBs: ", err)
		return nil, err
	}

	return db, nil
}
