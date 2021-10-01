package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDbPool() *sql.DB {
	dbName := os.Getenv("SOCIAL_DATABASE")
	db, err := sql.Open("mysql", dbName)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	connLifeTimeInMinutes, err := strconv.Atoi(os.Getenv("DB_CONN_LIFE_TIME"))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DB_CONN_LIFE_TIME isn't specified")
		os.Exit(1)
	}

	maxOpenConn, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTION"))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DB_MAX_OPEN_CONNECTION isn't specified")
		os.Exit(1)
	}

	maxIdleConn, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTION"))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DB_MAX_IDLE_CONNECTION isn't specified")
		os.Exit(1)
	}

	db.SetConnMaxLifetime(time.Minute * time.Duration(connLifeTimeInMinutes))
	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)

	//err = db.Ping()
	//if err != nil {
	//	_, _ = fmt.Fprintf(os.Stderr, "Couldn't ping database: %s", err)
	//	os.Exit(1)
	//}

	return db
}
