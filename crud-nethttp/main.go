package main

import (
	"log"
	"net/http"

	"learn.go/crud-net-http/handler"
	"learn.go/crud-net-http/model"
)

func main() {
	h := handler.HTTPHandler{Persons: []*model.Person{}}
	http.HandleFunc("/persons", h.HandleListAndCreate)
	http.HandleFunc("/persons/", h.HandleDetailAndModify)
	log.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
