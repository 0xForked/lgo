package todo

import (
	"encoding/json"
	"net/http"
)

type handler struct {
	todo TODO
}

func (h handler) fetch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todos := h.todo.Fetch()
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(todos)
}

func (h handler) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newTodo TODO
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo := h.todo.Create(newTodo.Body)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(todo)
}

func NewHandler(mux *http.ServeMux) {
	h := &handler{todo: TODO{}}
	mux.HandleFunc("GET /todos", h.fetch)
	mux.HandleFunc("POST /todos", h.create)
}
