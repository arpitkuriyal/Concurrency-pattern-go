# Confinement (Go Concurrency Pattern)

## What problem does this solve?

In concurrent programs, bugs happen when:
- Multiple goroutines
- Access the same data
- At the same time

This leads to race conditions and hard-to-debug issues.

---

## What is Confinement?

**Confinement is a concurrency design technique.**

> Confinement means a piece of data is accessed by only one goroutine.

If only one goroutine can touch the data:
- The program is automatically safe
- No locks or synchronization are needed

---

## Important Notes

- Confinement is NOT a Go keyword
- Confinement is NOT enforced by the language
- It is a way of **designing concurrent code**

---

## Why Confinement Works

Concurrency bugs happen because of **shared mutable data**.

Confinement removes sharing entirely:
- No race conditions
- No deadlocks
- No mutexes
- No atomics

---

## Immutable Data vs Confinement

### Immutable Data
- Data is never modified
- All goroutines can read it
- Safe because nothing changes

### Confinement
- Data can change
- But only one goroutine can access it

Both approaches reduce concurrency problems.

---

## Types of Confinement

### 1. Ad-hoc Confinement
- Safety is based on convention
- Developers promise not to misuse data
- Compiler does not help

❌ Easy to break  
❌ Not recommended for large codebases

---

### 2. Lexical Confinement (Recommended)
- Uses scope to enforce safety
- Compiler prevents incorrect access
- Makes misuse impossible

This is the preferred approach in Go.

---

## Confinement Using Channels

A common pattern is the **channel owner pattern**:
- One goroutine owns the channel
- Other goroutines get read-only or write-only access

This prevents:
- Writing from the wrong goroutine
- Closing the channel incorrectly

---

## Confinement with Non-Thread-Safe Data

Even unsafe data structures can be used safely if:
- Each goroutine gets its own data
- No shared memory exists
- Data is limited by scope

No synchronization is needed in this case.

---

## When Should You Use Confinement?

Always ask:

> Can this data belong to only one goroutine?

- If YES → use confinement
- If NO → use mutexes or other synchronization

---

## Mental Rule (IMPORTANT)

### Owner Rule
> Every piece of data should have exactly one owner goroutine.

If you cannot identify the owner, the design is unsafe.

---

## One-Line Summary

> If only one goroutine can access the data, the data is safe.

This is confinement.

---

## Key Takeaway

- Confinement is a design technique
- It reduces complexity and bugs
- It should be your first choice in Go concurrency
