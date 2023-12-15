package middleware

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	// Logging to a file.
	f, _ := os.Create("app.log")
	gin.DefaultWriter = io.MultiWriter(f)
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// custom format
		return fmt.Sprintf("%s - [%s] \n%s\n%s\n%s\n%d\n%s\n\"%s\"\n%s\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}
