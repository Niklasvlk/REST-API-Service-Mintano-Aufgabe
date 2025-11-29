package services

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"task-api/models"
	"time"
)

// MockTaskService implementiert TaskServiceInterface für Tests.
// Dient dazu, das Verhalten des TaskService zu simulieren, ohne eine echte Datenbank zu verwenden.
// Felder:
// - Tasks: Vordefinierte Tasks für Tests.
// - Err: Optionaler Fehler, der bei GetTaskByID zurückgegeben wird.
// - ShouldFail: Wenn true, schlagen alle Methoden absichtlich fehl.
type MockTaskService struct {
	Tasks      []*models.Task
	Err        error
	ShouldFail bool
}

// CreateTask simuliert das Erstellen eines Tasks.
// Gibt einen Task zurück oder einen internen Serverfehler, wenn ShouldFail=true ist.
func (m *MockTaskService) CreateTask(req models.CreateTaskRequest) (*models.Task, error) {
	if m.ShouldFail {
		return nil, fiber.ErrInternalServerError
	}
	return &models.Task{
		ID:          1,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
	}, nil
}

// GetAllTasks gibt alle Tasks im Mock zurück.
// Liefert einen Fehler, wenn ShouldFail=true ist.
func (m *MockTaskService) GetAllTasks() ([]*models.Task, error) {
	if m.ShouldFail {
		return nil, fiber.ErrInternalServerError
	}
	return m.Tasks, nil
}

// GetTaskByID gibt einen Task anhand der ID zurück.
// Liefert "not found", wenn kein Task existiert oder m.Err gesetzt ist.
func (m *MockTaskService) GetTaskByID(id int) (*models.Task, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	for _, t := range m.Tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

// UpdateTask simuliert das Aktualisieren eines Tasks.
// Felder, die im Request leer sind, bleiben unverändert.
// Liefert "not found" oder internen Serverfehler je nach Konfiguration.
func (m *MockTaskService) UpdateTask(id int, req models.CreateTaskRequest) (*models.Task, error) {
	if m.ShouldFail {
		return nil, fiber.ErrInternalServerError
	}

	var task *models.Task
	for _, t := range m.Tasks {
		if t.ID == id {
			task = t
			break
		}
	}
	if task == nil {
		return nil, fmt.Errorf("not found")
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	task.UpdatedAt = time.Now()

	return task, nil
}

// DeleteTask simuliert das Löschen eines Tasks anhand der ID.
// Liefert "not found", wenn kein Task existiert, oder einen Fehler, wenn ShouldFail=true ist.
func (m *MockTaskService) DeleteTask(id int) error {
	if m.ShouldFail {
		return fiber.ErrInternalServerError
	}

	for i, t := range m.Tasks {
		if t.ID == id {
			m.Tasks = append(m.Tasks[:i], m.Tasks[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("not found")
}
