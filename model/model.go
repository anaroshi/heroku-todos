package model

import (
	"time"
)

type Todo struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Completed bool `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type DBHandler interface {
	GetTodos(sessionId string) []*Todo
	AddTodo(sessionId, name string) *Todo
	RemoveTodo(id int) bool
	CompleteTodo(id int, complete bool) bool
	Close()
}

// 처음에 한번만 읽어들인다. 초기화시 사용
func NewDBHandler(dbConn string) DBHandler {
	return newPqHandler(dbConn)
}