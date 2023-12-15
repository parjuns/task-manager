package handlers

import (
	"github.com/Parjun2000/task-manager/models"
	"github.com/Parjun2000/task-manager/utils"

	"github.com/Parjun2000/task-manager/helpers"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Credentials
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary		Register a new user
// @Description	Register a new user with username and password
// @Tags			Register & Login
// @Accept			application/json
// @Param			Credentials	body		Credentials				true	"Username & Password"
// @Success		200			{object}	object{message=string}	"User registered successfully"
// @Failure		400			{object}	object{error=string}	"Invalid JSON"
// @Failure		500			{object}	object{error=string}	"Internal server error"
// @Router			/api/v1/auth/register [post]
func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	user.Password = string(hashedPassword)

	s, _ := c.Get("db")
	db := s.(utils.Storage)

	_, err = db.CreateUser(utils.User{Username: user.Username, Password: user.Password})
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": "User registered successfully"})
}

// @Summary		Login user
// @Description	Login user with username and password
// @Tags			Register & Login
// @Accept			application/json
// @Param			Credentials	body		Credentials				true	"Username & Password"
// @Success		200			{object}	object{token=string}	"User logged in successfully"
// @Failure		400			{object}	object{error=string}	"Invalid JSON"
// @Failure		401			{object}	object{error=string}	"Invalid credentials"
// @Failure		404			{object}	object{error=string}	"User not found"
// @Failure		500			{object}	object{error=string}	"Internal server error"
// @Router			/api/v1/auth/login [post]
func Login(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	s, _ := c.Get("db")
	db := s.(utils.Storage)
	storedUser, err := db.GetUserByUsername(user.Username) // Function to get user details from PostgreSQL
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := helpers.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}
