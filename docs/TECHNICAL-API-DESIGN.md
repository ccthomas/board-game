# Board Game - Technical API Design

## Overview

This service exposes operational and infrastructure-level APIs for:

* Service health monitoring
* PostgreSQL database version reporting
* Database migration management

The API is built using:

* Router: Gorilla Mux
* Database: PostgreSQL
* Migrations: golang-migrate

All endpoints are mounted under:

```
/api
```

---

# Architecture Overview

The API follows a layered architecture:

```
Controller → Service → Database → PostgreSQL
```

### Controller Layer

* HTTP request parsing
* Response formatting
* Status codes
* Logging context enrichment

### Service Layer

* Business logic
* Command validation
* Migration orchestration
* Database interaction abstraction

### Database Layer

* Connection management
* Migration execution
* Version retrieval

---

# Base URL

```
/api
```

---

# 1️⃣ Health API

## Endpoint

```
GET /api/health
```

## Purpose

Returns:

* Service status
* PostgreSQL version
* Current timestamp

## Controller

`HealthController.GetHealth`

## Service

`HealthService.GetDatabaseVersion`

---

## Response

### 200 OK

```json
{
  "service": "Healthy",
  "database": "PostgreSQL 16.1 on x86_64...",
  "timestamp": "2026-02-28T04:11:22.384Z"
}
```

### Behavior Details

| Field     | Source               |
| --------- | -------------------- |
| service   | Hardcoded "Healthy"  |
| database  | `database.Version()` |
| timestamp | `time.Now()`         |

If database version retrieval fails:

* The error string is returned as the `database` value
* HTTP status remains `200`

---

## Failure Scenarios

| Scenario          | Result                    |
| ----------------- | ------------------------- |
| JSON encode fails | 500 Internal Server Error |

---

# 2️⃣ Database Migration API

## Endpoint

```
POST /api/configuration/database/migrations
```

## Purpose

Allows controlled execution of database migrations via API.

---

## Request Body

```json
{
  "command": "up",
  "quantity": 2
}
```

### Fields

| Field    | Type   | Required | Description               |
| -------- | ------ | -------- | ------------------------- |
| command  | string | yes      | `"up"` or `"down"`        |
| quantity | int8   | no       | Number of migration steps |

---

## MigrationCommand Enum

```go
const (
    MigrationUp   = "up"
    MigrationDown = "down"
)
```

---

# Migration Execution Logic

### command = "up"

| quantity | Behavior                         |
| -------- | -------------------------------- |
| null     | Runs full `MigrationDown()`      |
| number   | Runs `MigrationSteps(+quantity)` |

### command = "down"

| quantity | Behavior                         |
| -------- | -------------------------------- |
| null     | Runs full `MigrationUp()`        |
| number   | Runs `MigrationSteps(-quantity)` |

> ⚠ Note: Current implementation appears logically inverted:
>
> * `"up"` with nil runs `MigrationDown`
> * `"down"` with nil runs `MigrationUp`
>
> This may be intentional or a defect.

---

## Success Response

```
200 OK
```

Empty body.

---

## Error Responses

### 400 Bad Request

Invalid JSON body.

```
invalid body
```

### 500 Internal Server Error

Returned when:

* Database connection fails
* Migration fails
* Unknown command received

---

# Database Layer Details

## Connection

Connection parameters are retrieved from environment variables:

```
DB_USER
DB_PASSWORD
DB_HOST
DB_NAME
DB_PORT
```

Connection string:

```
user=... password=... dbname=... host=... port=... sslmode=disable TimeZone=America/New_York
```

---

## Migration Engine

Migrations are loaded from:

```
db/migrations
```

Using:

```
file://<absolute_path>
```

via `migrate.NewWithDatabaseInstance`.

---

# Logging Strategy

Each layer enriches logs with:

* file_name
* class_name
* contextual fields
* command values
* migration step counts

Log Levels Used:

| Level | Usage                              |
| ----- | ---------------------------------- |
| Trace | Detailed execution flow            |
| Debug | Entry/exit points                  |
| Info  | Migration file discovery           |
| Warn  | Invalid commands / decode failures |
| Error | Runtime failures                   |

---

# Security Considerations

⚠ This API allows database schema mutation via HTTP.

Recommended protections:

* Restrict to internal network
* Require authentication (e.g. admin token)
* Disable in production builds
* Add IP allow-listing
* Add environment guard (only allow in non-prod)

---

# Operational Considerations

## Health Endpoint

Safe for:

* Kubernetes liveness probe
* Readiness probe
* Monitoring systems

## Migration Endpoint

Intended for:

* CI/CD automation
* Controlled administrative usage

Not intended for:

* Public access
* End-user exposure

---

# Example curl Commands

### Health

```bash
curl http://localhost:8080/api/health
```

### Run Full Migration Up

```bash
curl -X POST http://localhost:8080/api/configuration/database/migrations \
  -H "Content-Type: application/json" \
  -d '{"command":"up"}'
```

### Run 2 Steps Down

```bash
curl -X POST http://localhost:8080/api/configuration/database/migrations \
  -H "Content-Type: application/json" \
  -d '{"command":"down","quantity":2}'
```

---

# Future Improvements

* Fix migration command inversion logic
* Return structured success response
* Add authentication middleware
* Add OpenAPI / Swagger specification
* Add structured error responses
* Add correlation/request IDs
* Separate admin APIs from public APIs

---

# Summary

This API layer provides:

* Operational visibility (`/health`)
* Infrastructure control (`/configuration/database/migrations`)
* Layered architecture with clean separation
* Structured logging at every layer
* PostgreSQL-backed migration management
