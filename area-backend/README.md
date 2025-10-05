# AREAmirror Backend Documentation

## Quick Overview
- REST API written in Go using `gin` and `gorm` for persistence.
- PostgreSQL database bootstrapped automatically with migrations and demo seed data.
- JWT based authentication with profile management, media uploads, and OAuth account linking (GitHub, Google, Facebook).
- Full CRUD for services, actions, reactions, areas, applets, and roles, guarded by authorization middleware.
- Integration services: background scheduler, Gmail email delivery, Discord webhooks, and GitHub webhooks (creation and processing).

## Core Technologies
- Go 1.23
- Gin Gonic (`github.com/gin-gonic/gin`)
- Gorm (`gorm.io/gorm` + `gorm.io/driver/postgres`)
- PostgreSQL
- OAuth libraries (Google, GitHub, Facebook) and the Gmail API

## Project Structure
```
area-backend/
|-- main.go                 # Entry point (routes, middleware, scheduler bootstrap)
|-- controllers/            # HTTP controllers grouped by domain
|-- database/               # PostgreSQL connection, enum creation, seed data
|-- models/                 # Gorm models (users, services, actions, reactions, areas, roles)
|-- services/               # External integrations (Gmail, Discord, GitHub, scheduler)
|-- uploads/                # Uploaded files (profile pictures)
`-- Dockerfile              # Container build configuration for the backend
```

## Startup Flow (`main.go`)
1. Load environment variables with `godotenv` (requires a `.env` file).
2. Initialize PostgreSQL via `database.InitDB()`; creates enums `area_status` and `run_status` when missing.
3. Run `AutoMigrate` on all models, then seed demo data through `database.SeedData` (roles, users, services, actions, reactions).
4. Configure the `gin` server with permissive CORS to support the frontend.
5. Register routes (public, JWT protected, and admin protected via `RoleMiddleware`).
6. Optionally start the scheduler in a goroutine (`services.NewSchedulerService`).

## Data Model and Persistence
- `models.User`: user account, profile fields, OAuth attributes, primary role, and `many2many` relationship with `Role`.
- `models.Service`, `models.Action`, `models.Reaction`: automation catalog; store JSON configuration blobs.
- `models.Area`: core IFTTT scenario storing trigger/action metadata, JSONB configs, scheduling properties, and run statistics.
- `models.Role` and `models.UserRole`: permission management (see `models/role.go` for default permissions).
- `models.Applet`: lightweight user-owned applets.

`database/seed.go` ensures:
- Roles `admin` and `member` exist.
- Demo accounts (`admin@area.com` / `admin123`, `john@example.com`, etc.).
- Sample services, actions, and reactions for quick testing.

## Authentication and Authorization
- `POST /register`, `POST /login`: create and authenticate accounts, returning a 24h JWT.
- `AuthMiddleware`: checks the `Authorization: Bearer <token>` header and injects `userID` / `userEmail` in the context.
- `RoleMiddleware("admin")`: limits access to administrative endpoints (areas, roles, role assignments, etc.).
- Profile management (`/profile`, `/profile/image`) and OAuth linking endpoints (`/profile/github/link`, `/profile/google/link`, `/profile/facebook/link`).
- OAuth flows rely on provider credentials defined in env vars and prevent linking accounts that are already associated with another user.

## Domain Controllers
### Users (`controllers/user.go`)
- Full CRUD (`/users`, `/users/:id`) with Bcrypt hashing on create/update.
- Inspect roles (`GET /users/:id/roles`), assign/remove roles (`POST/DELETE /users/:id/roles`).
- Update primary role (`PUT /users/:id/role`).

### Services, Actions, Reactions
- CRUD endpoints for `/services`, `/actions`, `/reactions`.
- Filter helpers `/service/:id/actions` and `/service/:id/reactions`.
- Creation validates that the related service exists.

### Areas (`controllers/area.go`)
- Protected by `AuthMiddleware` and `RoleMiddleware` for write operations.
- `GET /areas` (admin) and `GET /user/me/areas` return scenarios for admins or the current user.
- `POST /areas`: create scenarios from `CreateAreaRequest` (trigger/action configs serialized as JSON, icon lookup helpers).
- `PATCH /areas/:id/toggle`: enable or disable an area.
- `GET /areas/popular` and `/areas/recommended`: curated lists of active public areas.
- Test utilities: `/test/email`, `/test/discord`, `/test/scheduler/:id`.

### Applets (`controllers/applet.go`)
- Basic CRUD scoped to a user (`/user/:id/applets`).

## Scheduler (`services/scheduler.go`)
- Runs every 30 seconds to execute areas with `TriggerService == "Google Calendar"`.
- Validates the scheduled window (`eventTime` in `TriggerConfig`) and avoids re-running within five minutes.
- Supports two actions:
  - `Gmail`: sends email via `EmailService` using templates with variables `{{eventTitle}}`, `{{eventTime}}`, `{{areaName}}`.
  - `Discord`: sends webhook messages.
- Updates `LastRunAt`, `RunCount`, and `LastRunStatus` on success.
- `TestScheduler` forces an immediate execution for a given area.

## External Integrations
- **Email / Gmail (`services/email.go`)**: uses Google OAuth (client ID/secret + tokens) and Gmail SMTP (`GMAIL_USER` / `GMAIL_PASSWORD`). Provides HTML templates for GitHub notifications.
- **Discord (`services/discord.go`)**: posts messages to Discord webhooks with validation and error handling.
- **GitHub**:
  - OAuth account linking (`/profile/github/link`) and repository retrieval (`GET /api/github/repositories`).
  - Preconfigured GitHub->Gmail area creation via `POST /api/areas/github-gmail`.
  - `services/GitHubIntegrationService`: manages repository webhooks using a server token (`GITHUB_TOKEN`).
- **Google**: account linking, email sharing, and profile image upload via OAuth v2 (`GOOGLE_CLIENT_ID`, etc.).
- **Facebook**: account linking via the Facebook Graph API (`FACEBOOK_CLIENT_ID`, etc.).

## GitHub Webhooks (`controllers/github_webhook.go`)
- Endpoint `POST /webhooks/github`.
- Verifies HMAC SHA256 signatures when `WEBHOOK_SECRET` is defined.
- Handles `push`, `pull_request`, and `issues` events (push events have full processing).
- Delegates to `GitHubEventProcessor`, which fetches active GitHub areas (`trigger_service = 'github'`) and sends email notifications based on each area configuration.

## Environment Variables
```
# Database
DB_HOST
DB_PORT (default 5432)
DB_USER
DB_PASSWORD
DB_NAME
DB_SSLMODE (default disable)

