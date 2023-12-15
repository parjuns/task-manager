package handlers

import (
	"os"
	"testing"

	"github.com/Parjun2000/task-manager/utils"
	"github.com/gin-gonic/gin"
)

var testRouter = gin.Default()

func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
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
