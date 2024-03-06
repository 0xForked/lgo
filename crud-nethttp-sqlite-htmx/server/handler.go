package server

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/0xForked/crud-nethttp-sqlite-htmx/model"
	"github.com/0xForked/crud-nethttp-sqlite-htmx/resource"
)

type handler struct {
	db *sql.DB
}

func (h *handler) fetch(w http.ResponseWriter, r *http.Request) {
	person := model.Person{}
	data, err := person.List(h.db)
	component := resource.Home(data, err)
	_ = component.Render(r.Context(), w)
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	// check if the request has a form
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// check if the form has a name field
	if !r.Form.Has("name") {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	// create a new person
	newPerson := model.Person{Name: r.Form.Get("name")}
	if _, err := newPerson.Insert(h.db); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/"):]
	id, err := strconv.Atoi(param)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newPerson := model.Person{ID: id}
	if err := newPerson.Delete(h.db); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func NewHandler(mux *http.ServeMux, db *sql.DB) {
	h := &handler{db: db}
	mux.HandleFunc("GET /", h.fetch)
	mux.HandleFunc("POST /", h.create)
	mux.HandleFunc("DELETE /{id}", h.delete)
}
