package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"strings"
	"task-api/services"
	"testing"
	"time"

	"task-api/handlers"
	"task-api/models"

	"github.com/gofiber/fiber/v2"
)

// setupFiberHandler initialisiert einen Fiber-App-Server mit allen TaskHandler-Routen
// (POST /tasks, GET /tasks, GET /tasks/:id, PUT /tasks/:id, DELETE /tasks/:id) unter Verwendung eines Mock-Service.
func setupFiberHandler(mockService *services.MockTaskService) *fiber.App {
	app := fiber.New()
	handler := handlers.TaskHandler{Service: mockService}
	app.Post("/tasks", handler.CreateTask)
	app.Get("/tasks/:id", handler.GetTaskByID)
	app.Get("/tasks", handler.GetAllTasks)
	app.Put("/tasks/:id", handler.UpdateTask)
	app.Delete("/tasks/:id", handler.DeleteTask)
	return app
}

// Test_CreateTask_Handler_Success prüft, dass ein Task erfolgreich erstellt wird (Status 201) und die
// Response den Task-Titel enthält.
func Test_CreateTask_Handler_Success(t *testing.T) {
	app := setupFiberHandler(&services.MockTaskService{})

	reqBody := models.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "todo",
		Priority:    "medium",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err, "request should not fail")
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode, "expected status 201 Created")

	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Contains(t, string(data), "Test Task", "response should contain task title")
}

// Test_CreateTask_Handler_ValidationError prüft Validierungsfehler: leere Titel, zu lange Titel/Beschreibung,
// ungültige Priority/Status → Status 400.
func Test_CreateTask_Handler_ValidationError(t *testing.T) {
	app := setupFiberHandler(&services.MockTaskService{})

	testCases := []struct {
		Name string
		Body models.CreateTaskRequest
		Err  string
	}{
		{"Empty Title", models.CreateTaskRequest{Title: ""}, "Title is required"},
		{"Too Long Title", models.CreateTaskRequest{Title: strings.Repeat("a", 201)}, "Title is required"},
		{"Too Long Description", models.CreateTaskRequest{Title: "ok", Description: strings.Repeat("d", 1001)}, "Description must be max 1000 characters"},
		{"Invalid Priority", models.CreateTaskRequest{Title: "ok", Priority: "urgent"}, "Priority must be one of: low, medium, high or nothing"},
		{"Invalid Status", models.CreateTaskRequest{Title: "ok", Status: "waiting"}, "Status must be one of: todo, in progress, done or nothing"},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			body, _ := json.Marshal(tc.Body)
			req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err, "request should not fail")
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "expected 400 Bad Request")

			data, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			assert.Contains(t, string(data), "validation error", "response should contain validation error")
		})
	}
}

// Test_CreateTask_Handler_ServiceError prüft, dass ein Fehler im Service korrekt als Status 400 zurückgegeben wird.
func Test_CreateTask_Handler_ServiceError(t *testing.T) {
	app := setupFiberHandler(&services.MockTaskService{ShouldFail: true})

	reqBody := models.CreateTaskRequest{
		Title: "Test Task",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err, "request should not fail")
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "expected 400 Bad Request")
}

// Test_GetTasks_Handler_Success prüft, dass alle Tasks erfolgreich abgerufen werden (Status 200) und die
// Response alle Task-Titel enthält.
func Test_GetTasks_Handler_Success(t *testing.T) {
	mockService := &services.MockTaskService{
		Tasks: []*models.Task{
			{ID: 1, Title: "Test 1", Status: "todo", Priority: "low"},
			{ID: 2, Title: "Test 2", Status: "in progress", Priority: "high"},
		},
	}

	app := setupFiberHandler(mockService)

	req := httptest.NewRequest("GET", "/tasks", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err, "Request should not fail")
	assert.Equal(t, fiber.StatusOK, resp.StatusCode, "expected 200 OK")

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)

	assert.Contains(t, body, "Test 1", "response should contain first task")
	assert.Contains(t, body, "Test 2", "response should contain second task")
}

