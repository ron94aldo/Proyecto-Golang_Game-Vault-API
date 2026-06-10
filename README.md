# Game Vault API

REST API built with Go, PostgreSQL, and the RAWG Video Games Database to manage a personal video game collection.

## Features

- 🔍 **Search Games**: Search for games in the RAWG database
- 📚 **Manage Library**: Add, update, and delete games from your personal collection
- 📊 **Statistics**: Track your gaming habits with collection statistics
- ✅ **Validation**: Comprehensive validation of input data
- 📝 **Unit Tests**: 80%+ code coverage with Go's native testing framework

## Tech Stack

- **Language**: Go 1.22+
- **HTTP**: Native `net/http` package
- **Database**: PostgreSQL
- **External API**: RAWG Video Games Database

## Installation

### Prerequisites

- Go 1.22 or higher
- PostgreSQL 12 or higher

### Setup

1. Clone the repository:

```bash
git clone https://github.com/ron94aldo/Proyecto-Golang_Game-Vault-API.git
cd Proyecto-Golang_Game-Vault-API
```

2. Set up PostgreSQL database:
```bash
createdb -h 127.0.0.1 -U postgres game_vault
psql -h 127.0.0.1 -U postgres -d game_vault -f db/schema.sql
```

3. Set environment variables (optional, defaults provided):
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=[You_Password]
export DB_NAME=game_vault
```

4. Download dependencies:
```bash
go mod download
```

5. Run the server:
```bash
go run main.go
```

The API will be available at http://localhost:8080

# API Endpoints

Search Endpoints

Search Games in RAWG:

```bash
GET /api/search?q=zelda
```

Get Game Details:

```bash
GET /api/games/3
```
Library Endpoints

List Library:

```bash
GET /api/library
GET /api/library?status=completado
```

Add Game to Library:

```bash
POST /api/library
Content-Type: application/json

{
  "rawg_id": 3,
  "title": "The Legend of Zelda: Breath of the Wild",
  "genre": "Adventure",
  "platform": "Nintendo Switch",
  "cover_url": "https://..."
}
```

Update Game:

```bash
PUT /api/library/1
Content-Type: application/json

{
  "personal_note": "Completed all dungeons!",
  "personal_score": 9,
  "status": "completado"
}
```

Valid status values: pendiente, jugando, completado, abandonado 
Valid personal score: 1-10

Delete Game:

```bash
DELETE /api/library/1
Response: 204 No Content
```

Get Statistics:

```bash
GET /api/library/stats
```

# Error Handling
All error responses follow this format:

```bash
{
  "code": "error_code",
  "error": "Error description"
}
```

HTTP Status Codes
Code	Meaning	Example
200	OK	Successful GET/PUT request
201	Created	Successful POST request
204	No Content	Successful DELETE request
400	Bad Request	Invalid body or missing fields
404	Not Found	Game not found in library
409	Conflict	Duplicate rawg_id
500	Internal Server Error	Database connection failure
502	Bad Gateway	RAWG API unreachable

Testing
Run all tests:

```bash
go test ./...
```

Run with coverage:

```bash
go test -coverprofile coverage.out ./...
go tool cover -html coverage.out
go tool cover -func coverage.out
```

# Environment Variables

DB_HOST	localhost	(PostgreSQL host)

DB_PORT	5432	(PostgreSQL port)

DB_USER	postgres	(Database user)

DB_PASSWORD	postgres	(Database password)

DB_NAME	game_vault	(Database name)

# API Key

The RAWG API key used in this project is publicly provided for educational purposes:

```bash
Key: 945d345a57cc4c3fb7b4f67211edd4c8
Base URL: https://api.rawg.io/api
Documentation: https://rawg.io/apidocs
```

# Author

Name: ron94aldo
Repository: Proyecto-Golang_Game-Vault-API

# License
This project is provided for educational purposes.

# Support
For issues or questions, please open an issue on the GitHub repository.
