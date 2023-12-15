package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Parjun2000/task-manager/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testRouter = gin.Default()

func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}
func TestAuthMiddleware(t *testing.T) {
	// Create a test Gin router
	testRouter.Use(AuthMiddleware())
	testRouter.GET("/tasks", func(c *gin.Context) {
		c.String(200, "Authorized")
	})

	// Perform a GET request to the protected endpoint without authentication
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// Assert response status code (should be unauthorized)
	assert.Equal(t, 401, w.Code)
	// Add more assertions as needed
}

func setup() {
	testRouter.Use(mockDB())
}

func teardown() {
	//
}

func mockDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Passing MockDatabase in context
		s := utils.NewMockDB()
		c.Set("db", s)
		c.Next()
	}
}
