# Polaris Server

A Go-based API server following the standard Go project layout.

## Project Structure

```
.
├── cmd/                # Main applications for this project
│   └── api/           # API server entry point
├── internal/          # Private application and library code
│   ├── api/          # API handlers and routes
│   └── database/     # Database connection and operations
├── pkg/              # Library code that's ok to use by external applications
├── api/              # OpenAPI/Swagger specs, JSON schema files
├── configs/          # Configuration file templates or default configs
├── deployments/      # IaaS, PaaS, system and container orchestration deployment configs
├── docs/             # Design and user documents
└── scripts/          # Various build, install, analysis, etc scripts
```

## Prerequisites

- Go 1.21 or later
- MySQL 8.0 or later
- Make (optional, for using Makefile commands)

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
