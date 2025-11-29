package services

import (
	"fmt"
	"task-api/models"
	"task-api/repository"
	"time"
)

// TaskService kapselt die Businesslogik für Tasks.
// Nutzt ein Repository (Postgres), um Daten zu speichern und abzurufen.
// Verantwortlich für Default-Werte und Fehlerbehandlung.
type TaskService struct {
	Repo repository.TaskRepositoryInterface
}

// CreateTask erstellt einen neuen Task anhand der übergebenen CreateTaskRequest.
// Setzt Default-Werte: Status="todo", Priority="medium", falls nicht angegeben.
// Gibt den gespeicherten Task zurück oder einen Fehler.
func (s *TaskService) CreateTask(req models.CreateTaskRequest) (*models.Task, error) {

	// Default Status/Priority
	if req.Status == "" {
		req.Status = "todo"
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
	}

	return s.Repo.Create(task)
}

// GetAllTasks gibt alle gespeicherten Tasks zurück.
// Gibt ein Slice von Tasks oder einen Fehler zurück.
func (s *TaskService) GetAllTasks() ([]*models.Task, error) {
	return s.Repo.GetAll()
}

// GetTaskByID gibt einen Task anhand der ID zurück.
// Gibt einen Fehler "not found", wenn keine Task existiert.
func (s *TaskService) GetTaskByID(id int) (*models.Task, error) {
	task, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("not found")
	}
	return task, nil
}

// UpdateTask aktualisiert einen bestehenden Task anhand der ID und der neuen Werte.
// Felder, die im Request leer bleiben, werden nicht verändert.
// Setzt UpdatedAt auf die aktuelle Zeit.
// Gibt den aktualisierten Task zurück oder einen Fehler.
func (s *TaskService) UpdateTask(id int, req models.CreateTaskRequest) (*models.Task, error) {
	task, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, err
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

	updatedTask, err := s.Repo.Update(task)
	if err != nil {
		return nil, err
	}

	return updatedTask, nil
}

// DeleteTask entfernt einen Task anhand der ID.
// Gibt einen Fehler "not found", falls der Task nicht existiert.
func (s *TaskService) DeleteTask(id int) error {
	task, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}
	if task == nil {
		return fmt.Errorf("not found")
	}

	return s.Repo.Delete(id)
}