// Test_GetTaskByID_Handler_Found prüft, dass ein Task anhand der ID gefunden wird (Status 200) und Response
// den korrekten Task enthält.
func Test_GetTaskByID_Handler_Found(t *testing.T) {
	mockService := &services.MockTaskService{
		Tasks: []*models.Task{
			{ID: 1, Title: "API bauen", Status: "todo", Priority: "high"},
		},
	}

	app := setupFiberHandler(mockService)

	req := httptest.NewRequest("GET", "/tasks/1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err, "request should not fail")
	assert.Equal(t, fiber.StatusOK, resp.StatusCode, "expected 200 OK")

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)
	assert.Contains(t, body, "API bauen", "response should contain the task")
}

// Test_GetTaskByID_Handler_NotFound prüft, dass ein nicht vorhandener Task korrekt mit Status 404 beantwortet wird.
func Test_GetTaskByID_Handler_NotFound(t *testing.T) {
	mockService := &services.MockTaskService{
		Tasks: []*models.Task{},
	}

	app := setupFiberHandler(mockService)

	req := httptest.NewRequest("GET", "/tasks/999", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err, "request should not fail")
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode, "expected 404 Not Found")

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)
	assert.Contains(t, body, "not found", "response should contain 'not found' error")
}

// Test_UpdateTask_Handler_Success prüft, dass ein existierender Task erfolgreich aktualisiert wird (Status 200) und
// Response die neuen Werte enthält.
func Test_UpdateTask_Handler_Success(t *testing.T) {
	mockService := &services.MockTaskService{
		Tasks: []*models.Task{
			{ID: 1, Title: "API implementieren", Description: "API implementieren", Status: "todo", Priority: "medium", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}

	app := setupFiberHandler(mockService)

	reqBody := models.CreateTaskRequest{
		Title:       "API implementieren (aktualisiert)",
		Description: "API implementieren (aktualisiert)",
		Status:      "in progress",
		Priority:    "high",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var updatedTask models.Task
	err = json.Unmarshal(data, &updatedTask)
	assert.NoError(t, err)

	assert.Equal(t, 1, updatedTask.ID)
	assert.Equal(t, reqBody.Title, updatedTask.Title)
	assert.Equal(t, reqBody.Description, updatedTask.Description)
	assert.Equal(t, reqBody.Status, updatedTask.Status)
	assert.Equal(t, reqBody.Priority, updatedTask.Priority)
}

// Test_UpdateTask_Handler_NotFound prüft, dass ein Update auf einen nicht existierenden Task korrekt Status 404 liefert.
func Test_UpdateTask_Handler_NotFound(t *testing.T) {
	mockService := &services.MockTaskService{
		Tasks: []*models.Task{},
	}

	app := fiber.New()
	handler := handlers.TaskHandler{Service: mockService}
	app.Put("/tasks/:id", handler.UpdateTask)

	reqBody := models.CreateTaskRequest{
		Title: "Nicht existierende Task",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/tasks/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	assert.Contains(t, string(data), "not found")
}

// Test_DeleteTask_Handler_Success prüft, dass ein existierender Task erfolgreich gelöscht wird (Status 204 No Content).
func Test_DeleteTask_Handler_Success(t *testing.T) {
	mockService := &services.MockTaskService{
		Tasks: []*models.Task{
			{ID: 1, Title: "Test Task"},
		},
	}
	app := setupFiberHandler(mockService)

	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

// Test_DeleteTask_Handler_NotFound prüft, dass das Löschen eines nicht existierenden Tasks korrekt Status 404 zurückgibt.
func Test_DeleteTask_Handler_NotFound(t *testing.T) {
	mockService := &services.MockTaskService{
		Tasks: []*models.Task{},
	}

	app := setupFiberHandler(mockService)

	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "not found")
}
