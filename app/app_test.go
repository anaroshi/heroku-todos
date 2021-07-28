package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"sundorWeb/todos/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodos(t *testing.T) {
	getSessionID = func (r *http.Request ) string {
		return "testsessionId"
	}
	os.Remove("./test.db")
	assert := assert.New(t)
	ah := MakeHandler("./test.db")
	defer ah.Close()
	
	ts := httptest.NewServer(ah)
	defer ts.Close()

	// addTodoHandler 테스트go
	// form value로 받기 때문에 PostForm 씀
	resp, err := http.PostForm(ts.URL+"/todos", url.Values{"name":{"Test todo"}})   
	assert.NoError(err)
	assert.Equal(http.StatusCreated,resp.StatusCode) // add한것으로 StatusCreated가 옴

	// json으로 온것을 읽기
	var todo model.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name,"Test todo")
	id1 := todo.ID

	resp, err = http.PostForm(ts.URL+"/todos", url.Values{"name":{"Test todo2"}})   
	assert.NoError(err)
	assert.Equal(http.StatusCreated,resp.StatusCode) // add한것으로 StatusCreated가 옴
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name,"Test todo2")
	id2 := todo.ID
	// --------------- 2개 저장했음

	// getTodoListHandler 테스트
	resp, err = http.Get(ts.URL+"/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	
	todos := []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)

	for _, t := range todos {
		if t.ID == id1 {
			assert.Equal("Test todo", t.Name)
		} else if t.ID == id2 {
			assert.Equal("Test todo2", t.Name)
		} else {
			assert.Error(fmt.Errorf("testID should be id1 or id2"))
		}		
	}

	
	// completeTodoHandler 테스트
	resp, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// completeTodoHandler 테스트 확인 -- getTodoListHandler 테스트
	resp, err = http.Get(ts.URL+"/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)

	for _, t := range todos {
		if t.ID == id1 {
			assert.True(t.Completed)
		}		
	}

	// removeTodoHandler 테스트
	req, _ := http.NewRequest("DELETE", ts.URL + "/todos/" + strconv.Itoa(id1), nil) 
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	

	// removeTodoHandler 테스트 확인 -- getTodoListHandler 테스트
	resp, err = http.Get(ts.URL+"/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)

	for _, t := range todos {
		assert.Equal(t.ID, id2)
	}
}