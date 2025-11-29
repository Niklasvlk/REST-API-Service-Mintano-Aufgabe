package models

import "time"

// Task repräsentiert eine Aufgabe in der API.
// Wird sowohl in Responses als auch intern verwendet.
type Task struct {
	ID          int       `json:"id"`          // Eindeutige ID der Task (automatisch vom System vergeben)
	Title       string    `json:"title"`       // Pflichtfeld, max 200 Zeichen
	Description string    `json:"description"` // Optional, max 1000 Zeichen
	Status      string    `json:"status"`      // Status der Task; erlaubt: "todo", "in_progress", "done"
	Priority    string    `json:"priority"`    // Priorität der Task; erlaubt: "low", "medium", "high"
	CreatedAt   time.Time `json:"created_at"`  // Erstellungszeitpunkt
	UpdatedAt   time.Time `json:"updated_at"`  // Letzter Änderungszeitpunkt
}

// CreateTaskRequest repräsentiert die Struktur, die beim Erstellen oder Aktualisieren
// einer Task vom Client an die API geschickt wird.
// Pflichtfeld: Title, optional: Description, Status, Priority.
type CreateTaskRequest struct {
	Title       string `json:"title"`       // Pflichtfeld, max 200 Zeichen
	Description string `json:"description"` // Optional, max 1000 Zeichen
	Status      string `json:"status"`      // Optional, erlaubt: "todo", "in_progress", "done"
	Priority    string `json:"priority"`    // Optional, erlaubt: "low", "medium", "high"
}
