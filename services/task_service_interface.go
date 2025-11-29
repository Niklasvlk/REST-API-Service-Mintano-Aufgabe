package services

import "task-api/models"

// TaskServiceInterface definiert die Methoden, die jeder TaskService implementieren muss.
// Dient dazu, unterschiedliche Implementierungen (z.B. echte Service-Logik oder Mocks) austauschbar zu machen.
type TaskServiceInterface interface {
	// CreateTask erstellt einen neuen Task basierend auf der übergebenen CreateTaskRequest.
	// Gibt den gespeicherten Task zurück oder einen Fehler.
	CreateTask(req models.CreateTaskRequest) (*models.Task, error)

	// GetAllTasks gibt alle gespeicherten Tasks zurück.
	// Liefert ein Slice von Tasks oder einen Fehler.
	GetAllTasks() ([]*models.Task, error)

	// GetTaskByID gibt einen Task anhand der ID zurück.
	// Gibt einen Fehler "not found", wenn keine Task mit dieser ID existiert.
	GetTaskByID(id int) (*models.Task, error)

	// UpdateTask aktualisiert einen bestehenden Task anhand der ID und der übergebenen Werte.
	// Nicht gesetzte Felder bleiben unverändert.
	// Liefert den aktualisierten Task oder einen Fehler.
	UpdateTask(id int, req models.CreateTaskRequest) (*models.Task, error)

	// DeleteTask entfernt einen Task anhand der ID.
	// Gibt einen Fehler "not found", falls der Task nicht existiert.
	DeleteTask(id int) error
}
