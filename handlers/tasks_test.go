package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetTasks(t *testing.T) {
	// Create a GET route
	testRouter.GET("/tasks", GetTasks)
	testRouter.Use(func(ctx *gin.Context) {
		ctx.Set("user_id", 1)
	})

	// This create a mock request to test the handler
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// Assert the HTTP status code
	assert.Equal(t, 500, recorder.Code)

}

func TestCreateTask(t *testing.T) {
	// Create a POST route
	testRouter.POST("/tasks", CreateTask)

	// This create a mock request to test the handler
	req, err := http.NewRequest("POST", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// Assert the HTTP status code
	assert.Equal(t, 400, recorder.Code)
}

func TestGetTaskByID(t *testing.T) {
	// Create a GET route
	testRouter.GET("/tasks/:id", GetTaskByID)

	// This create a mock request to test the handler
	req, err := http.NewRequest("GET", "/tasks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// Assert the HTTP status code
	assert.Equal(t, 400, recorder.Code)
}

func TestUpdateTask(t *testing.T) {
	// Create a PUT route
	testRouter.PUT("/tasks/:id", UpdateTask)

	// This create a mock request to test the handler
	req, err := http.NewRequest("PUT", "/tasks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// Assert the HTTP status code
	assert.Equal(t, 400, recorder.Code)
}

func TestDeleteTask(t *testing.T) {
	// Create a DELETE route
	testRouter.DELETE("/tasks/:id", DeleteTask)

	// This create a mock request to test the handler
	req, err := http.NewRequest("DELETE", "/tasks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// Assert the HTTP status code
	assert.Equal(t, 400, recorder.Code)
}

func TestMarkTasksDoneConcurrently(t *testing.T) {
	// Create a PUT route
	testRouter.PUT("/tasks/mark-done", MarkTasksDoneConcurrently)

	// This create a mock request to test the handler
	req, err := http.NewRequest("PUT", "/tasks/mark-done", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// Assert the HTTP status code
	assert.Equal(t, 400, recorder.Code)
}
