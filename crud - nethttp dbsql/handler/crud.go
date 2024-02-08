package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"learn.go/crud-net-http-db-sql/model"
)

type HTTPHandler struct {
	DB *sql.DB
}

func (h *HTTPHandler) HandleListAndCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.fetch(w, r)
	case http.MethodPost:
		h.add(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HTTPHandler) HandleDetailAndModify(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.show(w, r)
	case http.MethodPut:
		h.edit(w, r)
	case http.MethodDelete:
		h.destroy(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HTTPHandler) fetch(w http.ResponseWriter, _ *http.Request) {
	person := model.Person{}
	data, err := person.List(h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func (h *HTTPHandler) add(w http.ResponseWriter, r *http.Request) {
	var newPerson model.Person
	if err := json.NewDecoder(r.Body).Decode(&newPerson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := newPerson.Insert(h.DB); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.fetch(w, r)
}

func (h *HTTPHandler) show(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/persons/"):]
	id, err := strconv.Atoi(param)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	person := model.Person{ID: id}
	data, err := person.Detail(h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func (h *HTTPHandler) edit(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/persons/"):]
	id, err := strconv.Atoi(param)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var person model.Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	person.ID = id
	if err := person.Update(h.DB); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.fetch(w, r)
}

func (h *HTTPHandler) destroy(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/persons/"):]
	id, err := strconv.Atoi(param)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	person := model.Person{ID: id}
	if err := person.Delete(h.DB); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.fetch(w, r)
}
