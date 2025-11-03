# Kainos Core API

Microservice built with Go, Uber FX, and clean architecture principles.

## Architecture

### Dependency Injection (Uber FX)
- **Config Module**: Application configuration
- **Database Module**: PostgreSQL connection and store
- **NATS Module**: Message broker connection
- **Events Module**: Event publishing service
- **Temporal Module**: Workflow engine integration
- **Handlers Module**: HTTP request handlers
- **Server Module**: HTTP server and routing

### Event-Driven Architecture
- **User Events**: Published to NATS when users are created/updated/deleted via Clerk webhooks
- **Email Integration**: Events consumed by email service for notifications

### Folder Structure
```
core/
├── cmd/
│   └── main.go                    # Application entry point with FX
├── configs/
│   └── config.go                  # Configuration management
├── internal/
│   ├── database/
│   │   └── module.go              # Database connection module
│   ├── events/
│   │   └── publisher.go           # NATS event publisher
│   ├── fx/
│   │   └── modules.go             # FX dependency modules
│   ├── handlers/
│   │   ├── users/
│   │   │   └── handler.go         # Clerk webhook handler
│   │   └── workflow/
│   │       └── handler.go         # Workflow endpoints
│   ├── nats/
│   │   └── client.go              # NATS connection setup
│   ├── server/
│   │   └── server.go              # HTTP server with FX lifecycle
│   └── temporal/
│       └── module.go              # Temporal workflow engine
└── db/                            # Database schemas and queries
```

## Features

### Clerk Webhook Integration
- **User Created**: Publishes `user.created` event to NATS
- **User Updated**: Publishes `user.updated` event to NATS
- **User Deleted**: Publishes `user.deleted` event to NATS

### Event Publishing
Events are published to NATS with structure:
```json
{
  "id": "uuid",
  "type": "user.created|user.updated|user.deleted",
  "timestamp": "2024-01-01T00:00:00Z",
  "source": "core-api",
  "data": {
    "user_id": "clerk_user_id",
    "email": "user@example.com",
    "name": "John Doe",
    "first_name": "John",
    "last_name": "Doe"
  }
}
```

### Environment Variables
```bash
# Server
APP_SERVER_HOST=0.0.0.0
APP_SERVER_PORT=8081

# Database
APP_DATABASE_HOST=postgres
APP_DATABASE_PORT=5432
APP_DATABASE_NAME=kainos
APP_DATABASE_USERNAME=kainos
APP_DATABASE_PASSWORD=password

# NATS
APP_NATS_URL=nats://nats:4222

# Temporal
APP_TEMPORAL_HOSTPORT=temporal:7233
APP_TEMPORAL_NAMESPACE=default

# Clerk
CLERK_SECRET=your_clerk_secret

# JWT
APP_JWT_SECRET=your_jwt_secret
```

## Development

### Run Locally
```bash
go run cmd/main.go
```

### Build Docker Image
```bash
docker build -t core-service:latest .
```

### Test Webhook
```bash
curl -X POST http://localhost:8081/webhooks/clerk \
  -H "Content-Type: application/json" \
  -d '{
    "type": "user.created",
    "data": {
      "id": "user_123",
      "first_name": "John",
      "last_name": "Doe",
      "email_addresses": [
        {"email_address": "john@example.com"}
      ]
    }
  }'
```

## Dependencies

- **Uber FX**: Dependency injection framework
- **Gin**: HTTP web framework
- **NATS**: Message broker client
- **Temporal**: Workflow engine SDK
- **PostgreSQL**: Database driver (pgx)
- **Zerolog**: Structured logging
