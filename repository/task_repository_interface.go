package repository

import "task-api/models"

// TaskRepositoryInterface definiert die CRUD-Methoden, die jedes Repository implementieren muss.
type TaskRepositoryInterface interface {
	// Create speichert einen neuen Task und gibt den vollständigen Task zurück.
	Create(task *models.Task) (*models.Task, error)

	// GetAll gibt alle gespeicherten Tasks zurück.
	GetAll() ([]*models.Task, error)

	// GetByID gibt einen Task anhand seiner ID zurück.
	// Gibt nil, nil zurück, wenn kein Task gefunden wird.
	GetByID(id int) (*models.Task, error)

	// Update aktualisiert einen bestehenden Task und gibt den aktualisierten Task zurück.
	Update(task *models.Task) (*models.Task, error)

	// Delete entfernt einen Task anhand seiner ID.
	Delete(id int) error
}
