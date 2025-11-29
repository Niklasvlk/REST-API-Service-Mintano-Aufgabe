package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"task-api/models"
	"task-api/services"
)

// TaskHandler stellt die HTTP-Schicht dar und verbindet eingehende Requests
// mit der Businesslogik im TaskService. Jeder Handler entspricht einem API-Endpoint.
type TaskHandler struct {
	Service services.TaskServiceInterface
}

// Erlaubte Priorities für die Validierung der Task
// "" bedeutet kein gesetzter Wert
var allowedPriorities = map[string]bool{
	"low":    true,
	"medium": true,
	"high":   true,
	"":       true,
}

// Erlaubte Status für die Validierung der Task
// "" bedeutet kein gesetzter Wert
var allowedStatus = map[string]bool{
	"todo":        true,
	"in progress": true,
	"done":        true,
	"":            true,
}

// CreateTask verarbeitet POST /tasks.
// Erwartet einen JSON-Body mit Task-Daten. Diese muss nur zwingend einen Titel beinhalten.
// Antwort:
//
//	201 - Task erfolgreich erstellt (JSON)
//	400 - Fehlerhafte Anfrage / Validierungsfehler / Serverfehler beim Erstellen der Task
//
// Beispiel Request-Body:
//
//	{
//	  "Title": "Einkaufen",
//	  "Description": "Einkaufen gehen",
//	}
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	// Request Body einlesen & JSON → Struct parsen
	var req models.CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": "invalid request body",
		})
	}

	// Validierung
	if req.Title == "" || len(req.Title) > 200 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": "Title is required and must be max 200 characters",
		})
	}

	if len(req.Description) > 1000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": "Description must be max 1000 characters",
		})
	}

	if !allowedPriorities[req.Priority] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": "Priority must be one of: low, medium, high or nothing",
		})
	}

	if !allowedStatus[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": "Status must be one of: todo, in progress, done or nothing",
		})
	}

	// Service übernimmt persistente Logik (Clean Architecture)
	task, err := h.Service.CreateTask(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

// GetAllTasks verarbeitet GET /tasks.
// Gibt eine Liste aller gespeicherten Tasks zurück.
// Antwort:
//
//	200 - OK + Array von Tasks (Ohne die Description)
//	400 - Fehler beim Laden aus der Datenbank
func (h *TaskHandler) GetAllTasks(c *fiber.Ctx) error {
	// Ruft alle Tasks über den Service ab
	tasks, err := h.Service.GetAllTasks()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "internal error",
			"message": err.Error(),
		})
	}

	// Wandelt Task-Model in API-Response konformes JSON-Objekt um
	var respTasks []fiber.Map
	for _, t := range tasks {
		respTasks = append(respTasks, fiber.Map{
			"id":         t.ID,
			"title":      t.Title,
			"status":     t.Status,
			"priority":   t.Priority,
			"created at": t.CreatedAt,
		})
	}

	// Erfolgreiche Antwort → gibt Liste aller Tasks + Gesamtanzahl zurück
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tasks": respTasks,
		"total": len(respTasks),
	})
}

// GetTaskByID verarbeitet GET /tasks/:id.
// Parameter:
//
//	path: id (string/int) → ID des gewünschten Tasks
//
// Antwort:
//
//	200 - Task gefunden (JSON)
//	404 - Keine Task mit dieser ID vorhanden
func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	// Liest die ID aus der URL und wandelt sie in einen Integer um
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": "ID must be an integer",
		})
	}

	// Holt den Task über den Service anhand der ID
	task, err := h.Service.GetTaskByID(id)
	if err != nil {
		if err.Error() == "not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "not found",
				"message": fmt.Sprintf("Task with ID %d not found", id),
			})
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	// Erfolgreiche Antwort → gibt spezifischen Task zurück
	return c.Status(fiber.StatusOK).JSON(task)
}

// UpdateTask verarbeitet PUT /tasks/:id.
// Erwartet JSON-Body mit neuen Werten für den Task. Falls ein neuer Wert gesetzt ist, wird dieser in der Task aktualisiert.
//
// Antwort:
//
//	200 - Erfolgreich aktualisiert + neuer Task
//	400 - Ungültige Daten / Fehler beim Update
//	404 - Task nicht gefunden
//
// Beispiel Request-Body:
//
//	{
//	  "title": "Neue Beschreibung",
//	}
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	// Liest die ID aus der URL und wandelt sie in einen Integer um
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": "invalid id",
		})
	}

	// Validierung
	var req models.CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": "invalid request body",
		})
	}

	if len(req.Title) > 200 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": "Title is required and must be max 200 characters",
		})
	}

	if len(req.Description) > 1000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": "Description must be max 1000 characters",
		})
	}

	if !allowedPriorities[req.Priority] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": "Priority must be one of: low, medium, high or nothing",
		})
	}

	if !allowedStatus[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": "Status must be one of: todo, in progress, done or nothing",
		})
	}

	// Update der Task über den Service
	updatedTask, err := h.Service.UpdateTask(id, req)
	if err != nil {
		if err.Error() == "not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "not found",
				"message": fmt.Sprintf("Task with ID %d not found", id),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "internal error",
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(updatedTask)
}

// DeleteTask verarbeitet DELETE /tasks/:id.
// Löscht einen Task anhand seiner ID.
//
// Antwort:
//
//	204 - Erfolgreich gelöscht (Kein Body)
//	400 - Fehler beim Löschen
//	404 - Task existiert nicht
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	// Liest ID aus der URL und validiert sie als Integer
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "validation error",
			"message": "Invalid ID",
		})
	}

	// Service ruft Löschvorgang für den Task auf
	err = h.Service.DeleteTask(id)
	if err != nil {
		if err.Error() == "not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "not found",
				"message": "Task with ID " + idParam + " not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "internal error",
			"message": err.Error(),
		})
	}

	// Erfolgreich gelöscht → 204 No Content
	return c.SendStatus(fiber.StatusNoContent)
}
