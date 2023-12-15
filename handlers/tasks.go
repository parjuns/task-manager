package handlers

import (
	"errors"
	"strconv"
	"sync"

	"github.com/Parjun2000/task-manager/models"
	"github.com/Parjun2000/task-manager/utils"

	"github.com/gin-gonic/gin"
)

// TaskDetails
type TaskDetails struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// @Summary		Get tasks with pagination, sorting, and filtering
// @Description	Get tasks with pagination, sorting by status/created_at, and filtering by status
// @Tags			Tasks
// @Accept			application/json
// @Produce		application/json
// @Security		JWT
// @Param			page	query		int		false	"Page number"
// @Param			limit	query		int		false	"Items per page"
// @Param			sort_by	query		string	false	"Sort by title/status/description/created_at"
// @Param			order	query		string	false	"Sort order: asc/desc"
// @Param			status	query		string	false	"Filter by task status"
// @Success		200		{array}		utils.Task
// @Failure		400		{object}	object{error=string}	"Error Message"
// @Failure		500		{object}	object{error=string}	"Internal Server Error"
// @Failure		500		{object}	object{error=string}	"Failed to fetch tasks:Error"
// @Router			/api/v1/tasks [get]
func GetTasks(c *gin.Context) {
	page, limit, sortBy, order, status, err := extractPaginationParams(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	s, _ := c.Get("db")
	db := s.(utils.Storage)
	tasks, err := db.GetTasksWithParams(userId.(int), page, limit, sortBy, order, status)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch tasks" + err.Error()})
		return
	}

	c.JSON(200, tasks)
}

// Extract parameters for pagination, sorting, and filtering
func extractPaginationParams(c *gin.Context) (int, int, string, string, string, error) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		return 0, 0, "", "", "", errors.New("invalid page number")
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		return 0, 0, "", "", "", errors.New("invalid limit")
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")
	status := c.DefaultQuery("status", "")

	validSortOptions := map[string]bool{
		"title":       true,
		"description": true,
		"status":      true,
		"created_at":  true,
	}
	if !validSortOptions[sortBy] {
		return 0, 0, "", "", "", errors.New("invalid sort option")
	}

	return page, limit, sortBy, order, status, nil
}

// @Summary		Create a task
// @Description	Create a new task
// @Tags			Tasks
// @Accept			application/json
// @Security		JWT
// @Param			TaskDetails	body		TaskDetails				true	"Task Details"
// @Success		200			{object}	object{message=string}	"Task created successfully"
// @Failure		400			{object}	object{error=string}	"Invalid JSON"
// @Failure		400			{object}	object{error=string}	"Validation Error"
// @Failure		500			{object}	object{error=string}	"Internal Server Error"
// @Router			/api/v1/tasks [post]
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := task.Validate(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	var newTask = utils.Task{
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		UserID:      userId.(int),
	}

	s, _ := c.Get("db")
	db := s.(utils.Storage)
	taskId, err := db.CreateTask(newTask)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.Set("task_id", taskId)
	c.JSON(200, gin.H{"message": "Task created successfully"})
}

// @Summary		Get task by ID
// @Description	Get a task by its ID
// @Tags			Tasks
// @Produce		application/json
// @Security		JWT
// @Param			id	path		int	true	"Task ID"
// @Success		200	{object}	utils.Task
// @Failure		404	{object}	object{error=string}	"Task not found"
// @Failure		500	{object}	object{error=string}	"Internal Server Error"
// @Router			/api/v1/tasks/{id} [get]
func GetTaskByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Task Id"})
		return
	}

	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	s, _ := c.Get("db")
	db := s.(utils.Storage)
	task, err := db.GetTaskByID(userId.(int), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(200, task)
}

// @Summary		Update a task
// @Description	Update an existing task by ID
// @Tags			Tasks
// @Accept			application/json
// @Security		JWT
// @Param			id			path		int						true	"Task ID"
// @Param			TaskDetails	body		TaskDetails				true	"Updated Task Details"
// @Success		200			{object}	object{message=string}	"Task updated successfully"
// @Failure		400			{object}	object{error=string}	"Invalid JSON"
// @Failure		400			{object}	object{error=string}	"Invalid Task Id"
// @Failure		404			{object}	object{error=string}	"Task not found"
// @Failure		500			{object}	object{error=string}	"Internal Server Error"
// @Failure		500			{object}	object{error=string}	"Failed to update task"
// @Router			/api/v1/tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Task Id"})
		return
	}
	var updatedTask models.Task
	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := updatedTask.Validate(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	s, _ := c.Get("db")
	db := s.(utils.Storage)
	if _, err := db.GetTaskByID(userId.(int), id); err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	var task = utils.Task{
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Status:      updatedTask.Status,
	}

	if err = db.UpdateTaskByID(userId.(int), id, task); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update task"})
		return
	}
	c.JSON(200, gin.H{"message": "Task updated successfully"})
}

// @Summary		Delete a task
// @Description	Delete a task by ID
// @Tags			Tasks
// @Accept			application/json
// @Security		JWT
// @Param			id	path		int						true	"Task ID"
// @Success		200	{object}	object{message=string}	"Task deleted successfully"
// @Failure		400	{object}	object{error=string}	"Invalid Task Id"
// @Failure		404	{object}	object{error=string}	"Task not found"
// @Failure		500	{object}	object{error=string}	"Internal Server Error"
// @Failure		500	{object}	object{error=string}	"Failed to delete task"
// @Router			/api/v1/tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Task Id"})
		return
	}
	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	s, _ := c.Get("db")
	db := s.(utils.Storage)
	if _, err := db.GetTaskByID(userId.(int), id); err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}
	if err := db.DeleteTask(userId.(int), id); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete task"})
		return
	}
	c.JSON(200, gin.H{"message": "Task deleted successfully"})
}

// MarkTasksDoneConcurrently marks multiple tasks as done concurrently
//
//	@Summary		Mark tasks as done concurrently
//	@Description	Mark multiple tasks as done concurrently using Goroutines
//	@Tags			Tasks
//	@Accept			application/json
//	@Produce		application/json
//	@Security		JWT
//	@Param			task_ids	body		[]string	true	"Task IDs to mark as done"
//	@Success		200			{object}	object{message=string,updated_tasks=[]string}
//	@Failure		400			{object}	object{error=string}	"Invalid Task Id"
//	@Failure		500			{object}	object{error=string}	"Internal Server Error"
//	@Failure		500			{object}	object{error=string}	"Failed to Update task"
//	@Router			/api/v1/tasks/mark-done [put]
func MarkTasksDoneConcurrently(c *gin.Context) {
	var taskIDs []string
	if err := c.BindJSON(&taskIDs); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	var wg sync.WaitGroup
	resultCh := make(chan string)

	// Launch Goroutines to mark tasks as done concurrently
	for _, taskID := range taskIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			task_id, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid Task Id"})
				return
			}

			s, _ := c.Get("db")
			db := s.(utils.Storage)
			if err := db.UpdateTaskStatusDone(userId.(int), task_id); err != nil {
				c.JSON(500, gin.H{"error": "Failed to Update task"})
				return
			}
			resultCh <- id
		}(taskID)
	}

	// Wait for all Goroutines to finish
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect results from Goroutines
	updatedTasks := make([]string, 0)
	for id := range resultCh {
		updatedTasks = append(updatedTasks, id)
	}

	c.JSON(200, gin.H{"message": "Tasks marked as done concurrently", "updated_tasks": updatedTasks})
}
