# Vayload - Complete System Architecture

## Overview

Vayload is a **microkernel-based, plugin-extensible** system built with:
- **Go**: Core services and microkernel
- **Lua**: Plugin system for extensibility
- **Svelte**: Admin SPA (no SEO required)

---

## Architecture Diagrams

### Main System Architecture

```mermaid
graph TB
    %% ================= CLIENTS =================
    subgraph ClientLayer[Client Layer]
        WEB[Web Clients]
        MOBILE[Mobile Apps]
        ADMIN[Admin SPA - Svelte]
    end

    %% ================= TRANSPORTS =================
    subgraph TransportLayer[Transport Layer]
        HTTP[HTTP / REST API Gateway]
        GRPC[gRPC Internal Comms]
        WS[WebSocket Real-time]
        MCP[MCP Server AI Tools]
        CLI[CLI vayload-cli]
    end

    %% ================= GLOBAL MIDDLEWARES =================
    subgraph GlobalMiddlewareLayer[Global Middleware Layer]
        RATE[Rate Limiting]
        TRACE[Tracing / Observability]
    end

    %% ================= SERVICE FIREWALLS =================
    subgraph ServiceMiddlewareLayer[Service Firewalls / Policies]
        SVC_POLICY[Policy Enforcement / Auth / RBAC]
        SVC_GUARD[Origin Guards]
    end

    %% ================= MICROKERNEL =================
    subgraph MicrokernelCore[Microkernel Core - Go]
        direction TB
        BOOT[Bootstrap / Startup Order]
        REG[Service Registry / DI]
        DISPATCH[Request Dispatcher / Router]
        EVENTBUS[Event Bus Pub/Sub]
        LIFECYCLE[Lifecycle Manager]
    end

    %% ================= CORE SERVICES =================
    subgraph CoreServices[Core Services - Go]
        direction LR
        AUTH_SVC[Auth Service JWT / RBAC]
        USER[User Service]
        DB[Database Service]
        CACHE[Cache Service Redis / Memory]
        STORAGE[Storage Service S3 / Local]
        NOTIFY[Notification Service]
        METRICS[Metrics / Logging]
        JOBQUEUE[Job Queue / Background Tasks]
    end

    %% ================= DATABASE BACKENDS =================
    subgraph DatabaseBackends[Database Backends]
        SQLITE[(SQLite)]
        MYSQL[(MySQL / MariaDB)]
        POSTGRES[(PostgreSQL)]
    end

    %% ================= DATA LAYER =================
    subgraph DataLayer[Data Layer]
        REDIS[(Redis Cache / Queue)]
        S3[(S3 / Minio Object Storage)]
        WAL[(WAL Write-Ahead Log)]
    end

    %% ================= CLIENT → TRANSPORT =================
    WEB --> HTTP
    MOBILE --> HTTP
    ADMIN --> HTTP
    ADMIN -.WebSocket.-> WS
    ADMIN -.Agents.->MCP

    %% ================= TRANSPORT → GLOBAL MW =================
    HTTP --> RATE
    GRPC --> RATE
    WS --> RATE
    MCP --> RATE
    CLI --> RATE

    RATE --> TRACE

    %% ================= GLOBAL → SERVICE MW =================
    TRACE --> SVC_POLICY
    SVC_POLICY --> SVC_GUARD

    %% ================= SERVICE MW → KERNEL =================
    SVC_GUARD --> DISPATCH

    %% ================= BOOTSTRAP ORDER =================
    BOOT --> REG
    BOOT --> DB
    BOOT --> CACHE
    BOOT --> AUTH_SVC
    BOOT --> USER
    BOOT --> STORAGE
    BOOT --> NOTIFY
    BOOT --> METRICS
    BOOT --> JOBQUEUE
    BOOT --> LIFECYCLE

    %% ================= KERNEL ORCHESTRATION =================
    DISPATCH --> AUTH_SVC
    DISPATCH --> REG
    DISPATCH --> EVENTBUS
    REG --> LIFECYCLE

    %% ================= SERVICE REGISTRATION =================
    REG --> AUTH_SVC
    REG --> USER
    REG --> DB
    REG --> CACHE
    REG --> STORAGE
    REG --> NOTIFY
    REG --> METRICS
    REG --> JOBQUEUE

    %% ================= DATABASE SERVICE =================
    DB --> SQLITE
    DB --> MYSQL
    DB --> POSTGRES

    %% ================= DATA CONNECTIONS =================
    CACHE --> REDIS
    STORAGE --> S3
    METRICS --> POSTGRES
    JOBQUEUE --> REDIS
    JOBQUEUE --> WAL

    %% ================= EVENT BUS =================
    EVENTBUS -.Events.-> USER
    EVENTBUS -.Events.-> NOTIFY
    EVENTBUS -.Events.-> METRICS
    EVENTBUS -.Events.-> JOBQUEUE

    %% ================= STYLING =================
    classDef kernelClass fill:#ff6b6b,stroke:#c92a2a,stroke-width:3px,color:#fff,stroke-dasharray: 5 5
    classDef serviceClass fill:#4ecdc4,stroke:#0a8577,stroke-width:2px,color:#fff
    classDef transportClass fill:#a8dadc,stroke:#457b9d,stroke-width:2px,color:#000
    classDef middlewareClass fill:#ffa94d,stroke:#e8590c,stroke-width:2px,color:#000
    classDef dataClass fill:#9d84b7,stroke:#6a4c93,stroke-width:2px,color:#fff
    classDef clientClass fill:#95e1d3,stroke:#38ada9,stroke-width:2px,color:#000

    class BOOT,REG,DISPATCH,EVENTBUS,LIFECYCLE kernelClass
    class AUTH_SVC,USER,DB,CACHE,STORAGE,NOTIFY,METRICS,JOBQUEUE serviceClass
    class HTTP,GRPC,WS,MCP,CLI transportClass
    class RATE,TRACE,SVC_POLICY,SVC_GUARD middlewareClass
    class SQLITE,MYSQL,POSTGRES,REDIS,S3,WAL dataClass
    class WEB,MOBILE,ADMIN clientClass
```

