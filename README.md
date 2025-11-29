# Task API

Eine kleine REST-API für Task-Management, geschrieben in **Go** mit **Fiber**, **PostgreSQL** und **unit-testbaren Services**.

---

## 🛠️ Setup mit Docker Compose

### Voraussetzungen

- Docker
- Docker Compose
- Git

### 1. Repository klonen

```bash
git clone <REPO_URL>
cd task-api
```

### 2. .env-Datei

Die `.env`-Datei liegt im Projekt und wird automatisch von Docker Compose genutzt. Beispiel:

```bash
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secret
POSTGRES_DB=tasks
```

Hinweis: Laut Aufgabenstellung wird die `.env` mit gepusht.

### 3. Docker Compose starten
```bash
docker compose up --build
```

Der API-Server läuft auf `http://localhost:8080`

Die PostgreSQL-Datenbank läuft im Container `db`

### 4. Container stoppen
```bash
docker compose down
```

## 📦 API Endpoints

Alle Endpoints erwarten/geben **JSON**.

### Health Check
```bash
GET /health
```

#### Antwort:

- `200 OK` → `"OK"`

### Tasks erstellen
```bash
POST /tasks
```

#### Request Body:
```bash
{
"title": "Einkaufen",
"description": "Lebensmittel besorgen",
"status": "todo",
"priority": "medium"
}
```
#### Antwort:

- `201 Created` → JSON des erstellten Tasks
- `400 Bad Request` → Validierungsfehler

### Alle Tasks abrufen
```bash
GET /tasks
```
#### Antwort:

- `200 OK` → Liste aller Tasks, inkl. `id`, `title`, `status`, `priority`, `created_at`
- `400 Bad Request` → DB Fehler

### Task nach ID abrufen
```bash
GET /tasks/:id
```
#### Antwort:

- `200 OK` → Task als JSON
- `404 Not Found` → Task existiert nicht

### Task aktualisieren
```bash
PUT /tasks/:id
```


#### Request Body: *(nur zu ändernde Felder angeben)*
```bash
{
"title": "Neue Beschreibung",
"status": "in progress"
}
```

#### Antwort:

- `200 OK` → aktualisierter Task
- `400 Bad Request` → Validierungsfehler
- `404 Not Found` → Task existiert nicht

### Task löschen
```bash
DELETE /tasks/:id
```

#### Antwort:

- `204 No Content` → erfolgreich gelöscht
- `404 Not Found` → Task existiert nicht

nicht

## 🧾 Datenmodelle

### Task

| Feld        | Typ        | Beschreibung                       |
|------------|------------|------------------------------------|
| id         | int        | Eindeutige ID                      |
| title      | string     | Pflichtfeld, max 200 Zeichen       |
| description| string     | Optional, max 1000 Zeichen         |
| status     | string     | "todo", "in progress", "done"      |
| priority   | string     | "low", "medium", "high"            |
| created_at | time.Time  | Zeitpunkt der Erstellung           |
| updated_at | time.Time  | Zeitpunkt der letzten Änderung     |

Das `Task`-Modell repräsentiert einen einzelnen Task innerhalb der API. Es definiert 
alle Eigenschaften eines Tasks, die in der Datenbank gespeichert und über die API 
zurückgegeben werden. Dazu gehören Titel, Beschreibung, Status, Priorität sowie 
Zeitstempel für Erstellung und letzte Aktualisierung. Dieses Modell dient als 
zentrale Datenstruktur für alle CRUD-Operationen in der Anwendung.

### CreateTaskRequest

| Feld        | Typ    | Beschreibung                        |
|------------|--------|------------------------------------|
| title      | string | Pflichtfeld, max 200 Zeichen       |
| description| string | Optional, max 1000 Zeichen         |
| status     | string | "todo", "in progress", "done" (optional, default "todo") |
| priority   | string | "low", "medium", "high" (optional, default "medium")     |

Das `CreateTaskRequest`-Modell definiert die Datenstruktur, die benötigt wird, um einen 
neuen Task über die API zu erstellen. Es legt fest, welche Felder optional oder 
verpflichtend sind und welche Standardwerte verwendet werden, falls bestimmte Angaben 
fehlen. Damit wird sichergestellt, dass neue Tasks konsistent und fehlerfrei angelegt 
werden können.

## 🧪 Tests
### Handler-Tests

- Liegen in `handlers/handlers_test.go`
- Testen die HTTP-Endpunkte mit Mock Services
- Beispiel-Tests:
    - `Test_CreateTask_Handler_Fiber`
    - `Test_GetTaskByID_Handler_Found`
    - `Test_DeleteTask_Handler_NotFound`

### Service-Tests

- Liegen in `services/services_test.go`
- Testen die Businesslogik direkt mit Mock-Repositories
- Beispiel-Tests:
  - `Test_Service_CreateTask_Success_with_priority_and_status`
  - `Test_Service_GetTaskByID_NotFound`
  - `Test_Service_DeleteTask_Success`
```bash
go test ./... 
```

## 🔧 Mocks

- `MockTaskService` → simuliert `TaskServiceInterface`
- `MockTaskRepository` → simuliert `TaskRepositoryInterface`

Dient zum Testen der Services und Handler ohne echte Datenbank.

.

## 📌 Hinweise

- Validierte Felder:
    - `Title` (Pflicht, max 200 Zeichen)
    - `Description` (optional, max 1000 Zeichen)
    - `Status`: `"todo" | "in progress" | "done"`
    - `Priority`: `"low" | "medium" | "high"`

- Standardwerte beim Erstellen, falls leer:
  - Status: `"todo"`
  - Priority: `"medium"`