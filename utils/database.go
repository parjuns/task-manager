package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetUserByID(int) (User, error)
	GetUserByUsername(string) (User, error)
	CreateUser(User) (int, error)
	GetTasksWithParams(userId, page, limit int, sortBy, order, status string) ([]Task, error)
	GetTaskByID(userId, id int) (Task, error)
	CreateTask(newTask Task) (int, error)
	UpdateTaskStatusDone(userID, taskID int) error
	UpdateTaskByID(userID, taskID int, updatedTask Task) error
	DeleteTask(userID, id int) error
}
type PostgresDB struct {
	DB *sql.DB
}

// Create PostgresDB
func NewPostgresDB(db *sql.DB) *PostgresDB {
	return &PostgresDB{
		DB: db,
	}
}

// InitDB initializes the database connection.
func InitDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	store, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	// db = conn

	if err = store.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to database")

	return store, nil
}

// CloseDB closes the database connection.
func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
		log.Println("Database connection closed")
	}
}

// Task represents the task structure.
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int       `json:"user_id"`
}

// User represents the user structure.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetUserByID retrieves a user by ID from the database.
func (s *PostgresDB) GetUserByID(id int) (User, error) {
	var user User
	row := s.DB.QueryRow("SELECT id, username, password FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username from the database.
func (s *PostgresDB) GetUserByUsername(username string) (User, error) {
	var user User
	row := s.DB.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}
	return user, nil
}

// CreateUser creates a new user in the database.
func (s *PostgresDB) CreateUser(newUser User) (int, error) {
	var id int
	err := s.DB.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id",
		newUser.Username, newUser.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetTasksWithParmas retrieves all tasks from the database with offset,limit,sort and order.
func (s *PostgresDB) GetTasksWithParams(userId, page, limit int, sortBy, order, status string) ([]Task, error) {

	offset := (page - 1) * limit

	query := "SELECT id, title, description, status, created_at,user_id FROM tasks"
	query += " Where user_id = $1"
	if status != "" {
		query += " and status = $4"
	}
	query += " ORDER BY " + sortBy + " " + order
	query += " LIMIT $2 OFFSET $3"

	var rows *sql.Rows
	var err error

	if status != "" {
		rows, err = s.DB.Query(query, userId, limit, offset, status)
	} else {
		rows, err = s.DB.Query(query, userId, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UserID); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTaskByID retrieves a task by its ID from the database.
func (s *PostgresDB) GetTaskByID(userId, id int) (Task, error) {
	var task Task
	row := s.DB.QueryRow("SELECT id, title, description, status, created_at, user_id FROM tasks WHERE id = $1 and user_id = $2", id, userId)
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, errors.New("task not found")
		}
		return Task{}, err
	}
	return task, nil
}

// CreateTask creates a new task in the database.
func (s *PostgresDB) CreateTask(newTask Task) (int, error) {
	var id int
	err := s.DB.QueryRow("INSERT INTO tasks (title, description, status, created_at, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		newTask.Title, newTask.Description, newTask.Status, time.Now(), newTask.UserID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateTaskStatusDone updates the status of an existing task in the database with 'done'.
func (s *PostgresDB) UpdateTaskStatusDone(userID, taskID int) error {
	_, err := s.DB.Exec("UPDATE tasks SET status = $1 WHERE id = $2 and user_id = $3", "done", taskID, userID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTaskStatus updates the status of an existing task in the database by its ID.
func (s *PostgresDB) UpdateTaskByID(userID, taskID int, updatedTask Task) error {
	_, err := s.DB.Exec("UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4 and user_id=$5", updatedTask.Title, updatedTask.Description, updatedTask.Status, taskID, userID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTask deletes a task by its ID from the database.
func (s *PostgresDB) DeleteTask(userID, id int) error {
	result, err := s.DB.Exec("DELETE FROM tasks WHERE id = $1 and user_id= $2", id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("task not found")
	}
	return nil
}
