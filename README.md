# 🐷 AI Gamification Elearning-Server — Backend

> Open source backend for the **AI Gamification Elearning Server** e-learning platform: a gamified, AI-powered environment to learn web and backend development.

---

## 📖 About the project

**AI Gamification Elearning Server** is an open source, gamified e-learning platform designed to teach software development in a practical and engaging way. The platform's initial focus is **web and backend development**, with a roadmap to expand into mobile development and artificial intelligence.

The platform is built around two core ideas:

- **Learning by doing** — students write real code, solve real challenges, and progress through levels as they build skills.
- **AI-powered feedback** — an integrated AI code reviewer analyses the student's submissions and provides contextual, constructive feedback in real time, acting as a always-available mentor.

This repository is the **Go backend** of the AI Gamification Elearning Server platform. It exposes the REST API consumed by the frontend and handles all business logic, persistence, and integration with external services — including the AI code review engine.

This project is **free to use and open to contributions**. See the [license](#license) section for details.

---

## 🗺️ Roadmap highlights

| Phase | Scope |
|---|---|
| v1 | Web & backend development tracks |
| v2 | AI code reviewer integration |
| v3 | Mobile development track |
| v4 | AI / ML development track |

---

## 🏗️ Architecture

The backend follows a **Hexagonal Architecture** (also known as Ports and Adapters), adapted for Go idioms. The guiding principle is that dependency arrows always point **inward** — outer layers know about inner layers, never the reverse.

```
Adapters ──► Infrastructure ──► Application ──► Domain
```

The domain layer has zero knowledge of Gin, GORM, Redis, or any other external library.

---

## 📁 Full project structure

```
ai-gamification-elearning-server/
├── cmd/
│   └── main.go                    # Entry point. Loads config, builds the DI container, starts the server.
│                                  # Intentionally thin — zero business logic lives here.
│
├── configs/
│   └── config.yaml                # Default configuration values. Production secrets are injected via
│                                  # environment variables — never committed to the repository.
│
├── pkg/                           # [GLOBAL] Cross-cutting utilities. Can be imported by any layer.
│   │                              # Rule: if two or more layers need the same helper, it belongs here.
│   │
│   ├── logger/                    # Structured logger wrapper (e.g. over zap or slog).
│   │                              # Centralises log level, format (JSON/text), and output target.
│   │                              # Rule: initialise once in main.go and inject. Never use fmt.Println in production.
│   │
│   ├── timeutil/                  # Date and time helpers: ISO 8601 parsing, timestamp formatting,
│   │                              # timezone conversion, and age calculation.
│   │                              # Rule: never hardcode time.Parse layouts outside this package.
│   │
│   ├── pagination/                # Shared pagination structs (PageRequest, PageResponse) and helpers
│   │                              # for computing offsets, total pages, and building paginated responses.
│   │                              # Rule: pagination logic lives here — not inside GORM queries or Gin handlers.
│   │
│   ├── crypto/                    # Hashing utilities (bcrypt, SHA-256), secure random token generation,
│   │                              # and JWT signing/verification helpers.
│   │                              # Rule: no raw crypto calls outside this package.
│   │
│   ├── validator/                 # Wrapper around go-playground/validator. Registers custom validation
│   │                              # rules (e.g. valid-username, strong-password).
│   │                              # Rule: controllers only call validator.Validate() — tag logic stays here.
│   │
│   └── config/                    # Configuration management using Viper and dotenv.
│                                  # Handles environment variable substitution, YAML parsing,
│                                  # and mapping external settings to internal Go structures.
│                                  # Rule: centralized source of truth for app settings. Avoid 
│                                  # hardcoding environment variable lookups (os.Getenv) elsewhere.
│
├── internal/                      # Go internal boundary. Prevents any code outside this module from
│   │                              # importing these packages — enforcing the hexagonal boundary at language level.
│   │
│   ├── domain/                    # [DOMAIN LAYER] The innermost layer. Defines what the system IS,
│   │   │                          # not how it works technically. Zero third-party imports.
│   │   │                          # Rule: importing any external library here is an architecture violation.
│   │   │
│   │   ├── entities/              # Plain Go structs representing core business concepts:
│   │   │                          # User, Course, Challenge, Submission, Badge.
│   │   │                          # No GORM tags, no JSON tags — pure business data shapes.
│   │   │                          # Rule: entities have no methods that call external systems.
│   │   │
│   │   ├── ports/                 # Go interfaces defining contracts between the domain and the outside world.
│   │   │                          # Examples: UserRepository, CachePort, NotificationPort.
│   │   │                          # The application layer depends on these; infrastructure implements them.
│   │   │                          # Rule: ports are named by what they do, not how (UserRepository not PostgresRepo).
│   │   │
│   │   └── errors/                # Typed, semantic domain errors: ErrUserNotFound, ErrEmailAlreadyExists,
│   │                              # ErrInvalidSubmission. Carry business meaning only — no HTTP status codes.
│   │                              # Rule: infrastructure maps these errors to HTTP codes, not the other way around.
│   │
│   ├── application/               # [APPLICATION LAYER] Orchestrates the domain to fulfil use cases.
│   │   │                          # Knows about entities and ports. Completely unaware of HTTP or GORM.
│   │   │                          # Rule: no *gin.Context or *gorm.DB anywhere in this layer.
│   │   │
│   │   ├── usecases/              # One subfolder per aggregate (user/, course/, challenge/).
│   │   │                          # Each file is a single use case struct with an Execute() method.
│   │   │                          # Business rules, invariant checks, and port calls live here.
│   │   │                          # Rule: one use case = one file = one public Execute() method.
│   │   │
│   │   ├── dtos/                  # Data Transfer Objects — the explicit contract between HTTP and business logic.
│   │   │                          # Controllers map JSON → Request DTO. Use cases return Response DTOs.
│   │   │                          # Prevents domain entities from leaking into API responses.
│   │   │                          # Rule: DTOs are flat structs with json/validate tags. No logic inside them.
│   │   │
│   │   ├── services/              # Application-level orchestration: NotificationService, ScoringService.
│   │   │                          # Coordinate multiple use cases or ports. Distinct from domain services —
│   │   │                          # these depend on ports, not on business invariants.
│   │   │                          # Rule: services here orchestrate; they do not contain business rules.
│   │   │
│   │   └── mappers/               # Explicit mapping functions between domain entities and DTOs.
│   │                              # EntityToResponse(), RequestToEntity() — centralised to avoid scattered conversions.
│   │                              # Rule: all entity ↔ DTO transformations happen here, nowhere else.
│   │
│   ├── infrastructure/            # [INFRASTRUCTURE LAYER] All concrete implementations of domain ports.
│   │   │                          # The only layer that imports GORM, Gin, Redis, Kafka, or external SDKs.
│   │   │                          # Rule: nothing in domain or application imports from this layer.
│   │   │
│   │   ├── persistence/           # Everything related to SQL database access via GORM.
│   │   │   │                      # Rule: GORM models are internal — never return a *gorm.Model upward.
│   │   │   │
│   │   │   ├── db.go              # GORM connection factory. Reads DSN from config, runs AutoMigrate,
│   │   │   │                      # and configures connection pool (MaxOpenConns, MaxIdleConns, ConnMaxLifetime).
│   │   │   │
│   │   │   ├── models/            # GORM model structs with database-specific tags, indexes, and hooks.
│   │   │   │                      # Intentionally separate from domain entities — schema changes here do not
│   │   │   │                      # ripple into the business layer.
│   │   │   │                      # Rule: models embed gorm.Model; map to/from entities inside repositories only.
│   │   │   │
│   │   │   ├── repositories/      # Concrete implementations of domain port interfaces using GORM.
│   │   │   │                      # Each file corresponds to a port (GormUserRepository → ports.UserRepository).
│   │   │   │                      # Handles mapping between GORM models and domain entities internally.
│   │   │   │                      # Rule: repository methods return domain entities, never GORM models.
│   │   │   │
│   │   │   └── helpers/           # Persistence-specific utilities: reusable GORM scope builders,
│   │   │                          # soft-delete helpers, paginated query builders.
│   │   │                          # Rule: used only within the persistence sub-package.
│   │   │
│   │   ├── http/                  # HTTP delivery layer using Gin. The boundary between the outside world
│   │   │   │                      # and the application layer — translates HTTP ↔ DTOs.
│   │   │   │                      # Rule: controllers never contain business logic. Bind, validate, call, respond.
│   │   │   │
│   │   │   ├── controllers/       # Gin handler structs, one per aggregate.
│   │   │   │                      # Each method: (1) binds JSON to a DTO, (2) calls the use case,
│   │   │   │                      # (3) maps the result or domain error to an HTTP response.
│   │   │   │                      # Rule: domain error → HTTP status code mapping lives here via mapDomainError().
│   │   │   │
│   │   │   ├── routes/            # Gin router registration, separated from controller logic.
│   │   │   │                      # One file per aggregate: RegisterUserRoutes(), RegisterCourseRoutes().
│   │   │   │                      # Groups routes under /api/v1 and applies route-level middleware.
│   │   │   │                      # Rule: no handler logic in route files — only router.METHOD(path, handler).
│   │   │   │
│   │   │   └── helpers/           # HTTP-specific utilities: standardised JSON response builders
│   │   │                          # (Success(), Error(), Paginated()), and request context extractors.
│   │   │                          # Rule: all JSON response shapes defined here — never gin.H{} directly in controllers.
│   │   │
│   │   └── external/              # Adapters for third-party services: Redis, Kafka, email, Anthropic AI API.
│   │       │                      # Each file implements a domain port interface.
│   │       │                      # Rule: if the provider changes, only this folder changes.
│   │       │
│   │       ├── redis_cache.go     # Implements ports.CachePort. Handles connection, serialisation, and TTL.
│   │       │
│   │       ├── kafka_publisher.go # Publishes domain events to Kafka. Handles producer lifecycle and retries.
│   │       │
│   │       └── helpers/           # External-service utilities: retry with exponential backoff,
│   │                              # circuit breaker wrapper, response deserialisers.
│   │                              # Rule: retry and circuit-breaker logic lives here — not duplicated per client file.
│   │
│   └── adapters/                  # [ADAPTERS LAYER] Outermost layer. Bootstraps the application,
│       │                          # wires all dependencies, and defines entry points from the outside world.
│       │                          # Rule: instantiates things — contains no business logic of its own.
│       │
│       ├── server/                # Gin engine bootstrap. Creates the engine, registers global middleware
│       │                          # (CORS, recovery, request ID, rate limiting), and exposes Start().
│       │                          # Rule: no routes registered here — delegated to infrastructure/http/routes/.
│       │
│       ├── messaging/             # Kafka consumer entry points. Each consumer goroutine listens on a topic,
│       │                          # deserialises the message, and calls the appropriate use case with a DTO.
│       │                          # Acts as the async equivalent of an HTTP controller.
│       │                          # Rule: consumers call use cases with DTOs — no business logic here.
│       │
│       ├── middleware/            # Gin middleware: JWT authentication, role-based authorisation,
│       │                          # request logging (via pkg/logger), CORS headers, and panic recovery.
│       │                          # Rule: named by single responsibility — not grouped under a generic "utils".
│       │
│       └── di/                    # Composition root. The single place where every dependency is instantiated
│                                  # and wired together manually — no reflection-based containers.
│                                  # Reading this file reveals the entire dependency graph of the application.
│                                  # Rule: adding a feature should require wiring in exactly one place here.
│
├── go.mod                         # Go module definition. Declares module path and pinned dependency versions.
└── Makefile                       # Developer shortcuts: make run, make test, make lint, make migrate, make docker-up.
```

---
 
## 🔵 Domain layer
 
**Location:** `internal/domain/`
 
The innermost layer. It defines what the business *is*, not how it works technically. This layer has no third-party dependencies and can be tested in complete isolation.
 
```go
// SAMPLE

// entities/user.go — no GORM tags, no JSON tags
type User struct {
    ID    string
    Name  string
    Email string
    XP    int
    Level int
}
 
// ports/user_repository.go — contract, not implementation
type UserRepository interface {
    FindByID(id string) (*entities.User, error)
    Save(user *entities.User) error
}
 
// errors/domain_errors.go — business meaning, no HTTP codes
var ErrUserNotFound       = errors.New("user not found")
var ErrEmailAlreadyExists = errors.New("email already exists")
```
 
---
 
## 🟢 Application layer
 
**Location:** `internal/application/`
 
Orchestrates the domain to fulfil specific use cases. DTOs isolate the API contract from internal domain changes — a schema tweak in the database or a field rename in an entity never breaks the HTTP contract.
 
```go
// SAMPLE

// dtos/create_user_request.go
type CreateUserRequest struct {
    Name     string `json:"name"     validate:"required"`
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
 
// usecases/user/create_user.go
type CreateUserUseCase struct {
    repo  ports.UserRepository
    cache ports.CachePort
}
 
func (uc *CreateUserUseCase) Execute(req dtos.CreateUserRequest) (*dtos.UserResponse, error) {
    // 1. Validate business rules
    // 2. Build domain entity
    // 3. Persist via port (unaware of GORM)
    // 4. Return response DTO
    // No *gin.Context, no *gorm.DB
}
```
 
---
 
## 🟡 Infrastructure layer
 
**Location:** `internal/infrastructure/`
 
The only layer that imports GORM and Gin. GORM models are kept strictly separate from domain entities — the mapping is explicit and lives inside the repository.
 
```go
// SAMPLE

// persistence/models/user_model.go
type UserModel struct {
    gorm.Model
    Name  string `gorm:"not null"`
    Email string `gorm:"uniqueIndex;not null"`
    XP    int    `gorm:"default:0"`
    Level int    `gorm:"default:1"`
}
 
// http/controllers/user_controller.go
func (c *UserController) CreateUser(ctx *gin.Context) {
    var req dtos.CreateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, helpers.Error(err.Error()))
        return
    }
    resp, err := c.createUserUseCase.Execute(req)
    if err != nil {
        ctx.JSON(helpers.MapDomainError(err))
        return
    }
    ctx.JSON(http.StatusCreated, helpers.Success(resp))
}
 
// http/routes/user_routes.go
func RegisterUserRoutes(r *gin.Engine, ctrl *UserController) {
    v1 := r.Group("/api/v1")
    v1.POST("/users", ctrl.CreateUser)
    v1.GET("/users/:id", ctrl.GetUser)
}
```
 
---
 
## 🔴 Adapters layer

**Location:** `internal/adapters/`
 
Bootstraps and wires. The composition root (`di/wire.go`) is the single file that reveals the entire dependency graph — no magic, no reflection.

```go
// SAMPLE

// adapters/di/wire.go
func BuildContainer(cfg *config.Config) (*gin.Engine, error) {
    db       := persistence.NewDB(cfg)
    userRepo := repositories.NewGormUserRepository(db)
    cache    := external.NewRedisCache(cfg)
 
    createUser := userUsecase.NewCreateUserUseCase(userRepo, cache)
    getUser    := userUsecase.NewGetUserUseCase(userRepo, cache)
 
    userCtrl := controllers.NewUserController(createUser, getUser)
    server   := ginserver.NewGinServer(cfg)
    routes.RegisterUserRoutes(server, userCtrl)
 
    return server, nil
}
```

---

## ⚙️ Tech stack

| Concern | Technology |
|---|---|
| Language | Go 1.22+ |
| HTTP framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) |
| Cache | Redis |
| Messaging | Kafka |
| AI integration | Anthropic API (Claude) |
| Logger | [zap](https://github.com/uber-go/zap) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |

---

## 🚀 Getting started

```bash
# Clone the repository
git clone https://github.com/your-org/ai-gamification-elearning-server.git
cd ai-gamification-elearning-server

# Copy and configure environment
cp .env.sample .env.local
# Edit .env with your values (optional: create .env.local for local overrides)

# Run the server
make run

# Run with live-reload
make run/watch
```

---

## 🤝 Contributing

AI Gamification Elearning Server is an open source project and contributions are welcome. Please open an issue before submitting a pull request so we can discuss the proposed change first.

When contributing, make sure your code respects the layer boundaries described in this document. A PR that imports GORM into the domain layer or places business logic inside a controller will not be merged.

---

## 📄 License

This project is released under the [MIT License](LICENSE). You are free to use, modify, and distribute it.
