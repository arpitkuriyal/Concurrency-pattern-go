# The `context` Package in Go

This document explains the **theory, intent, and best practices** of Go’s `context` package. It is meant for **conceptual understanding, interviews, and design decisions**, not just syntax.

---

## Why the `context` Package Exists

Modern Go programs are concurrent and layered:

* HTTP handlers
* service layers
* database calls
* RPC calls
* goroutines

We need a **standard way** to:

* cancel work
* enforce timeouts
* pass request-related metadata

The `context` package solves this.

---

## Core Responsibilities of Context

Context is designed for **three things only**:

1. **Cancellation**
2. **Deadlines / Timeouts**
3. **Request-scoped metadata**

Anything beyond this is usually misuse.

---

## Cancellation (Most Important Feature)

Cancellation allows one signal to stop:

* many goroutines
* many function calls
* entire pipelines

### Key idea

> When a context is canceled, **all functions and goroutines using it should stop immediately**.

This prevents:

* goroutine leaks
* wasted CPU
* hanging requests

---

## Deadlines and Timeouts

Contexts can carry time constraints:

* A **deadline** = absolute time
* A **timeout** = duration from now

When the time expires:

* `ctx.Done()` is closed
* all dependent work should stop

---

## Request-Scoped Data (Controversial Part)

The controversial feature of context is:

```go
context.WithValue(ctx, key, value)
```

This allows attaching arbitrary data to a context.

The official Go guideline says:

> Use context values only for request-scoped data that transits processes and API boundaries, not for passing optional parameters to functions.

---

## What Is Request-Scoped Data?

Request-scoped data:

* is created when a request starts
* lives only for the lifetime of that request
* flows through many layers
* is immutable

Examples:

* Request ID
* Trace ID
* User ID
* Authorization token

---

## Why Context Values Are Controversial

Problems with `WithValue`:

* not type-safe (`interface{}`)
* easy to misuse
* can become hidden global state

Because of this, **clear rules are required**.

---

## Heuristics for Storing Data in Context

These are **guidelines**, not strict laws.

### 1. Data should cross process or API boundaries

If data never leaves your service, it usually does **not** belong in context.

---

### 2. Data should be immutable

Mutable data implies shared state, which breaks context’s purpose.

---

### 3. Data should be simple

Prefer:

* strings
* integers
* byte slices

Avoid:

* large structs
* database connections
* complex objects

---

### 4. Store data, not behavior

Context should contain **information**, not objects with methods.

Logic belongs in functions, not in context values.

---

### 5. Context should decorate behavior, not drive it

Context may:

* add logging metadata
* add tracing information

Context should **not**:

* change algorithms
* act as optional parameters

If program logic changes based on context values, it is likely misused.

---

## Examples of Common Context Data

| Data                | Appropriate for Context | Reason             |
| ------------------- | ----------------------- | ------------------ |
| Request ID          | Yes                     | Logging, tracing   |
| User ID             | Yes                     | Identity metadata  |
| URL                 | Sometimes               | Depends on usage   |
| Authorization token | Debatable               | Team decision      |
| DB connection       | No                      | Large, mutable     |
| Config values       | No                      | Not request-scoped |

---

## Invisible Dependency Trade-off

Using context creates an **implicit dependency**:

* functions do not declare exactly what data they need

Alternative:

* explicit parameters (more verbose, more explicit)

There is no universal right answer — teams must decide case by case.

---

## Best Practices

* Always pass context as the **first parameter**
* Never store business logic or mutable state in context
* Use unique, unexported key types
* Cancel contexts as soon as work is no longer needed
* Prefer `context.Context` over custom `done` channels in production code

---

## Relationship to Go Concurrency Patterns

Context replaces or simplifies:

* `done` channels
* or-channels
* manual cancellation wiring

It integrates naturally with:

* pipelines
* fan-out / fan-in
* worker pools

---

## Final Takeaway

> Use `context.Context` primarily for cancellation and deadlines. Use context values sparingly and only for request-scoped metadata.

Even if you avoid `WithValue`, **do not avoid context** — it is a core part of modern Go.

---

## One-Line Interview Summary

> The context package provides a standard way to propagate cancellation, deadlines, and request-scoped metadata across API boundaries in concurrent Go programs.
