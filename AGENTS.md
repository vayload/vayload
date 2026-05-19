# Vayload Agent Guidelines

**Vayload** — A modular, extensible CMS built in Go with a Svelte 5 frontend.

Lightweight kernel, plugin system in embedded Lua scripting, multi-database support (MySQL, PostgreSQL, SQLite), abstract storage, and event-driven architecture.

## Quick Reference

```bash
# Development
go run ./cmd/vayload                 # Run the main server (dev mode)
go test ./... -short                 # Run unit tests
go test ./...                        # Full tests (including integration)
go fmt ./...                         # Format code
golangci-lint run                    # Lint (recommended)
make dev                             # Start with hot reload (using Air)
make build                           # Build production binary
make test                            # Run all tests + lint
make migrate                         # Run database migrations

# Useful commands
go generate ./...                    # Run code generation if any
air                                  # Live reload (configured in .air.toml)
```

## Project Structure

```bash
.
├── cmd/                    # Application entrypoints
├── internal/               # All private code (never import outside)
│   ├── kernel/             # Core system (DI, events, lifecycle, Lua VM)
│   ├── modules/            # Business modules (auth, database, storage, plugin-manager, ...)
│   ├── shared/             # Shared utilities, errors, security, snowflake IDs
│   ├── vayload/            # Background jobs, queues, and workers
│   └── types/              # Project-wide types
├── pkg/                    # Public packages (if published)
├── public/                 # Built Svelte 5 frontend (shadcn + Tailwind)
├── config/                 # Configuration examples
├── migrations/             # Database migrations
├── tests/                  # Integration and e2e tests
├── benchmarks/
├── scripts/
├── docs/
└── .agents/
```

### Detailed `internal/` Structure

```bash
internal/
├── kernel/                 # Foundational layer — keep lightweight
│   ├── container.go
│   ├── events.go
│   ├── kernel.go
│   ├── lifecycle.go
│   ├── lua_vm.go
│   └── services.go
│
├── modules/
│   ├── auth/                   # Authentication & RBAC
│   │   ├── domain/
│   │   ├── application/        # Use cases (login, registration, recovery, analytics)
│   │   ├── infrastructure/
│   │   └── transport/          # HTTP handlers, routes, middlewares
│   ├── database/               # Query builder, schema, multi-driver support
│   ├── plugin-manager/         # Dynamic plugin loading
│   ├── storage/                # Local + S3 storage abstraction
│   └── ...                     # Future: content, media, forms, etc.
│
├── shared/
│   ├── cache/
│   ├── ds/
│   ├── errors/
│   ├── security/
│   └── snowflake/
│
├── vayload/                # Queue system and background runtime
│   ├── queue.go
│   ├── runtime.go
│   └── ...
│
└── types/
```

## Where to Look

| Task                              | Location                                      | Notes |
|-----------------------------------|-----------------------------------------------|-------|
| Dependency Injection              | `internal/kernel/container.go`                | Central container |
| Event system                      | `internal/kernel/events.go`                   | Use for inter-module communication |
| Application lifecycle             | `internal/kernel/lifecycle.go`                | Start/Stop hooks |
| Lua scripting                     | `internal/kernel/lua_vm.go`                   | For user extensions |
| Auth domain & use cases           | `internal/modules/auth/domain/` & `application/` | Entities, errors, login/registration |
| HTTP routes & handlers            | `internal/modules/auth/transport/http/`       | All API endpoints |
| Database drivers & query builder  | `internal/modules/database/`                  | MySQL, Postgres, SQLite |
| Storage engines                   | `internal/modules/storage/engines/`           | local.go, s3.go |
| Background jobs / queues          | `internal/vayload/queue.go` & `runtime.go`    | Worker system |
| Shared errors & security          | `internal/shared/errors/` & `security/`       | JWT, hashing, policies |
| ID generation                     | `internal/shared/snowflake/`                  | Snowflake IDs |
| Plugin loading                    | `internal/modules/plugin-manager/`            | Manifest & loader |

## Core Principles

1. **Clean Architecture First** — Strict separation: `domain` → `application` → `infrastructure` → `transport`
2. **Modules are independent** — Communicate only via kernel events or the DI container. No direct imports between modules.
3. **Business logic belongs in `application/`** — Never in handlers or infrastructure.
4. **Correctness over cleverness** — Prefer readable, maintainable code.
5. **Extensibility via Lua** — When possible, expose features through the embedded Lua VM instead of hard-coding.
6. **Every significant change needs tests** — Unit + integration when appropriate.
7. **Keep the kernel lightweight** — It should not depend on any specific module.
8. **Use events for loose coupling** — Especially between modules.

## Development Workflow

- Follow the folder structure of the `auth` module when creating new modules.
- Always inject dependencies through the kernel container.
- Prefer composition over inheritance.
- Use `snake_case` for file and package names.
- Interface names should end with `er` (e.g. `UserRepository`, `StorageEngine`).
- All errors should be defined in the `domain/errors.go` of the relevant module when possible.

## Guides (to be created)

- Testing Strategy
- Creating a New Module
- Adding Database Support
- Plugin Development
- Frontend Integration

---