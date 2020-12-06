package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	MysqlDB *sql.DB
)

const (
	DatabaseError = "Internal server error, database disconnected. Contact Administrator"
)

func connect() *sql.DB {
	user := "user"
	password := "user@123"
	host := "157.245.55.51"
	database := "pymnt_db"

	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, database)

	db, err := sql.Open("mysql", connection)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

//check database current connection
func pingDb(db *sql.DB) error {
	err := db.Ping()
	return err
}
