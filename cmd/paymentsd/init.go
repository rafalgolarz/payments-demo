package main

import (
	"os"

	log "github.com/rafalgolarz/payments-demo/pkg/log/logrus"
	storage "github.com/rafalgolarz/payments-demo/pkg/storage"
)

var (
	apiPort string
	// compile passing -ldflags "-X main.Build <build sha1>"
	build string

	logConfig = &log.Config{
		LogLevel: "error",
	}

	logHandler *log.LogHandler
)

var paymentsStorage = &storage.DBHandler{
	Cfg: &storage.Config{
		User: "rafal",
		//Password:         os.Getenv("MYSQL_ROOT_PASSWORD"),
		Password:         "password",
		Port:             "3306",
		Db:               "payments_demo",
		Host:             "localhost",
		AdditionalParams: "?charset=utf8&parseTime=true&loc=UTC",
	},
}

// mysqlConnector
var dbConn storage.Connector

func init() {
	apiPort = os.Getenv("DEFAULT_API_PORT")
	if apiPort == "" {
		apiPort = ":8080"
	}
}

func setup(sc *storage.DBHandler) (storage.Connector, error) {

	var err error

	dbConn := storage.InitMysql(paymentsStorage)
	err = dbConn.Connect()
	conn := sc.DBConn()

	if err != nil {
		logHandler.Error("problem connecting to database", log.Fields{"dbname": sc.Cfg.Db, "func": "setup"})
		return nil, err
	}

	switch dbLogMode := os.Getenv("DB_LOG_MODE"); dbLogMode {
	case "true":
		conn.SetLogMode(true)
	default:
		conn.SetLogMode(false)
	}
	conn.SetLogMode(true)
	conn.SetConnMaxLifetime(300)
	conn.SetMaxIdleConns(0)
	conn.SetMaxOpenConns(1)

	return dbConn, nil
}
