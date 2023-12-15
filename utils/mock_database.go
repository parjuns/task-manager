package utils

import "time"

type MockDB struct {
	Tasks []Task
	Users []User
}

func NewMockDB() *MockDB {
	return &MockDB{
		Tasks: []Task{{1, "title1", "description1", "todo", time.Now(), 1}},
		Users: []User{{1, "user1", "pass1"}},
	}
}
func (m *MockDB) GetUserByID(id int) (User, error) {
	return User{
		ID:       id,
		Username: "randomuser",
		Password: "randompass",
	}, nil
}
func (m *MockDB) GetUserByUsername(username string) (User, error) {
	return User{
		ID:       2,
		Username: username,
		Password: "randompass",
	}, nil
}
func (m *MockDB) CreateUser(newUser User) (int, error) {
	return 1, nil
}
func (m *MockDB) GetTasksWithParams(userId, page, limit int, sortBy, order, status string) ([]Task, error) {
	return m.Tasks, nil
}
func (m *MockDB) GetTaskByID(userId, id int) (Task, error) {
	return m.Tasks[0], nil
}
func (m *MockDB) CreateTask(newTask Task) (int, error) {
	m.Tasks = append(m.Tasks, newTask)
	return newTask.ID, nil
}
func (m *MockDB) UpdateTaskStatusDone(userID, taskID int) error {
	return nil
}
func (m *MockDB) UpdateTaskByID(userID, taskID int, updatedTask Task) error {
	return nil
}
func (m *MockDB) DeleteTask(userID, id int) error {
	return nil
}