---

### Plugin System Architecture

```mermaid
graph TB
    %% ================= MICROKERNEL =================
    subgraph MicrokernelCore[Microkernel Core - Go]
        direction TB
        REG[Service Registry / DI]
        DISPATCH[Request Dispatcher / Router]
        EVENTBUS[Event Bus Pub/Sub]
        LIFECYCLE[Lifecycle Manager]
        BRIDGE[Plugin API Bridge]
    end

    %% ================= PLUGIN SYSTEM =================
    subgraph PluginSystem[Plugin System - Lua]
        direction TB
        PLUGMGR[Plugin Manager Load/Unload/Reload]
        LUAVM[Lua VM Pool gopher-lua]
        SANDBOX[Sandbox / Resource Limits CPU/Memory/IO]

        subgraph LuaPlugins[Lua Plugins]
            direction LR
            P1[Payment Plugin]
            P2[Webhook Plugin]
            P3[Custom Logic]
            P4[Integrations]
        end
    end

    %% ================= BRIDGE CAPABILITIES =================
    subgraph BridgeCapabilities[Bridge Exposed APIs]
        direction TB
        BRIDGE_DB[Database Access]
        BRIDGE_CACHE[Cache Access]
        BRIDGE_HTTP[HTTP Client]
        BRIDGE_EVENTS[Event Pub/Sub]
        BRIDGE_TRANSPORT[Transport Registration REST/gRPC/WS]
        BRIDGE_STORAGE[Storage Access]
        BRIDGE_NOTIFY[Notification API]
        BRIDGE_JOBS[Job Queue Enqueue/Status]
    end

    %% ================= DATA LAYER =================
    subgraph DataLayer[Data Layer]
        REDIS[(Redis Cache / Queue)]
        S3[(S3 / Minio)]
        WAL[(WAL Write-Ahead Log)]
    end

    %% ================= DATABASE BACKENDS =================
    subgraph DatabaseBackends[Database Backends]
        SQLITE[(SQLite)]
        MYSQL[(MySQL)]
        POSTGRES[(PostgreSQL)]
    end

    %% ================= KERNEL ORCHESTRATION =================
    DISPATCH --> REG
    DISPATCH --> EVENTBUS
    REG --> LIFECYCLE

    %% ================= PLUGIN SYSTEM =================
    REG --> PLUGMGR
    PLUGMGR --> LUAVM
    LUAVM --> SANDBOX
    SANDBOX --> P1
    SANDBOX --> P2
    SANDBOX --> P3
    SANDBOX --> P4

    %% ================= BRIDGE CONNECTIONS =================
    BRIDGE --> BRIDGE_DB
    BRIDGE --> BRIDGE_CACHE
    BRIDGE --> BRIDGE_HTTP
    BRIDGE --> BRIDGE_EVENTS
    BRIDGE --> BRIDGE_TRANSPORT
    BRIDGE --> BRIDGE_STORAGE
    BRIDGE --> BRIDGE_NOTIFY
    BRIDGE --> BRIDGE_JOBS

    %% ================= BRIDGE TO RESOURCES =================
    BRIDGE_DB --> SQLITE
    BRIDGE_DB --> MYSQL
    BRIDGE_DB --> POSTGRES
    BRIDGE_CACHE --> REDIS
    BRIDGE_STORAGE --> S3
    BRIDGE_EVENTS --> EVENTBUS
    BRIDGE_JOBS --> REDIS
    BRIDGE_JOBS --> WAL

    %% ================= PLUGINS USE BRIDGE ONLY =================
    P1 -->|All access via Bridge| BRIDGE
    P2 -->|All access via Bridge| BRIDGE
    P3 -->|All access via Bridge| BRIDGE
    P4 -->|All access via Bridge| BRIDGE

    %% ================= EVENT BUS TO PLUGINS =================
    EVENTBUS -.Events.-> P1
    EVENTBUS -.Events.-> P2
    EVENTBUS -.Events.-> P3
    EVENTBUS -.Events.-> P4

    %% ================= TRANSPORT REGISTRATION =================
    BRIDGE_TRANSPORT -.Registers.-> DISPATCH

    %% ================= STYLING =================
    classDef kernelClass fill:#ff6b6b,stroke:#c92a2a,stroke-width:3px,color:#fff
    classDef pluginClass fill:#ffe66d,stroke:#f5a623,stroke-width:2px,color:#000
    classDef bridgeClass fill:#51cf66,stroke:#2f9e44,stroke-width:3px,color:#fff
    classDef dataClass fill:#9d84b7,stroke:#6a4c93,stroke-width:2px,color:#fff

    class REG,DISPATCH,EVENTBUS,LIFECYCLE kernelClass
    class PLUGMGR,LUAVM,SANDBOX,P1,P2,P3,P4 pluginClass
    class BRIDGE,BRIDGE_DB,BRIDGE_CACHE,BRIDGE_HTTP,BRIDGE_EVENTS,BRIDGE_TRANSPORT,BRIDGE_STORAGE,BRIDGE_NOTIFY,BRIDGE_JOBS bridgeClass
    class SQLITE,MYSQL,POSTGRES,REDIS,S3,WAL dataClass
```

