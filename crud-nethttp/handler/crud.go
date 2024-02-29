package handler

import (
	"encoding/json"
	"net/http"

	"learn.go/crud-net-http/model"
	"learn.go/crud-net-http/util"
)

type HTTPHandler struct {
	Persons []*model.Person
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
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(h.Persons)
}

func (h *HTTPHandler) add(w http.ResponseWriter, r *http.Request) {
	var newPerson model.Person
	if err := json.NewDecoder(r.Body).Decode(&newPerson); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newPerson.ID = util.GenID(6)
	h.Persons = append(h.Persons, &newPerson)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(h.Persons)
}

func (h *HTTPHandler) show(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/persons/"):]
	var person *model.Person
	for _, p := range h.Persons {
		if param != p.ID {
			continue
		}
		person = p
		break
	}
	if person == nil {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(person)
}

func (h *HTTPHandler) edit(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/persons/"):]
	var person model.Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, p := range h.Persons {
		if param != p.ID {
			continue
		}
		p.Name = person.Name
		break
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(h.Persons)
}

func (h *HTTPHandler) destroy(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/persons/"):]
	var newPersons []*model.Person
	for _, p := range h.Persons {
		if param == p.ID {
			continue
		}
		newPersons = append(newPersons, p)
	}
	h.Persons = newPersons
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(h.Persons)
}
