# NotesGo

Notes API built with **Go + Gin + MySQL**, organized in **hexagonal architecture**
(adapters → core ports → services), with JWT authentication and RBAC
authorization (`HasAccess`) enforced per route.

## Prerequisites

| Tool | Version | Install |
|---|---|---|
| Go | `1.26.0` (see `go.mod`) | https://go.dev/dl |
| MySQL | 8.x | Homebrew: `brew install mysql` + `brew services start mysql` |
| goose CLI | `v3.29.0` (see `go.mod`) | `go install github.com/pressly/goose/v3/cmd/goose@latest` |
| swag CLI | `v1.16.x` (see `go.mod`) | `go install github.com/swaggo/swag/cmd/swag@latest` |
| sqlc | only needed to regenerate `db/schemas` | `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` (or check `~/go/bin/sqlc`) |
| air | optional, live-reload | `go install github.com/air-verse/air@latest` |
| make | macOS Xcode CLT | `xcode-select --install` |

> The `Makefile` prepends `GOPATH/bin` to `PATH` by itself, so `make migrate-*`
> works even if `~/go/bin` is not on your shell `PATH`. All `make` targets also
> work as raw `goose` commands (see [Migrations](#migrations--seeds)).

## Quick start

```sh
# 1. Clone and enter the project
git clone <repo-url> notes_app && cd notes_app

# 2. Configure environment
cp .env.example .env
# edit .env: DBUSER, DBPASS, DBADDRESS (e.g. 127.0.0.1:3306), DBNAME, JWT_SECRET

# 3. Install Go dependencies
go mod download

# 4. Create the database (once)
mysql -h 127.0.0.1 -P 3306 -u <user> -p -e "CREATE DATABASE IF NOT EXISTS notesdb;"

# 5. Run migrations + seeds (tables + roles/actions/permissions)
make migrate-up

# 6. Run the API (auto-migrates on boot unless RUN_MIGRATIONS=false)
go run .
# or with live-reload:
air

# 7. Smoke test
curl localhost:8080/api/v1/health
# {"message":"Everything's good"}
```

## Configuration

Copy `.env.example` to `.env`. Every variable:

| Variable | Example | Effect |
|---|---|---|
| `DBUSER` | `root` | MySQL user |
| `DBPASS` | `"secret#2025"` | MySQL password — **quote it** if it contains `#`, `$`, spaces |
| `DBNET` | `tcp` | Network for `go-sql-driver/mysql` |
| `DBADDRESS` | `127.0.0.1:3306` | Host:port of MySQL |
| `DBNAME` | `notesdb` | Database (must exist before `migrate-up`) |
| `JWT_SECRET` | `change-me-in-production` | HMAC secret for login tokens (2h expiry) |
| `RUN_MIGRATIONS` | `true` | `false` disables auto-migrate on boot (`main.go`, `db/migrate.go`) |
| `MIGRATION_DIR` | `db/migrations` | Override migration directory |

> Quoting note: values with `#` **must** stay quoted in `.env`. The shell,
> `godotenv`, and the `Makefile` recipes (which source `.env` via shell) handle
> this correctly. Do **not** read `.env` with `make include` — `make` treats `#`
> as a comment starter even inside quotes.

## Migrations & seeds

Migrations live in `db/migrations` and run with **goose** (MySQL dialect).

| File | Content |
|---|---|
| `00001_init.sql` | Baseline schema (`roles`, `actions`, `permissions`, `users`, `notes`), `IF NOT EXISTS`, safe on existing DBs |
| `00002_seed_rbac.sql` | Seed data required **before any signup**: 3 roles, 9 actions, 9 permission rows (all `INSERT IGNORE`, pinned IDs) |
| `00003_seed_admin_user.sql` | Default admin login `admin@gmail.com` / `admin` (dev only — weak password, change or delete before sharing the DB) |

Seeded matrix:

| Role (id) | Actions |
|---|---|
| Administrator (1) | `Admin:Add`, `Admin:Read`, `Admin:Update`, `Admin:Delete` |
| Manager (2) | `Manager:Add`, `Manager:Read`, `Manager:Update` |
| User (3) | `User:Add`, `User:Read` |

`make` targets (`.env` is sourced automatically):

| Command | Effect |
|---|---|
| `make migrate-up` | Apply pending migrations (`... up`) |
| `make migrate-down` | Roll back one migration (`... down 1`) |
| `make migrate-status` | Show applied/pending migrations |
| `make migrate-create` | Prompts for a name, scaffolds a new `db/migrations/*.sql` |

Raw equivalents (no `make` needed):

```sh
set -a; source .env; set +a
goose -dir db/migrations mysql "$DBUSER:$DBPASS@tcp($DBADDRESS)/$DBNAME?parseTime=true" up
# replace `up` with `status`, `down 1`, or `create <name> sql`
```

The app also migrates on boot (`db.MigrateUp` in `main.go:25`, skip with
`RUN_MIGRATIONS=false`). Seeds are idempotent — running `up` twice is a no-op.

### Migration rules (read before touching `db/migrations`)

1. **Never edit an already-applied migration.** Goose tracks applied versions in
   `goose_db_version` and never re-runs them, so edits are silently ignored
   everywhere the migration already ran. Add a **new** migration instead
   (`make migrate-create`), e.g. `00003_seed_admin_user.sql` was created this
   way instead of editing `00002`.
2. **Seeds must be idempotent.** Use `INSERT IGNORE` with explicit pinned IDs
   (roles `1-3`, actions `1-9`) so re-running and fresh installs converge to the
   same rows — permission pairs like `(1,1)` only stay valid if the IDs match.
3. **Seed password hashes, never plaintext.** Login compares with
   `bcrypt.CompareHashAndPassword` (`utils/hash.go`, cost 12), so a literal
   string in `password_hash` can never log in. Generate the value with the
   app's own hasher (e.g. a throwaway `go run` calling `utils.HashPassword`)
   and paste the `$2a$...` result into the migration.
4. **Every `Up` needs a `Down`.** The down migration must undo exactly what the
   up did (e.g. `DELETE FROM users WHERE id = 1`), in reverse FK order for
   schema changes — otherwise `migrate-down` leaves orphan rows or fails.

## Running the app

```sh
go run .                    # dev, serves :8080 (see main.go:37)
air                         # dev with live-reload (config .air.toml, binary ./tmp/main)
go build -o ./tmp/main .    # production-style binary
./tmp/main
```

Public endpoint (no token):

```sh
curl localhost:8080/api/v1/health
```

Interactive API docs (Swagger UI, no token needed to browse):

```sh
open http://localhost:8080/swagger/index.html
```

## API documentation (Swagger)

Annotations live next to each handler (`@Summary`, `@Tags`, `@Param`,
`@Success`, `@Router`, `@Security`); generated output is committed under
`docs/` (`docs.go`, `swagger.json`, `swagger.yaml`).

```sh
go install github.com/swaggo/swag/cmd/swag@latest  # once, lands in ~/go/bin
make swagger      # regenerates docs/ — run after any handler/annotation change
```

Flow: `POST /auth/login` → copy `data.token` → Swagger UI **Authorize** button
(API key `Authorization`, **raw token, no `Bearer` prefix**) → try any endpoint.
All responses use the `ResponseStruct` envelope
`{message, success, error, data}` — each operation's description names the
concrete DTO inside `data`.

## Auth & RBAC

**Signup** — creates a `User` with the default `User` role (`role_id = 3`):

```sh
curl -X POST localhost:8080/api/v1/auth/signup \
  -H 'Content-Type: application/json' \
  -d '{"name":"Ada","email":"ada@example.com","password":"secret123","phoneNumber":"1234567"}'
```

**Login** — returns a JWT (`{"data":{"token":"..."}}`):

```sh
curl -X POST localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"ada@example.com","password":"secret123"}'
```

**Authenticated requests** — send the raw token (no `Bearer` prefix) in the
`Authorization` header (`middlewares/auth.go`):

```sh
TOKEN=<token-from-login>
curl localhost:8080/api/v1/notes/ -H "Authorization: $TOKEN"
```

**Authorization model** — `middlewares.RBAC(userService, permissions...)` calls
`UserService.HasAccess`, which counts rows joining
`users → roles → permissions → actions` and requires the user's role to own at
least one of the listed `<Role>:<Action>` strings. Missing/invalid token → `401`,
valid token without permission → `403`.

### Endpoints × required permission

Base path `/api/v1`. All except `health`/`auth` require `Authorization` header.

| Method & path | Permission (any listed) |
|---|---|
| `GET /health` | public |
| `POST /auth/signup`, `POST /auth/login` | public |
| `GET /actions/`, `GET /actions/:id` | `Admin:Read` |
| `GET /roles/`, `GET /roles/:id` | `Admin:Read` |
| `GET /permissions/` | `Admin:Read` |
| `GET /users/`, `GET /users/:id` | `Admin:Read`, `Manager:Read` |
| `PUT /users/:id` (`name,email,password,phoneNumber`) | `Admin:Update`, `Manager:Update` |
| `POST /users/change-role` (`{"userId":1,"roleId":2}`) | `Admin:Update` |
| `DELETE /users/:id` | `Admin:Delete` |
| `GET /notes/`, `GET /notes/:id` | `Admin:Read`, `Manager:Read` |
| `GET /notes/user` | `Admin:Read`, `Manager:Read`, `User:Read` |
| `POST /notes/` | `Admin:Add`, `Manager:Add`, `User:Add` |
| `PUT /notes/:id` | `Admin:Update`, `Manager:Update` |
| `DELETE /notes/:id` | `Admin:Delete` |

## Project structure

```
.
├── main.go                          # gin setup, auto-migrate, DI wiring, :8080
├── Makefile                         # migrate-up/down/status/create (goose)
├── sqlc.yaml                        # sqlc: queries db/sqlc/query.sql, schema = migration 00001
├── .air.toml                        # live-reload config (kills stale binary, cleans on exit)
├── db/
│   ├── conn.go                      # *sql.DB from DBUSER/DBPASS/DBNET/DBADDRESS/DBNAME
│   ├── migrate.go                   # DSN builder + goose.Up wrapper (RUN_MIGRATIONS gate)
│   ├── migrations/00001_init.sql    # baseline schema
│   ├── migrations/00002_seed_rbac.sql # roles/actions/permissions seed
│   ├── sqlc/{schema.sql,query.sql}  # sqlc inputs (schema.sql kept as reference)
│   └── schemas/                     # sqlc-generated models + queries (DO NOT EDIT)
├── internal/
│   ├── adapters/
│   │   ├── https/handlers/v1/       # Gin handlers + Register*Routes (route + RBAC wiring)
│   │   ├── https/routes/routes.go   # /api/v1 groups; Authenticate on authed group
│   │   ├── https/dtos/              # request/response shapes + JSON mapping
│   │   └── repository/              # output-port implementations (sqlc-backed)
│   └── core/
│       ├── ports/input/             # service interfaces (driven by handlers/middlewares)
│       ├── ports/output/            # repository interfaces (driven by services)
│       ├── service/                 # business rules (validation, hashing, HasAccess, ChangeRoles)
│       ├── domain/                  # domain types (e.g. Login)
│       └── errors/                  # shared errors (ErrIdNotProvided, ...)
├── middlewares/
│   ├── auth.go                      # Authenticate: verifies JWT → sets userId/userEmail
│   └── rbac.go                      # RBAC(service, perms...): 401/403 or Next()
└── utils/                           # JWT generate/verify, password hashing
```

Conventions: handlers only parse HTTP, map DTOs, and translate service errors to
status codes; existence/permission invariants live in `service` (e.g. `Update`
checks `GetById` first, like `ChangeRoles` checks user + role). Services depend
only on `ports/output` interfaces — never on concrete repositories or each other.

## Command reference

| Command | Purpose |
|---|---|
| `go mod download` / `go mod tidy` | Install / reconcile dependencies |
| `go build ./...` | Compile everything |
| `go vet ./...` | Static checks (run before committing) |
| `go run .` | Run the API locally |
| `air` | Run with live-reload |
| `make migrate-up` / `-down` / `-status` / `-create` | goose migrations (see above) |
| `make swagger` | Regenerate Swagger `docs/` from handler annotations |
| `sqlc generate` | Regenerate `db/schemas` after changing queries/migrations |
| `go install github.com/pressly/goose/v3/cmd/goose@latest` | Install goose CLI (lands in `~/go/bin`) |
| `lsof -i :8080 -sTCP:LISTEN -Pn` | Find what holds port 8080 |
| `pkill -f "notes_app/tmp/main"` | Kill stale dev binary holding the port |

## Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| `listen tcp :8080: bind: address already in use` | Stale `./tmp/main` from a previous `air`/`go run` still listening | `lsof -i :8080 ...` then `pkill -f "notes_app/tmp/main"`; current `.air.toml` (`send_interrupt`, `kill_delay: 1000ms`, `clean_on_exit`) prevents recurrence |
| `/bin/sh: goose: command not found` | `~/go/bin` not on shell `PATH` | Use `make migrate-*` (Makefile prepends `GOPATH/bin` itself) or `export PATH="$PATH:$(go env GOPATH)/bin"` |
| `Access denied for user ... (using password: YES)` via make | Password mangled (e.g. `#` treated as comment) | Keep `DBPASS` quoted in `.env`; current `Makefile` sources `.env` via shell so quoted values survive |
| `Error 1452 ... fk_roles_in_users` on change-role | `role_id` value doesn't exist in `roles` | Seed first (`make migrate-up` provides ids 1-3); service validates both ids before writing |
| `Authenticate must run before RBAC` | `RBAC` registered on the public group instead of the authed group | Register RBAC routes under `serverAuth` (which runs `Authenticate` first) — see `routes.go` |
