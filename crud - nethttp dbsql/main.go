package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/glebarez/go-sqlite"
	"learn.go/crud-net-http-db-sql/handler"
)

func createDBConn() *sql.DB {
	driver, source := "sqlite", "./person.sqlite3"
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	db := createDBConn()
	h := handler.HTTPHandler{DB: db}
	http.HandleFunc("/persons", h.HandleListAndCreate)
	http.HandleFunc("/persons/", h.HandleDetailAndModify)
	log.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
