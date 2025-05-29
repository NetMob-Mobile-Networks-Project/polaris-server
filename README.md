# Polaris Server

A Go server application using MySQL with sqlc for type-safe database operations.

## Project Structure

```
.
├── cmd/                    # Main applications for this project
│   └── polaris-server/    # Server entry point and configuration
├── internal/              # Private application and library code
│   ├── api/              # API handlers and routes
│   └── database/         # Database related code
│       ├── query/        # SQL queries for sqlc
│       │   └── users.sql
│       ├── schema/       # Database schema
│       │   └── schema.sql
│       ├── mysql.go      # Database connection and queries wrapper
│       └── sqlc.yaml     # sqlc configuration
├── pkg/                  # Library code that's ok to use by external applications
├── api/                  # OpenAPI/Swagger specs, JSON schema files
├── configs/              # Configuration file templates or default configs
├── deployments/          # IaaS, PaaS, system and container orchestration deployment configs
├── docs/                 # Design and user documents
└── scripts/              # Various build, install, analysis, etc scripts
```

## Configuration

The server uses environment variables for configuration. The following environment variables are supported:

```bash
# Server Configuration
PORT=8000

# Database Configuration
DB_USER=root
DB_PASS=secret
DB_HOST=localhost
DB_PORT=3306
DB_NAME=polaris

# API Configuration
API_VERSION=v1
API_PREFIX=/api

# JWT Configuration
JWT_SECRET=your-secret-key-here
```

You can set these variables in your environment or use a `.env` file.

## Database Setup

The project uses [sqlc](https://sqlc.dev/) for type-safe database operations. The database schema and queries are defined in the `internal/database` directory.

### Prerequisites

- Go 1.23 or later
- MySQL
- sqlc CLI tool

### Installation

1. Install sqlc:
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

2. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Generate sqlc code:
```bash
cd internal/database
sqlc generate
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
   go run cmd/polaris-server/main.go
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
go build -o bin/polaris-server cmd/polaris-server/main.go
```

## License

MIT
