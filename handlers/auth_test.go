package handlers

import (
	"bytes"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	// Create a POST route
	testRouter.POST("/auth/register", Register)

	// Perform a POST request to the endpoint without body
	req, err := http.NewRequest("POST", "/auth/register", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// Assert response status code
	assert.Equal(t, 400, recorder.Code)

	// Perform a POST request to the endpoint with body payload
	body := []byte(`{"password":"pass1","username":"user1"}`)
	req, err = http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// Assert response status code
	assert.Equal(t, 200, recorder.Code)
}

func TestLogin(t *testing.T) {
	// Create a POST route
	testRouter.POST("/auth/login", Login)

	// Perform a POST request to the endpoint without body
	req, err := http.NewRequest("POST", "/auth/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// Assert response status code
	assert.Equal(t, 400, recorder.Code)

	// Perform a POST request to the endpoint with body payload
	body := []byte(`{"password":"pass1","username":"user1"}`)
	req, err = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// Assert response status code
	assert.Equal(t, 401, recorder.Code)
}
