package main

import (
	"math/rand"
	"time"
)

type Todo struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Created_at  time.Time `json:"created_at"`
}

type UpdateTodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type CreateTodoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func NewTodo(name, description string, completed bool) (*Todo, error) {
	return &Todo{
		ID:          int64(rand.Intn(10000)),
		Name:        name,
		Description: description,
		Completed:   completed,
		Created_at:  time.Now().UTC(),
	}, nil
}
