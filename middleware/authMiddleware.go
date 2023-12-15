package middleware

import (
	"os"
	"time"

	"github.com/Parjun2000/task-manager/helpers"
	"github.com/Parjun2000/task-manager/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

// @Failure	401	{object}	object{error=string}	"Unauthorized: Missing token"
// @Failure	401	{object}	object{error=string}	"Unauthorized: Invalid token"
// @Failure	401	{object}	object{error=string}	"Unauthorized"
// @Failure	403	{object}	object{error=string}	"Unauthorized: Token expired"
// @Failure	500	{object}	object{error=string}	"Internal Server Error"
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Unauthorized: Missing token"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &helpers.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(401, gin.H{"error": "Unauthorized: Invalid token"})
				c.Abort()
				return
			}
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(401, gin.H{"error": "Unauthorized: Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*helpers.Claims)
		if !ok {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if time.Now().Unix() > claims.ExpiresAt {
			c.JSON(403, gin.H{"error": "Unauthorized: Token expired"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		s, _ := c.Get("db")
		user, err := s.(utils.Storage).GetUserByUsername(claims.Username)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}
		c.Set("user_id", user.ID)
		c.Next()
	}
}
