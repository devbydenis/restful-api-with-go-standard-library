# Simple REST API with Go

This is a lightweight RESTful API for managing tasks, built entirely using the Go standard library (`net/http`). It demonstrates how to build a web server, handle HTTP routes, parse JSON, and manage in-memory state without third-party frameworks like Gin or Echo.

## Features

- **Standard Library Only**: Built using pure Go (`net/http`, `encoding/json`, etc.).
- **CRUD Operations**: Create, Read, and Delete tasks.
- **Filtering**:
  - Get tasks by specific **Tag**.
  - Get tasks by **Due Date** (Year/Month/Day).
- **Concurrency Safe**: Uses `sync.Mutex` to safely manage the in-memory data store across multiple requests.

## Project Structure

```
├── internal
│   └── taskstore   # In-memory "database" logic and data structures
├── model           # Structs for JSON Request and Response bodies
└── main.go         # Entry point, HTTP handlers, and server configuration
```

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) (version 1.22 or later recommended for `net/http` routing enhancements)

### Installation & Running

1. **Clone the repository:**

   ```bash
   git clone <repository-url>
   cd <project-directory>
   ```

2. **Run the server:**
   You must set the `SERVERPORT` environment variable before running.

   **Linux/macOS:**

   ```bash
   SERVERPORT=8080 go run main.go
   ```

   **Windows (PowerShell):**

   ```powershell
   $env:SERVERPORT="8080"; go run main.go
   ```

   The server will start at `http://localhost:8080` (or whichever port you specified).

## API Endpoints

### 1. Create a Task

**POST** `/task/`

**Request Body:**

```json
{
  "text": "Learn Go Standard Library",
  "tags": ["go", "learning", "backend"]
}
```

> Note: The `due` date is currently automatically set to the server's current time upon creation.

**Response:**

```json
{
  "status": 201,
  "message": "Task with 0 is created"
}
```

### 2. Get All Tasks

**GET** `/tasks/`

**Response:**

```json
{
  "status": 200,
  "message": "Fetched successfully",
  "data": [
    {
      "id": 0,
      "text": "Learn Go Standard Library",
      "tags": ["go", "learning", "backend"],
      "due": "2024-05-20T10:00:00Z"
    }
  ]
}
```

### 3. Get Task by ID

**GET** `/task/{id}/`

- Example: `/task/0/`

**Response:**

```json
{
  "id": 0,
  "text": "Learn Go Standard Library",
  "tags": ["go", "learning", "backend"],
  "due": "2024-05-20T10:00:00Z"
}
```

### 4. Get Tasks by Tag

**GET** `/tag/{tag}/`

- Example: `/tag/learning/`

**Response:**

```json
{
  "status": 200,
  "message": "Succesfully fetched data by tag",
  "data": [
    {
      "id": 0,
      "text": "Learn Go Standard Library",
      "tags": ["go", "learning", "backend"],
      "due": "2024-05-20T10:00:00Z"
    }
  ]
}
```

### 5. Get Tasks by Date

**GET** `/due/{year}/{month}/{day}/`

- Example: `/due/2024/5/20/`

**Response:**

```json
{
  "status": 200,
  "message": "Successfully fetch data by due",
  "data": [
    {
      "id": 0,
      "text": "Learn Go Standard Library",
      "tags": ["go", "learning", "backend"],
      "due": "2024-05-20T10:00:00Z"
    }
  ]
}
```

### 6. Delete Task by ID

**DELETE** `/task/{id}/`

- Example: `/task/0/`

**Response:**

```json
{
  "status": 204,
  "message": "Task with id 0 success to remove"
}
```

### 7. Delete All Tasks

**DELETE** `/tasks/`

**Response:**

```json
{
  "status": 204,
  "message": "Successfully remove all task!"
}
```

## Development Notes

- **Routing**: This project leverages the `http.NewServeMux` which has been enhanced in recent Go versions to support method-based routing (e.g., `"POST /task/"`) and path values (e.g., `"{id}"`).
- **Storage**: Data is stored in memory (`map[int]Task`). **All data will be lost when the server restarts.**
