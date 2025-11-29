package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"task-api/models"
)

// PostgresTaskRepository implementiert die Persistenzschicht für Tasks
// und kapselt alle CRUD-Operationen gegen eine PostgreSQL-Datenbank.
type PostgresTaskRepository struct {
	DB *sql.DB
}

// Create speichert einen neuen Task in der Datenbank.
// Gibt den vollständigen Task inklusive ID, CreatedAt und UpdatedAt zurück.
func (r *PostgresTaskRepository) Create(task *models.Task) (*models.Task, error) {
	query := `INSERT INTO tasks (title, description, status, priority)
	          VALUES ($1, $2, $3, $4)
	          RETURNING id, created_at, updated_at`

	err := r.DB.QueryRow(query, task.Title, task.Description, task.Status, task.Priority).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetAll gibt alle Tasks aus der Datenbank zurück.
// Liefert ein Slice von Task-Pointern oder einen Fehler.
func (r *PostgresTaskRepository) GetAll() ([]*models.Task, error) {
	rows, err := r.DB.Query(`SELECT id, title, description, status, priority, created_at, updated_at FROM tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		t := &models.Task{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// GetByID gibt einen Task anhand der ID zurück.
// Gibt nil zurück, wenn kein Task mit der ID existiert.
func (r *PostgresTaskRepository) GetByID(id int) (*models.Task, error) {
	query := `SELECT id, title, description, status, priority, created_at, updated_at
	          FROM tasks WHERE id=$1`

	task := &models.Task{}
	err := r.DB.QueryRow(query, id).Scan(
		&task.ID, &task.Title, &task.Description, &task.Status,
		&task.Priority, &task.CreatedAt, &task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return task, nil
}

// Update ändert die Felder eines bestehenden Tasks in der Datenbank.
// Gibt den aktualisierten Task zurück oder einen Fehler.
func (r *PostgresTaskRepository) Update(task *models.Task) (*models.Task, error) {
	query := `UPDATE tasks 
              SET title=$1, status=$2, priority=$3, updated_at=NOW()
              WHERE id=$4
              RETURNING id, title, status, priority, created_at, updated_at`

	err := r.DB.QueryRow(query, task.Title, task.Status, task.Priority, task.ID).
		Scan(&task.ID, &task.Title, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Delete entfernt einen Task anhand der ID aus der Datenbank.
// Gibt einen Fehler zurück, falls die Löschung fehlschlägt.
func (r *PostgresTaskRepository) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}
