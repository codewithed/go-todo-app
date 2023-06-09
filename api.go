package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ApiServer struct {
	ListenAddr string
	Store      Storage
}

func NewApiServer(addr string, store Storage) *ApiServer {
	return &ApiServer{
		ListenAddr: addr,
		Store:      store,
	}
}

func (s *ApiServer) Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "welcome"}`))
	})
	r.HandleFunc("/todos", makeHTTPHandleFunc(s.handleGetTodos))
	r.HandleFunc("/todos/{id}", makeHTTPHandleFunc(s.handleGetTodoByID))
	http.ListenAndServe(s.ListenAddr, r)
}

func (s *ApiServer) handleGetTodos(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleGetTodoByID(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
