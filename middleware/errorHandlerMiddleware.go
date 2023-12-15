package middleware

import (
	"github.com/gin-gonic/gin"
)

// @Failure	500	{object}	object{error=string}	"Internal Server Error"
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
			}
		}()

		c.Next()
	}
}