---

## System Components

### 1. Client Layer
- **Web Clients**: Browser-based applications
- **Mobile Apps**: Native iOS/Android applications
- **Admin SPA**: Svelte-based admin panel (no SEO requirements)

### 2. Transport Layer
Multiple protocol support for different use cases:

| Transport | Use Case | Protocol |
|-----------|----------|----------|
| **HTTP/REST** | Public API, external integrations | HTTP/1.1, HTTP/2 |
| **gRPC** | Internal service-to-service communication | HTTP/2 |
| **WebSocket** | Real-time bidirectional communication | WS/WSS |
| **MCP** | AI tool integration for agents | MCP Protocol |
| **CLI** | Administrative tools and automation | Direct process |

### 3. Middleware Layers

#### Global Middleware
Applied to **all** incoming requests:
- **Rate Limiting**: DDoS protection, quota enforcement
- **Tracing/Observability**: Distributed tracing, metrics collection

#### Service Middleware (Firewalls)
Applied per-service basis:
- **Policy Enforcement**: RBAC, permission checks
- **Auth**: JWT validation, session management
- **Origin Guards**: CORS, allowed origins validation

### 4. Microkernel Core (Go)

The minimal, essential kernel that coordinates everything:

```go
type Kernel struct {
    Bootstrap  *BootstrapManager  // Startup order orchestration
    Registry   *ServiceRegistry   // Dependency injection
    Dispatcher *RequestDispatcher // Request routing
    EventBus   *EventBus          // Pub/Sub messaging
    Lifecycle  *LifecycleManager  // Health checks, graceful shutdown
}
```

**Bootstrap Order**:
1. Database Service (supports SQLite, MySQL, PostgreSQL)
2. Cache Service (Memory/Redis)
3. Auth Service (JWT/RBAC)
4. Core Services (User, Storage, Notify, Metrics)
5. Job Queue (WAL/Redis backed)
6. Plugin Manager
7. Lifecycle Manager

### 5. Core Services (Go)

Production-critical services implemented in Go for performance:

| Service | Responsibility |
|---------|---------------|
| **Auth Service** | JWT token management, RBAC enforcement |
| **User Service** | User CRUD, profile management |
| **Database Service** | Multi-backend abstraction (SQLite/MySQL/PostgreSQL) |
| **Cache Service** | Redis/Memory caching, TTL management |
| **Storage Service** | S3/Minio object storage |
| **Notification Service** | Email, SMS, push notifications |
| **Metrics/Logging** | OpenTelemetry, structured logging |
| **Job Queue** | Background task processing, scheduled jobs |

### 6. Plugin System (Lua)

Extends functionality without recompiling the core:

#### Plugin Manager
```go
type PluginManager struct {
    vmPool   *LuaVMPool       // Pool of gopher-lua VMs
    sandbox  *Sandbox         // Resource limits enforcement
    plugins  map[string]*Plugin
}

// Hot reload without downtime
func (pm *PluginManager) Reload(name string) error
```

#### Sandbox Limits
- **CPU**: Max execution time per request
- **Memory**: Heap size limits
- **I/O**: Rate limiting on network/disk operations
- **API calls**: Quota on bridge API calls

#### Available Plugins
- **Payment Plugin**: Stripe, PayPal integrations
- **Webhook Plugin**: Outbound webhook delivery
- **Custom Logic**: Business-specific rules
- **Integrations**: Third-party API wrappers

---

## Plugin API Bridge

The **ONLY** interface plugins use to interact with the system.

### Bridge Capabilities

```lua
-- Database Access (all supported backends)
local user = bridge.db.query("SELECT * FROM users WHERE id = ?", user_id)
bridge.db.exec("UPDATE users SET name = ? WHERE id = ?", "Alice", 1)
bridge.db.transaction(function(tx)
    tx.exec("INSERT INTO logs ...")
    tx.exec("UPDATE counters ...")
end)

-- Cache Access
bridge.cache.set("session:123", data, 3600) -- TTL in seconds
local cached = bridge.cache.get("session:123")
bridge.cache.delete("session:123")

-- HTTP Client
local response = bridge.http.get("https://api.example.com/data")
local result = bridge.http.post("https://webhook.site", {
    body = json.encode({event = "payment"}),
    headers = {["Content-Type"] = "application/json"}
})

-- Event Pub/Sub
bridge.events.publish("user.created", {id = 123, name = "Alice"})
bridge.events.subscribe("payment.received", function(event)
    -- Handle event
end)

-- Transport Registration (Expose REST/gRPC/WS endpoints)
bridge.transport.register_http("GET", "/plugins/payment/status", function(req)
    return {
        status = 200,
        body = json.encode({status = "active"})
    }
end)

bridge.transport.register_grpc("PaymentService", "GetStatus", function(req)
    return {status = "active"}
end)

bridge.transport.register_websocket("/ws/notifications", {
    on_connect = function(conn) end,
    on_message = function(conn, msg) end,
    on_close = function(conn) end
})

-- Storage Access
local url = bridge.storage.upload("invoices/2024/invoice.pdf", data)
local content = bridge.storage.download("invoices/2024/invoice.pdf")
bridge.storage.delete("invoices/2024/invoice.pdf")

-- Notifications
bridge.notify.email({
    to = "user@example.com",
    subject = "Payment Received",
    body = "Your payment was processed"
})

bridge.notify.sms({
    to = "+1234567890",
    body = "OTP: 123456"
})

-- Job Queue (Background Tasks)
local job_id = bridge.jobs.enqueue("send_invoice_email", {
    user_id = 123,
    invoice_id = 456
}, {
    delay = 300,      -- 5 minutes delay
    retry = 3,        -- retry up to 3 times
    timeout = 60      -- 60 seconds timeout
})

-- Check job status
local status = bridge.jobs.status(job_id)
-- status = {state = "pending|running|completed|failed", attempts = 1, ...}

-- Schedule recurring job (cron-like)
bridge.jobs.schedule("daily_report", "0 0 * * *", function()
    -- Generate and send daily report
end)
```

### Security Model

1. **Sandboxing**: Each plugin runs in isolated Lua VM
2. **Resource Limits**: CPU/memory/I/O quotas enforced
3. **API Quotas**: Bridge calls rate-limited per plugin
4. **No Direct Access**: Plugins cannot access Go internals
5. **Capability-based**: Plugins declare required capabilities

---

## Data Layer

### Supported Databases
- **SQLite**: Embedded, single-node deployments
- **MySQL/MariaDB**: Traditional RDBMS
- **PostgreSQL**: Advanced features, JSON support

