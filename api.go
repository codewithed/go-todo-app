package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	r.HandleFunc("/todos", makeHTTPHandleFunc(s.handleTodos))
	r.HandleFunc("/todos/{id}", makeHTTPHandleFunc(s.handleTodosByID))
	http.ListenAndServe(s.ListenAddr, r)
}

func (s *ApiServer) handleTodos(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAllTodos(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateTodo(w, r)
	}
	return nil
}

func (s *ApiServer) handleTodosByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetTodoByID(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteTodoByID(w, r)
	}

	if r.Method == "PUT" || r.Method == "PATCH" {
		return s.handleUpdateTodoByID(w, r)
	}
	return nil
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func (s *ApiServer) handleGetAllTodos(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.Store.GetTodos()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *ApiServer) handleCreateTodo(w http.ResponseWriter, r *http.Request) error {
	createTodoRequest := new(CreateTodoRequest)
	// filling up the createTodoRequest with the JSON from the request body
	if err := json.NewDecoder(r.Body).Decode(createTodoRequest); err != nil {
		return err
	}

	todo, err := NewTodo(createTodoRequest.Name, createTodoRequest.Description, createTodoRequest.Completed)
	if err != nil {
		return err
	}

	if err := s.Store.CreateTodo(todo); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, todo)
}

func (s *ApiServer) handleGetTodoByID(w http.ResponseWriter, r *http.Request) error {
	// get id from the request body
	id, err := getID(r)
	if err != nil {
		return err
	}

	// perform the delete operation using the s.Store.DeleteTodoById operation
	todo, err := s.Store.GetTodoByID(id)
	if err != nil {
		return err
	}

	// return success
	return WriteJSON(w, http.StatusOK, todo)

}

func (s *ApiServer) handleDeleteTodoByID(w http.ResponseWriter, r *http.Request) error {
	// get id from the request body
	id, err := getID(r)
	if err != nil {
		return err
	}

	// perform the delete operation using the s.Store.DeleteTodoById operation
	if err := s.Store.DeleteTodo(id); err != nil {
		return err
	}

	// return success
	return WriteJSON(w, http.StatusOK, fmt.Sprintf("todo with id %d deleted", id))

}

func (s *ApiServer) handleUpdateTodoByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	// create a new updateTodorequest from the request body
	updateTodoRequest := new(UpdateTodoRequest)
	if err := json.NewDecoder(r.Body).Decode(updateTodoRequest); err != nil {
		return err
	}

	// perform the database operatyion using the s.Store.UpdateTodo
	if err := s.Store.UpdateTodo(id, updateTodoRequest); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, updateTodoRequest)
}

// HELPER FUNCTIONS
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

func getID(r *http.Request) (int, error) {
	id := chi.URLParam(r, "id")

	// validate the id
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return idInt, fmt.Errorf("Invalid id given: %v", id)
	}

	return idInt, nil
}