# Auth
JWT_SECRET

# Gmail / Google
GOOGLE_CLIENT_ID
GOOGLE_CLIENT_SECRET
GOOGLE_ACCESS_TOKEN
GOOGLE_REFRESH_TOKEN
GOOGLE_REDIRECT_URI
GMAIL_USER
GMAIL_PASSWORD

# GitHub
GITHUB_CLIENT_ID
GITHUB_CLIENT_SECRET
GITHUB_TOKEN (used for server API calls and account linking)
WEBHOOK_URL (public URL for the backend webhook)
WEBHOOK_SECRET (HMAC signature for webhook payloads)

# Facebook OAuth
FACEBOOK_CLIENT_ID
FACEBOOK_CLIENT_SECRET
FACEBOOK_REDIRECT_URI
```

## Running the Backend Locally
1. Create a `.env` file (copy values from the Dockerfile block or define them manually).
2. Start PostgreSQL (for example with the repo-level docker-compose).
3. From `area-backend/`, install dependencies: `go mod download`.
4. Run the server: `go run .` (migrations and seed data execute automatically).
5. Optional: build a container `docker build -t areamirror-backend .` then `docker run --env-file .env -p 8080:8080 areamirror-backend`.

Seeded credentials allow immediate login (`admin@area.com` / `admin123`).

## Key API Endpoints
### Auth and Profile
- `POST /register`, `POST /login`
- `GET|PUT /profile`
- `POST /profile/image`
- `POST|DELETE /profile/github/link` / `/profile/github/unlink`
- `POST|DELETE /profile/google/link` / `/profile/google/unlink`
- `POST|DELETE /profile/facebook/link` / `/profile/facebook/unlink`

### Users and Roles (admin)
- `GET|POST /users`, `GET|PUT|DELETE /users/:id`
- `GET /users/:id/roles`, `POST|DELETE /users/:id/roles`
- `PUT /users/:id/role`
- `GET|POST /roles`, `GET|PUT|DELETE /roles/:id`

### Services / Actions / Reactions
- `GET|POST /services`, `GET|PUT|DELETE /services/:id`
- `GET /service/:id/actions`, `GET /service/:id/reactions`
- `GET|POST /actions`, `GET|PUT|DELETE /actions/:id`
- `GET|POST /reactions`, `GET|PUT|DELETE /reactions/:id`

### Areas
- `GET /areas` (JWT)
- `GET /areas/:id` (JWT)
- `POST /areas` (JWT + admin)
- `PUT /areas/:id` (JWT + admin)
- `DELETE /areas/:id` (JWT + admin)
- `PATCH /areas/:id/toggle` (JWT + admin)
- `GET /user/me/areas` (JWT)
- `GET /areas/popular`, `GET /areas/recommended`

### Applets
- `POST /user/:id/applets`
- `GET /user/:id/applets`
- `GET|PUT|DELETE /user/:id/applets/:id`

### Integrations and Tests
- `GET /api/github/repositories` (JWT)
- `POST /api/areas/github-gmail` (JWT)
- `POST /webhooks/github`
- `POST /test/email`, `POST /test/discord`, `POST /test/scheduler/:id`

## Operational Notes
- Email and Discord services require valid environment variables; without them the scheduler logs warnings but keeps running.
- `AuthMiddleware` expects a standard `Authorization` header. Ensure the frontend sends the `Bearer` prefix.
- GitHub webhooks need a public URL (ngrok, reverse proxy, or deployed environment) and a shared `WEBHOOK_SECRET`.
- OAuth flows must use redirect URIs registered with each provider.

Keep this documentation updated whenever enums (`createEnums()`), JSONB contract (`TriggerConfig`, `ActionConfig`), or external service requirements evolve.
