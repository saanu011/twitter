package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
	"twitter/api"
)

const (
	appName = "twitter"
	// dbOpenConnsMax controls the maximum number of allowed, simultaneous open DB connections.
	dbOpenConnsMax = 5
)

var (
	// DBServer is the MySQL server queries should be sent to
	DBServer = os.Getenv("DB_HOST")
	DBUser   = os.Getenv("DB_PORT")
	DBPass   = os.Getenv("DB_PASSWORD")
	DBPort   = os.Getenv("DB_PORT")
	DBName   = os.Getenv("DB_NAME")

	secret = os.Getenv("JWT_SECRET")

	dbx *sqlx.DB
)

func main() {

	err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("ENVIRONMENT_NAME")))
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	fmt.Println(" starting things ")

	// setup write database
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&clientFoundRows=true",
		DBUser,
		DBPass,
		DBServer,
		DBPort,
		DBName,
	)
	writerLogFields := logrus.Fields{
		"dbServer": DBServer,
		"dbUser":   DBUser,
		"dbPort":   DBPort,
	}

	dbx, err = sqlx.Open("mysql", url)
	if err != nil {
		logrus.WithError(err).WithFields(writerLogFields).Error("failed to connect to mysql server")
		return
	}
	defer dbx.Close()

	// ensure write database is online
	if err := dbx.Ping(); err != nil {
		logrus.WithError(err).WithFields(writerLogFields).Error("failed to ping mysql server")
		return
	}
	logrus.WithFields(writerLogFields).Info("successfully connected to mysql server")

	dbx = dbx.Unsafe()
	dbx.SetMaxIdleConns(0)
	dbx.SetMaxOpenConns(dbOpenConnsMax)
	dbx.SetConnMaxLifetime(time.Minute * time.Duration(5))

	// set package level databases
	if err = api.SetDBs(dbx); err != nil {
		logrus.WithError(err).Error("failed to set mysql DB")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	// start the flourish API server
	router := NewRouter()

	port := "8088"
	twitterServer := &http.Server{
		Addr:         ":" + port,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,

		// Pass our instance of gorilla/mux in
		Handler: &MyServer{router},
	}
	go func() {
		logrus.WithError(twitterServer.ListenAndServe()).
			WithField("port", port).
			Error("failed to start flourish API server")
		return
	}()

	// graceful shutdown when termination signals received
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)

	sig := <-sigquit
	logrus.WithField("signal", sig).Info("caught interrupt signal, gracefully shutting down server")
	cancel()

	// shutdown the API server, waiting for any outstanding requests to complete
	twitterServer.Shutdown(ctx)
	logrus.Info("graceful server shutdown complete, exiting")
	return
}
