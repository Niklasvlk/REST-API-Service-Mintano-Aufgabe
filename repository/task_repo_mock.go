package repository

import "task-api/models"

// MockTaskRepository ist ein Mock des TaskRepositoryInterface für Tests.
// Jede Methode wird durch eine Funktion ersetzt, die individuell gesetzt werden kann.
// So kann man gewünschtes Verhalten für Unit-Tests simulieren.
type MockTaskRepository struct {
	// CreateFunc simuliert das Erstellen eines Tasks.
	CreateFunc func(task *models.Task) (*models.Task, error)

	// GetAllFunc simuliert das Abrufen aller Tasks.
	GetAllFunc func() ([]*models.Task, error)

	// GetByIdFunc simuliert das Abrufen eines Tasks anhand der ID.
	GetByIdFunc func(id int) (*models.Task, error)

	// UpdateFunc simuliert das Aktualisieren eines Tasks.
	UpdateFunc func(task *models.Task) (*models.Task, error)

	// DeleteFunc simuliert das Löschen eines Tasks anhand der ID.
	DeleteFunc func(id int) error
}

// Create ruft CreateFunc auf und gibt das Ergebnis zurück.
func (m *MockTaskRepository) Create(task *models.Task) (*models.Task, error) {
	return m.CreateFunc(task)
}

// GetAll ruft GetAllFunc auf und gibt das Ergebnis zurück.
func (m *MockTaskRepository) GetAll() ([]*models.Task, error) {
	return m.GetAllFunc()
}

// GetByID ruft GetByIdFunc auf und gibt das Ergebnis zurück.
func (m *MockTaskRepository) GetByID(id int) (*models.Task, error) {
	return m.GetByIdFunc(id)
}

// Update ruft UpdateFunc auf und gibt das Ergebnis zurück.
func (m *MockTaskRepository) Update(task *models.Task) (*models.Task, error) {
	return m.UpdateFunc(task)
}

// Delete ruft DeleteFunc auf und gibt das Ergebnis zurück.
func (m *MockTaskRepository) Delete(id int) error {
	return m.DeleteFunc(id)
}