### Caching
- **Redis**: Distributed cache, pub/sub, queues
- **Memory**: In-process cache for single-node

### Object Storage
- **S3**: AWS S3 compatible
- **Minio**: Self-hosted object storage

### Job Queue / Background Tasks
- **Backend Options**:
  - **Redis**: Distributed queue with pub/sub (recommended for multi-node)
  - **WAL (Write-Ahead Log)**: Embedded, file-based queue (single-node, SQLite-like durability)
- **Features**:
  - Delayed jobs
  - Retries with exponential backoff
  - Priority queues
  - Cron-like scheduling
  - Dead letter queue
  - Job status tracking

---

## Event Bus Architecture

Asynchronous, decoupled communication:

```go
type Event struct {
    Topic     string                 // e.g., "user.created"
    Payload   map[string]interface{} // Event data
    Timestamp time.Time
    TraceID   string                 // For distributed tracing
}

// Core services publish events
eventBus.Publish("user.created", map[string]interface{}{
    "id": 123,
    "email": "user@example.com",
})

// Plugins subscribe to events
bridge.events.subscribe("user.created", function(event)
    -- Send welcome email via plugin
end)
```

**Event Flow**:
1. Core service emits event → Event Bus
2. Event Bus broadcasts to all subscribers
3. Plugins/services react asynchronously
4. No tight coupling between components

---

## Request Flow

### 1. External Request Flow
```
Client → HTTP Transport → Rate Limiting → Tracing →
Policy Enforcement → Origin Guards → Dispatcher →
Auth Service (verify JWT) → Registry → Target Service → Response
```

### 2. Plugin-Initiated Request Flow
```
Lua Plugin → Bridge API → Registry → Core Service →
Database/Cache/Storage → Response → Plugin
```

### 3. Event-Driven Flow
```
Core Service → Event Bus → [Multiple Subscribers] →
Plugins/Services react independently
```

---

## Deployment Patterns

### Single Node (SQLite + Memory Cache)
```bash
vayload serve \
  --db-type=sqlite \
  --db-path=./data.db \
  --cache-type=memory \
  --queue-type=wal \
  --queue-path=./queue.wal \
  --plugins-dir=./plugins
```

### Multi-Node (PostgreSQL + Redis)
```bash
vayload serve \
  --db-type=postgres \
  --db-url=postgres://user:pass@db:5432/vayload \
  --cache-type=redis \
  --cache-url=redis://cache:6379 \
  --queue-type=redis \
  --queue-url=redis://queue:6379 \
  --storage-type=s3 \
  --storage-bucket=vayload-prod \
  --plugins-dir=./plugins
```

---

## Plugin Development Example

### Payment Plugin (`plugins/payment/main.lua`)

```lua
-- Plugin metadata
local plugin = {
    name = "payment",
    version = "1.0.0",
    capabilities = {"db", "http", "events", "transport"}
}

-- Initialize plugin
function plugin.init()
    -- Subscribe to events
    bridge.events.subscribe("order.created", handle_order)

    -- Register HTTP endpoint
    bridge.transport.register_http("POST", "/plugins/payment/process", process_payment)
end

-- Event handler
function handle_order(event)
    local order_id = event.payload.order_id
    local amount = event.payload.amount

    -- Query database
    local order = bridge.db.query_one(
        "SELECT * FROM orders WHERE id = ?",
        order_id
    )

    -- Call external payment API
    local response = bridge.http.post("https://api.stripe.com/v1/charges", {
        headers = {
            ["Authorization"] = "Bearer " .. STRIPE_KEY,
            ["Content-Type"] = "application/x-www-form-urlencoded"
        },
        body = "amount=" .. amount .. "&currency=usd"
    })

    -- Update database
    if response.status == 200 then
        bridge.db.exec(
            "UPDATE orders SET status = 'paid' WHERE id = ?",
            order_id
        )

        -- Publish success event
        bridge.events.publish("payment.completed", {
            order_id = order_id,
            transaction_id = response.body.id
        })
    else
        bridge.events.publish("payment.failed", {
            order_id = order_id,
            error = response.body.error
        })
    end
end

-- HTTP endpoint handler
function process_payment(req)
    local body = json.decode(req.body)

    -- Process payment logic
    local result = charge_card(body.amount, body.token)

    return {
        status = 200,
        headers = {["Content-Type"] = "application/json"},
        body = json.encode(result)
    }
end

return plugin
```

---

## Configuration

### Main Config (`config.yaml`)

