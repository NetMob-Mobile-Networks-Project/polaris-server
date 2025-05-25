# Polaris Server

A Go server application using MySQL with sqlc for type-safe database operations.

## Project Structure

```
.
├── cmd/                # Main applications for this project
│   └── api/           # API server entry point
├── internal/          # Private application and library code
│   ├── api/          # API handlers and routes
│   └── database/     # Database related code
│       ├── query/    # SQL queries for sqlc
│       │   └── users.sql
│       ├── schema/   # Database schema
│       │   └── schema.sql
│       ├── mysql.go  # Database connection and queries wrapper
│       └── sqlc.yaml # sqlc configuration
├── pkg/              # Library code that's ok to use by external applications
├── api/              # OpenAPI/Swagger specs, JSON schema files
├── configs/          # Configuration file templates or default configs
├── deployments/      # IaaS, PaaS, system and container orchestration deployment configs
├── docs/             # Design and user documents
└── scripts/          # Various build, install, analysis, etc scripts
```

## Database Setup

The project uses [sqlc](https://sqlc.dev/) for type-safe database operations. The database schema and queries are defined in the `internal/database` directory.

### Prerequisites

- Go 1.16 or later
- MySQL
- sqlc CLI tool

### Installation

1. Install sqlc:
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

2. Set up environment variables:
```bash
export DB_USER=your_username
export DB_PASS=your_password
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=your_database
```

3. Generate sqlc code:
```bash
cd internal/database
sqlc generate
```

### Database Schema

The current schema includes a `users` table with the following structure:

```sql
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Available Queries

The following queries are available through sqlc:

- `GetUser`: Retrieve a single user by ID
- `CreateUser`: Create a new user
- `ListUsers`: List users with pagination

## Usage

To use the database in your code:

```go
import "github.com/themhh/polaris-server/internal/database"

// Create a new database connection
db, err := database.NewMySQLConnection()
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Use the generated queries
user, err := db.GetUser(ctx, 1)
if err != nil {
    log.Fatal(err)
}
```

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/themhh/polaris-server.git
   cd polaris-server
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Update the `.env` file with your database credentials and other configurations.

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Run the server:
   ```bash
   go run cmd/api/main.go
   ```

## API Endpoints

### Health Check
- `GET /health` - Check server health status

### Users API (v1)
- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/{id}` - Get a specific user
- `PUT /api/v1/users/{id}` - Update a user
- `DELETE /api/v1/users/{id}` - Delete a user

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o bin/polaris-server cmd/api/main.go
```

## License

MIT
