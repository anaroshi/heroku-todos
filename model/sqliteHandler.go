package model

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteHandler struct {
	db *sql.DB
}

func ChkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (s *sqliteHandler) GetTodos(sessionId string) []*Todo {
	todos := []*Todo{}
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos WHERE sessionId=?", sessionId)
	ChkErr(err)
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt)
		todos = append(todos, &todo)
	}
	return todos
}

func (s *sqliteHandler) AddTodo(name, sessionId string) *Todo {
	stmt, err := s.db.Prepare("INSERT INTO todos (sessionId, name, completed, createdAt) VALUES (?, ?, ?, datetime('now'))")
	ChkErr(err)
	
	rst, err := stmt.Exec(sessionId, name, false)
	ChkErr(err)
	id, _ := rst.LastInsertId()
	var todo Todo
	todo.ID = int(id)
	todo.Name = name
	todo.Completed = false
	todo.CreatedAt = time.Now()
	return &todo
}

func (s *sqliteHandler) RemoveTodo(id int) bool {
	stmt, err := s.db.Prepare("DELETE FROM todos WHERE id=?")
	ChkErr(err)
	rst, err := stmt.Exec(id)
	ChkErr(err)
	cnt, err := rst.RowsAffected()
	ChkErr(err)	
	return cnt>0
}

func (s *sqliteHandler) CompleteTodo(id int, complete bool) bool {
	stmt, err := s.db.Prepare("UPDATE todos SET completed=? WHERE id=?")
	ChkErr(err)
	rst, err := stmt.Exec(complete, id)
	ChkErr(err)
	cnt, err := rst.RowsAffected()
	ChkErr(err)
	return cnt>0
}

func (s *sqliteHandler) Close() {
	s.db.Close()
}

func newSqliteHandler(filepath string) DBHandler {
	
	database, err := sql.Open("sqlite3", filepath)
	ChkErr(err)
	
	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			sessionId STRING,
			name TEXT,
			completed BOOLEAN,
			createdAt DATETIME
		);
		CREATE INDEX IF NOT EXISTS sessionIdIndexOnTodos ON todos (
			sessionId ASC
		);`)
	statement.Exec()
	return &sqliteHandler{db: database}
}