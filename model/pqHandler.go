package model

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type pqHandler struct {
	db *sql.DB
}

func chkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (s *pqHandler) GetTodos(sessionId string) []*Todo {
	todos := []*Todo{}
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos WHERE sessionId=?", sessionId)
	chkErr(err)
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt)
		todos = append(todos, &todo)
	}
	return todos
}

func (s *pqHandler) AddTodo(name, sessionId string) *Todo {
	stmt, err := s.db.Prepare("INSERT INTO todos (sessionId, name, completed, createdAt) VALUES (?, ?, ?, datetime('now'))")
	chkErr(err)
	
	rst, err := stmt.Exec(sessionId, name, false)
	chkErr(err)
	id, _ := rst.LastInsertId()
	var todo Todo
	todo.ID = int(id)
	todo.Name = name
	todo.Completed = false
	todo.CreatedAt = time.Now()
	return &todo
}

func (s *pqHandler) RemoveTodo(id int) bool {
	stmt, err := s.db.Prepare("DELETE FROM todos WHERE id=?")
	chkErr(err)
	rst, err := stmt.Exec(id)
	chkErr(err)
	cnt, err := rst.RowsAffected()
	chkErr(err)	
	return cnt>0
}

func (s *pqHandler) CompleteTodo(id int, complete bool) bool {
	stmt, err := s.db.Prepare("UPDATE todos SET completed=? WHERE id=?")
	chkErr(err)
	rst, err := stmt.Exec(complete, id)
	chkErr(err)
	cnt, err := rst.RowsAffected()
	chkErr(err)
	return cnt>0
}

func (s *pqHandler) Close() {
	s.db.Close()
}

func newPqHandler(dbConn string) DBHandler {
	
	database, err := sql.Open("postgres", dbConn)
	chkErr(err)
	
	statement, err := database.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sessionId STRING,
			name TEXT,
			completed BOOLEAN,
			createdAt DATETIME
		);
		CREATE INDEX IF NOT EXISTS sessionIdIndexOnTodos ON todos (
			sessionId ASC
		);`)
	chkErr(err)	
	_, err = statement.Exec()
	chkErr(err)
	return &sqliteHandler{db: database}
}