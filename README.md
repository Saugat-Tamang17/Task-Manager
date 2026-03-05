# 📝 Task Manager API

A lightweight REST API built with **Go** and **PostgreSQL** for managing tasks. Supports full CRUD operations — create, read, update, and delete tasks via HTTP endpoints.

---

## 🛠 Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Driver:** github.com/lib/pq

---

## ⚙️ Prerequisites

Make sure you have the following installed:

- [Go](https://golang.org/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)

---

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/Saugat-Tamang17/Task-Manager.git
cd Task-Manager
```

### 2. Set up the database

Open PostgreSQL and create a database:

```sql
CREATE DATABASE user_api;
```

Then create the tasks table:

```sql
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200),
    description VARCHAR(500),
    status VARCHAR(15),
    created_at DATE
);
```

### 3. Set your database password as an environment variable

**Windows (PowerShell):**
```bash
$env:DB_PASSWORD="your_postgres_password"
```

**Mac/Linux:**
```bash
export DB_PASSWORD=your_postgres_password
```

### 4. Install dependencies

```bash
go mod tidy
```

### 5. Run the server

```bash
go run main.go
```

Server will start at `http://localhost:9090`

---

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tasks` | Fetch all tasks |
| POST | `/tasks` | Create a new task |
| PUT | `/tasks/{id}` | Update an existing task |
| DELETE | `/tasks/{id}` | Delete a task |

---

## 📦 Example Requests

### Create a Task (POST)
```json
{
    "title": "Learn Go",
    "description": "Complete the backend tutorial",
    "status": "pending",
    "created_at": "2026-03-04T00:00:00Z"
}
```

### Update a Task (PUT) — `/tasks/1`
```json
{
    "title": "Learn Go (Updated)",
    "description": "Completed the backend tutorial",
    "status": "completed",
    "created_at": "2026-03-04T00:00:00Z"
}
```

### Delete a Task (DELETE) — `/tasks/1`
No request body needed.

---

## 📁 Project Structure

```
Task-Manager/
├── main.go        # Main application entry point
├── go.mod         # Go module file
├── go.sum         # Dependency checksums
└── README.md      # Project documentation
```

---

## 👨‍💻 Author

**Saugat Tamang**  
[GitHub](https://github.com/Saugat-Tamang17)
