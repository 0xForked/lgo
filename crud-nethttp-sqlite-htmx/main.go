package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/0xForked/crud-nethttp-sqlite-htmx/server"
	_ "github.com/glebarez/go-sqlite"
)

// server port
var port *string

// database driver and source
const driver, source = "sqlite", "./person.sqlite3"

func init() {
	// parse the flags
	port = flag.String("port",
		":8000", "server port")
	flag.Parse()
}

func newSQLiteConn() *sql.DB {
	db, err := sql.Open(
		driver, source)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	// create a new server mux
	mux := http.NewServeMux()
	// create a new database connection
	db := newSQLiteConn()
	// handle the person api routes
	server.NewHandler(mux, db)
	// start the server
	if err := http.ListenAndServe(
		*port, mux,
	); err != nil {
		log.Fatal("server failed to start: ", err)
	}
}
