package main

import (
	"log"
	"net/http"
)

func main() {
	MysqlDB = connect()
	defer MysqlDB.Close()
	router := server()
	log.Fatal(http.ListenAndServe(":8000", router))
}
