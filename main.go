package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"task-api/handlers"
	"task-api/repository"
	"task-api/services"
)

// main ist der Einstiegspunkt der Anwendung.
// - Stellt die PostgreSQL-Datenbankverbindung her
// - Initialisiert Repository-, Service- und Handler-Layer
// - Registriert alle HTTP-Routen
// - Startet den Fiber Webserver unter Port 8080
func main() {
	app := fiber.New()

	// Erstellen des Connection-Strings für Postgres.
	// Werte werden über Umgebungsvariablen eingelesen.
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	// Öffnet die Datenbankverbindung
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Testet, ob die Verbindung möglich ist
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Dependency-Injection:
	// Repository -> Service -> Handler
	repo := &repository.PostgresTaskRepository{DB: db}
	service := &services.TaskService{Repo: repo}
	handler := &handlers.TaskHandler{Service: service}

	// ---------------------- ROUTES ----------------------
	// POST /tasks  -> Erstellt einen neuen Task
	app.Post("/tasks", handler.CreateTask)

	// GET /tasks -> Liefert eine Liste aller Tasks zurück
	app.Get("/tasks", handler.GetAllTasks)

	// GET /tasks/:id -> Liefert einen Task anhand seiner ID zurück
	app.Get("/tasks/:id", handler.GetTaskByID)

	// PUT /tasks/:id -> Aktualisiert einen bestehenden Task
	app.Put("/tasks/:id", handler.UpdateTask)

	// DELETE /tasks/:id -> Löscht einen Task anhand der ID
	app.Delete("/tasks/:id", handler.DeleteTask)

	// GET /health -> Health Check Endpoint, liefert "OK"
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Startet Server unter Port 8080 (Blockierend)
	log.Fatal(app.Listen(":8080"))
}
