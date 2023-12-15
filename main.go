package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/Parjun2000/task-manager/docs"

	"github.com/Parjun2000/task-manager/handlers"
	"github.com/Parjun2000/task-manager/middleware"
	"github.com/Parjun2000/task-manager/utils"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title						Go-Gin Task Manager API
//	@version					1.0
//	@description				This API allow users to create, read, update, and delete tasks. Users can register, log in, and access their own tasks only when authenticated.
//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
func main() {

	//Load environment varibles
	if err := godotenv.Load("app.env"); err != nil {
		log.Panic("Error loading .env file")
	}

	// Initialize database connection
	db, err := utils.InitDB()
	if err != nil {
		log.Fatal("Database Initialize Error: ", err)
	}
	store := utils.NewPostgresDB(db)
	defer db.Close()

	// Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.DatabaseMiddleware(store))

	v1 := router.Group("/api/v1")
	// Auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Protected Tasks Routes
	tasks := v1.Group("/tasks")
	tasks.Use(middleware.AuthMiddleware())
	{
		tasks.GET("/", handlers.GetTasks)
		tasks.POST("/", handlers.CreateTask)
		tasks.GET("/:id", handlers.GetTaskByID)
		tasks.PUT("/:id", handlers.UpdateTask)
		tasks.DELETE("/:id", handlers.DeleteTask)
		tasks.PUT("/mark-done", handlers.MarkTasksDoneConcurrently)
	}

	// Swagger documentation route
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize Server
	port := os.Getenv("SERVER_PORT")
	err = router.Run(fmt.Sprint(":", port))
	if err != nil {
		log.Fatal("Server Initialize Error: ", err)
	}
}
