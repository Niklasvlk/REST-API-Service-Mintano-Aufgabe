package services

import (
	"github.com/stretchr/testify/assert"
	"task-api/models"
	"task-api/repository"
	"testing"
)

// Diese Datei enthält Unit-Tests für den TaskService.
// Es werden sowohl erfolgreiche als auch fehlerhafte Szenarien getestet.
// Für alle Tests wird ein MockTaskRepository verwendet, um DB-Zugriffe zu simulieren.

// Test_Service_CreateTask_Success_with_priority_and_status prüft, dass ein Task erfolgreich erstellt wird, wenn Status
// und Priority explizit gesetzt sind.
func Test_Service_CreateTask_Success_with_priority_and_status(t *testing.T) {

	mockRepo := &repository.MockTaskRepository{
		CreateFunc: func(task *models.Task) (*models.Task, error) {
			task.ID = 1
			return task, nil
		},
	}

	service := TaskService{Repo: mockRepo}

	req := models.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Desc",
		Status:      "in progress",
		Priority:    "high",
	}

	task, err := service.CreateTask(req)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, 1, task.ID)
	assert.Equal(t, req.Title, task.Title)
	assert.Equal(t, req.Description, task.Description)
	assert.Equal(t, req.Status, task.Status)
	assert.Equal(t, req.Priority, task.Priority)
}

// Test_Service_CreateTask_Success_without_priority_and_status prüft, dass ein Task erfolgreich erstellt wird,
// wenn Status und Priority leer sind. Default-Werte ("todo" und "medium") werden gesetzt.
func Test_Service_CreateTask_Success_without_priority_and_status(t *testing.T) {
	mockRepo := &repository.MockTaskRepository{
		CreateFunc: func(task *models.Task) (*models.Task, error) {
			task.ID = 1
			return task, nil
		},
	}

	service := TaskService{Repo: mockRepo}

	req := models.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Desc",
	}

	task, err := service.CreateTask(req)

	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, 1, task.ID)
	assert.Equal(t, req.Title, task.Title)
	assert.Equal(t, req.Description, task.Description)
	assert.Equal(t, "todo", task.Status)
	assert.Equal(t, "medium", task.Priority)
}

// Test_Service_GetTaskByID_Success prüft, dass ein Task anhand der ID erfolgreich zurückgegeben wird.
func Test_Service_GetTaskByID_Success(t *testing.T) {
	mockRepo := &repository.MockTaskRepository{
		GetByIdFunc: func(id int) (*models.Task, error) {
			return &models.Task{ID: id, Title: "Test"}, nil
		},
	}

	service := TaskService{Repo: mockRepo}

	task, err := service.GetTaskByID(1)

	assert.Nil(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, 1, task.ID)
}

// Test_Service_GetTaskByID_NotFound prüft, dass ein Fehler zurückgegeben wird, wenn die Task-ID nicht existiert.
func Test_Service_GetTaskByID_NotFound(t *testing.T) {
	mockRepo := &repository.MockTaskRepository{
		GetByIdFunc: func(id int) (*models.Task, error) {
			return nil, nil
		},
	}

	service := TaskService{Repo: mockRepo}

	task, err := service.GetTaskByID(99)

	assert.Nil(t, task)
	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
}

// Test_Service_UpdateTask_Success prüft, dass ein bestehender Task erfolgreich aktualisiert wird.
func Test_Service_UpdateTask_Success(t *testing.T) {
	mockRepo := &repository.MockTaskRepository{
		GetByIdFunc: func(id int) (*models.Task, error) {
			return &models.Task{
				ID:          id,
				Title:       "Alt",
				Description: "Alt",
				Status:      "todo",
				Priority:    "medium",
			}, nil
		},
		UpdateFunc: func(task *models.Task) (*models.Task, error) {
			return task, nil
		},
	}

	service := TaskService{Repo: mockRepo}

	req := models.CreateTaskRequest{
		Title:       "Neu",
		Description: "Neu",
		Status:      "in progress",
		Priority:    "high",
	}

	updated, err := service.UpdateTask(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, 1, updated.ID)
	assert.Equal(t, "Neu", updated.Title)
	assert.Equal(t, "Neu", updated.Description)
	assert.Equal(t, "in progress", updated.Status)
	assert.Equal(t, "high", updated.Priority)
}

// Test_Service_UpdateTask_NotFound prüft, dass ein Fehler zurückgegeben wird, wenn die Task-ID zum
// Update nicht existiert.
func Test_Service_UpdateTask_NotFound(t *testing.T) {
	mockRepo := &repository.MockTaskRepository{
		GetByIdFunc: func(id int) (*models.Task, error) {
			return nil, nil
		},
		UpdateFunc: func(task *models.Task) (*models.Task, error) {
			return task, nil
		},
	}

	service := TaskService{Repo: mockRepo}

	req := models.CreateTaskRequest{
		Title:  "Neu",
		Status: "in progress",
	}

	updated, err := service.UpdateTask(99, req)

	assert.Nil(t, updated)
	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
}

// Test_Service_DeleteTask_Success prüft, dass ein bestehender Task erfolgreich gelöscht wird.
func Test_Service_DeleteTask_Success(t *testing.T) {
	mockRepo := &repository.MockTaskRepository{
		GetByIdFunc: func(id int) (*models.Task, error) {
			return &models.Task{ID: id, Title: "Test Task"}, nil
		},
		DeleteFunc: func(id int) error {
			return nil
		},
	}

	service := TaskService{Repo: mockRepo}

	err := service.DeleteTask(1)
	assert.Nil(t, err)
}

// Test_Service_DeleteTask_NotFound prüft, dass ein Fehler zurückgegeben wird, wenn die Task-ID zum
// Löschen nicht existiert.
func Test_Service_DeleteTask_NotFound(t *testing.T) {
	mockRepo := &repository.MockTaskRepository{
		GetByIdFunc: func(id int) (*models.Task, error) {
			return nil, nil
		},
	}

	service := TaskService{Repo: mockRepo}

	err := service.DeleteTask(1)
	assert.NotNil(t, err)
	assert.Equal(t, "not found", err.Error())
}
