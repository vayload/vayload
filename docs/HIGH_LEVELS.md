# Architecture Model

## Overview

The platform follows a microkernel architecture composed of three execution layers:

1. Kernel
2. Internal Modules
3. External Plugins

The kernel is the only mandatory runtime component.

Business capabilities are implemented as modules that communicate through the kernel event bus.

External extensions are implemented as sandboxed Lua plugins executed through a controlled runtime.

---

# Layer 1 — Kernel

The kernel is the public runtime contract of the platform.

Responsibilities:

* module lifecycle
* service registry
* dependency resolution
* event bus
* capability routing
* security context propagation
* observability
* internal contracts
* runtime orchestration

The kernel does not implement business logic.

Modules depend on the kernel, never the opposite.

---

# Layer 2 — Internal Modules

Modules implement domain capabilities.

Examples:

* auth
* content
* workflow
* storage
* billing
* search

Modules are isolated runtime units managed by the kernel.

A module may:

* expose APIs
* publish events
* consume events
* register commands for call by IPC (unix: socket, win32: named pipes)
* expose contracts
* declare dependencies

Modules communicate through the kernel event bus.

Direct module-to-module coupling is discouraged.

---

# Layer 3 — External Plugins

Plugins are untrusted external extensions written in Lua.

Plugins are executed through a sandboxed runtime controlled by the platform.

Plugins cannot directly access:

* filesystem
* network
* secrets
* kernel internals

Plugins interact only through restricted APIs exposed by the plugin runtime.

---

# Runtime Patterns

Patterns are reusable runtime capabilities, not business services.

Examples:

* event bus
* queue
* scheduler
* rate limiter
* cache
* workflow state machine

Patterns provide operational behavior shared across modules.

---

# Infrastructure Adapters

Infrastructure integrations are implemented through adapters.

Examples:

* postgres adapter
* redis adapter
* s3 adapter
* grpc adapter
* nats adapter

The platform depends on internal contracts, not vendor implementations.

Adapters can be replaced without affecting modules.

---

# Communication Model

## Internal Communication

Preferred mechanisms:

* event bus
* internal contracts
* async messaging
* IPC for CLI

Optional:

* gRPC
* in-memory runtime calls

---

## External Communication

Public interfaces:

* HTTP APIs
* WebSocket channels
* MCP tools

---

# Security Model

The platform follows a zero-trust internal architecture.

Every operation validates:

* identity
* tenant
* permissions
* execution context

Plugins operate with explicit capabilities and least privilege.

No component is trusted implicitly.
