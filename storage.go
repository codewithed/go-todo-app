package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateTodo(*Todo) error
	DeleteTodo(int) error
	UpdateTodo(int, *UpdateTodoRequest) error
	GetTodos() ([]*Todo, error)
	GetTodoByID(int) (*Todo, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gotodoapp sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.createTodosTable()
}

func (s *PostgresStore) createTodosTable() error {
	query := `CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY NOT NULL,
		name varchar(1000) NOT NULL,
		description varchar(10000) NOT NULL,
		completed bool NOT NULL,
		created_at timestamp NOT NULL
	)`

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) CreateTodo(todo *Todo) error {
	query := `INSERT INTO todos (name, description, completed, created_at) VALUES(
		$1, $2, $3, $4
	)`
	_, err := s.db.Query(
		query,
		todo.Name,
		todo.Description,
		todo.Completed,
		todo.Created_at,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetTodoByID(id int) (*Todo, error) {
	rows, err := s.db.Query(`SELECT * FROM todos WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoTodo(rows)
	}

	return nil, fmt.Errorf("todo with id %d not found", id)
}

func (s *PostgresStore) DeleteTodo(id int) error {
	_, err := s.db.Query(`DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateTodo(id int, req *UpdateTodoRequest) error {

	if req.Name != "" {
		_, err := s.db.Exec(`UPDATE todos 
		SET name = $1 
		WHERE id = $2`, req.Name, id)
		if err != nil {
			return err
		}
	}

	if req.Description != "" {
		_, err := s.db.Exec(`UPDATE todos 
		SET description = $1
		WHERE id = $2`, req.Description, id)
		if err != nil {
			return err
		}
	}

	if req.Completed == true || req.Completed == false {
		_, err := s.db.Exec(`UPDATE todos 
		SET completed = $1
		WHERE id = $2`, req.Completed, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresStore) GetTodos() ([]*Todo, error) {

	rows, err := s.db.Query(
		`SELECT * FROM todos`,
	)
	if err != nil {
		return nil, err
	}

	todos := []*Todo{}
	for rows.Next() {
		todo, err := scanIntoTodo(rows)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func scanIntoTodo(rows *sql.Rows) (*Todo, error) {
	todo := new(Todo)
	err := rows.Scan(
		&todo.ID,
		&todo.Name,
		&todo.Description,
		&todo.Completed,
		&todo.Created_at,
	)

	return todo, err
}
