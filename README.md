# VetSys

A veterinary management system built with Go and PostgreSQL.

## Features

- **Client Management** - Register and manage pet owners
- **Patient Management** - Track pets and their medical records
- **Consultation Management** - Schedule and record veterinary consultations
- **User Authentication** - Secure login with session management
- **Rate Limiting** - API protection with request rate limiting

## Tech Stack

- **Language**: Go 1.25.6
- **Database**: PostgreSQL
- **Router**: Go standard library 
- **Authentication**: Session-based with bcrypt password hashing

## Prerequisites

- Go 1.25.6 or higher
- PostgreSQL
- Make (optional)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd vetsys
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your environment variables:
```bash
cp .env.example .env
```

Edit `.env` with your configuration:
```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/vetsys?sslmode=disable
TEST_DATABASE_URL=postgres://postgres:postgres@localhost/vetsys_test?sslmode=disable
ENV=development
PORT=8888
```

4. Create the database:
```bash
createdb vetsys
createdb vetsys_test
```

## Running the Application

```bash
go run cmd/vetsys/main.go
```

The server will start on `http://localhost:8888`

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Project Structure

```
vetsys/
├── cmd/
│   └── vetsys/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database layer and repositories
│   ├── domain/          # Domain models
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # Authentication & rate limiting
│   ├── router/          # Route definitions
│   ├── server/          # Server setup
│   └── utils/           # Utility functions
└── go.mod
```