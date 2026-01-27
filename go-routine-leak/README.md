# Goroutine Blocking and Leaks 

## What does **block** mean in Go?

**Blocking means a goroutine is waiting and cannot move forward.**

* The goroutine is **alive**
* It is **not executing code**
* It is **waiting for something to happen**

Blocking is **normal** in Go.

---

## Common reasons a goroutine blocks

### 1. Blocked while **receiving** from a channel

```go
value := <-ch
```

The goroutine blocks if:

* No value is sent on `ch`
* `ch` is `nil`
* Sender never sends

---

### 2. Blocked while **sending** to a channel

```go
ch <- value
```

The goroutine blocks if:

* Channel is unbuffered
* No goroutine is receiving
* Buffer is full

---

### 3. Blocked while using `range`

```go
for v := range ch {
}
```

The goroutine blocks until:

* The channel is **closed**

If the channel is never closed → blocks forever.

---

## Does blocking mean "stuck in memory"?

**Not exactly.**

| Term      | Meaning                      |
| --------- | ---------------------------- |
| Blocked   | Goroutine is waiting         |
| In memory | Goroutine still exists       |
| Leaked    | Goroutine is blocked forever |

---

## What is a goroutine leak?

A **goroutine leak** happens when:

```text
Blocked + No way to exit = Goroutine Leak
```

* Goroutine is waiting
* No signal to stop
* It never exits

This wastes memory and resources.

---

## What is the `done` channel?

* `done` is a **cancellation signal**
* It tells a goroutine: **"Stop working and exit"**
* Closing `done` broadcasts the stop signal

---

## Important point about `done`

`done` does NOT block a goroutine
`done` **unblocks** a goroutine

---

## How `done` prevents blocking forever

```go
select {
case v := <-ch:
    // do work
case <-done:
    return
}
```

* If `ch` never receives → goroutine would block
* If `done` is closed → goroutine wakes up and exits

---

## Simple textbook rule

> If a goroutine can block forever, it must listen to `done`.

---

## Do we always need `done`?

### No `done` needed (finishes naturally)

```go
go func() {
    fmt.Println("Hello")
}()
```

---

### `done` needed (may block forever)

```go
go func() {
    for {
        <-ch
    }
}()
```

---

## Final one-line summary (for exams)

> Blocking means a goroutine is waiting on an operation; if it has no cancellation path, it leaks. The `done` channel provides that cancellation path.
