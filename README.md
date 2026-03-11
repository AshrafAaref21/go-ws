# Go WebSockets Webchat — Backend

A real-time chat backend built with Go, using WebSockets for live messaging and SQLite for persistence. Supports private conversations, file sharing, JWT-based authentication, and graceful shutdown.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25+ |
| WebSockets | `coder/websocket` |
| Database | SQLite (`modernc.org/sqlite`) |
| Auth | JWT (`golang-jwt/jwt/v5`) |
| Config | `cleanenv` + `.env` files |
| Password | `bcrypt` (`golang.org/x/crypto`) |

---

## Project Structure

```
cmd/api/          → Entry point (main.go)
internal/
  config/         → Config loading from .env
  db/             → SQLite init and connection
  dto/            → Request/response data transfer objects
  handlers/       → HTTP & WebSocket handler functions
  middlewares/    → Auth, CORS, logging middleware
  models/         → Database models (User, Message, Private)
  realtime/       → WebSocket hub and client management
  routes/         → Route registration
  utils/          → JWT, password hashing, API responses
config/
  dev.env         → Environment configuration
sqlite/dev/       → SQLite database files (gitignored)
```

---

## API Endpoints

### Auth
| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/api/auth/register-email` | No | Register with email & password |
| `POST` | `/api/auth/login-email` | No | Login, returns access + refresh tokens |
| `POST` | `/api/auth/refresh-session` | No | Refresh access token |
| `POST` | `/api/auth/logout` | Yes | Invalidate session |
| `POST` | `/api/auth/current-user` | Yes | Get authenticated user info |

### Users
| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/api/users/{id}` | Yes | Get user by ID |

### Conversations
| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/api/conversations` | Yes | List all conversations |
| `POST` | `/api/conversations/privates/create` | Yes | Create or join a private conversation |
| `GET` | `/api/conversations/privates/{private_id}` | Yes | Get private conversation details |
| `GET` | `/api/conversations/privates/{private_id}/messages` | Yes | Get messages in a conversation |

### Files
| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/api/files/{private_id}` | Yes | Upload a file to a conversation |
| `GET` | `/api/files/` | Yes | Retrieve a file |

### WebSocket
| Endpoint | Description |
|---|---|
| `/api/ws` | WebSocket connection (authenticated via `Authorization` header) |

### Health
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/health-check-http` | HTTP health check |
| `GET` | `/api/health-check-ws` | WebSocket health check |

---

## Configuration

Configuration is loaded from a `.env` file. Default path: `config/dev.env`.

Override with a CLI flag or environment variable:
```sh
# CLI flag
./api -config path/to/custom.env

# Environment variable
CONFIG_PATH=path/to/custom.env ./api
```

**Available variables:**

| Variable | Default | Description |
|---|---|---|
| `ENV` | `dev` | Runtime environment |
| `HTTP_ADDRESS` | `localhost:8082` | Server listen address |
| `DB_PATH` | `sqlite/dev` | Directory for the SQLite database |
| `DB_NAME` | `api.db` | SQLite database filename |
| `JWT_KEY` | `supersecretkey` | Secret key for signing JWTs — **change in production** |

---

## Running the Server

```sh
# Build and run
go run ./cmd/api

# Or build binary first
go build -o bin/api ./cmd/api
./bin/api -config config/dev.env
```

The server supports graceful shutdown on `SIGINT` / `SIGTERM` with a 10-second timeout.

---

## Authentication Flow

1. Register or login to receive an **access token** (short-lived) and a **refresh token**.
2. Pass the access token in the `Authorization: Bearer <token>` header on protected routes.
3. Use `/api/auth/refresh-session` with the refresh token to obtain a new access token.
4. WebSocket connections authenticate via the same `Authorization` header on the upgrade request.

---

## WebSocket Events

The hub manages online clients per user ID and supports broadcasting events across connections. Events are JSON-encoded and dispatched to target user IDs on actions like new messages or file uploads.