```yaml
server:
  host: 0.0.0.0
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

database:
  type: postgres  # sqlite, mysql, postgres
  url: postgres://user:pass@localhost:5432/vayload
  max_connections: 100
  max_idle: 10

cache:
  type: redis  # memory, redis
  url: redis://localhost:6379
  ttl: 3600

storage:
  type: s3  # local, s3, minio
  bucket: vayload-uploads
  region: us-east-1

queue:
  type: redis  # wal, redis
  url: redis://localhost:6379  # only for redis
  path: ./queue.wal             # only for wal
  workers: 4                    # concurrent workers
  retry:
    max_attempts: 3
    backoff: exponential        # exponential, linear
    initial_delay: 1s
    max_delay: 5m

auth:
  jwt_secret: your-secret-key
  token_ttl: 24h
  refresh_ttl: 168h

plugins:
  dir: ./plugins
  hot_reload: true
  sandbox:
    max_memory_mb: 128
    max_cpu_ms: 1000
    max_io_mb: 10

middleware:
  rate_limit:
    enabled: true
    requests_per_minute: 60
  tracing:
    enabled: true
    exporter: jaeger
    endpoint: http://localhost:14268/api/traces

observability:
  metrics:
    enabled: true
    port: 9090
  logging:
    level: info
    format: json
```

---

## Advantages of This Architecture

### ✅ **Microkernel Benefits**
- Minimal core = easier to maintain
- Core is stable, plugins evolve independently
- Easy to reason about system behavior

### ✅ **Plugin Extensibility**
- Add features without recompiling
- Hot reload plugins without downtime
- Lua is simple for non-Go developers

### ✅ **Security**
- Plugins sandboxed with resource limits
- Bridge API is the only interface
- No direct access to Go internals

### ✅ **Multi-Transport**
- One business logic, multiple protocols
- Easy to add new transports
- Consistent middleware pipeline

### ✅ **Database Flexibility**
- Single abstraction layer
- Easy to switch backends
- Multi-tenancy support

### ✅ **Event-Driven**
- Loose coupling between components
- Asynchronous processing
- Easy to scale horizontally

### ✅ **Production Ready**
- OpenTelemetry for observability
- Distributed tracing
- Health checks and graceful shutdown
- Hot reload without downtime

---

## Development Workflow

### 1. Core Service Development (Go)
```bash
# Add new core service
cd internal/services
mkdir analytics
cd analytics
# Implement Service interface
# Register in kernel bootstrap
```

### 2. Plugin Development (Lua)
```bash
# Create plugin directory
mkdir -p plugins/my-plugin
cd plugins/my-plugin

# Create main.lua
cat > main.lua <<EOF
local plugin = {name = "my-plugin", version = "1.0.0"}
function plugin.init()
    bridge.transport.register_http("GET", "/my-endpoint", handler)
end
return plugin
EOF

# Hot reload
vayload plugins reload my-plugin
```

### 3. Admin SPA (Svelte)
```bash
cd admin-spa
npm run dev
# Connects to backend via HTTP/WebSocket
```

---

## Monitoring & Observability

### Metrics (Prometheus)
```
vayload_requests_total{method="GET",path="/api/users",status="200"} 1523
vayload_plugin_execution_seconds{plugin="payment"} 0.045
vayload_db_query_duration_seconds{query="select_user"} 0.003
```

### Tracing (Jaeger)
```
Request ID: trace-123
├─ HTTP Handler (2ms)
├─ Auth Middleware (1ms)
├─ Dispatcher (0.5ms)
├─ User Service (5ms)
│  └─ Database Query (3ms)
└─ Response (0.5ms)
Total: 9ms
```

### Logging (Structured JSON)
```json
{
  "level": "info",
  "ts": "2026-01-23T10:30:00Z",
  "caller": "dispatcher/handler.go:45",
  "msg": "request processed",
  "trace_id": "trace-123",
  "method": "GET",
  "path": "/api/users",
  "duration_ms": 9,
  "status": 200
}
```

---

## Future Enhancements

- [ ] gRPC streaming support
- [ ] GraphQL transport layer
- [ ] Multi-region replication
- [ ] Plugin marketplace
- [ ] Visual plugin builder (low-code)
- [ ] Advanced RBAC with attribute-based policies
- [ ] Kubernetes operator for auto-scaling
- [ ] Built-in API versioning
- [ ] Rate limiting per user/tenant
- [ ] Circuit breaker pattern for external calls

---

## License

MIT License - Copyright (c) 2026 Alex Zweiter
