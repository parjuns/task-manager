package middleware

import (
	"github.com/Parjun2000/task-manager/utils"

	"github.com/gin-gonic/gin"
)

func DatabaseMiddleware(s utils.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Passing Database in context
		c.Set("db", s)
		c.Next()
	}
}
