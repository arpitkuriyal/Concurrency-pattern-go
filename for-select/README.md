# for–select Pattern to Prevent Goroutine Leaks

This note explains **why `for + select` is needed** and how it prevents goroutine leaks, with simple examples.

---

## Why `for` alone is dangerous

### Receiver leak (no select)

```go
for {
    <-ch // blocks forever if no sender
}
```

* Goroutine waits on `<-ch`
* If nobody sends, execution **stops here forever**
* Goroutine can never exit → **leak**

---

### Sender leak (no select)

```go
for {
    ch <- value // blocks forever if no receiver
}
```

* Sending also blocks
* If nobody receives, goroutine is stuck forever
* Another **leak**

---

## What `select` does

`select` lets a goroutine **wait on multiple events at the same time**.

This allows the goroutine to say:

> "I will wait for work, but I can also stop if told to."

---

## `done` channel

* `done` is a **cancellation signal**
* Closing `done` tells goroutines to exit
* `done` does NOT block
* `done` **unblocks** blocked goroutines

---

## Receiver using `for + select` (SAFE)

```go
for {
    select {
    case v := <-ch:
        fmt.Println(v)
    case <-done:
        return
    }
}
```

### Why this works

* If `ch` has data → process it
* If `ch` blocks → goroutine waits
* If `done` is closed → goroutine exits cleanly

---

## Sender using `for + select` (SAFE)

```go
for {
    select {
    case ch <- value:
        fmt.Println("sent")
    case <-done:
        return
    }
}
```

### Why this works

* Send only happens when receiver is ready
* Cancellation is always checked
* No permanent blocking

---

## Blocking vs Leak

| Term     | Meaning                 |
| -------- | ----------------------- |
| Blocking | Goroutine is waiting    |
| Leak     | Goroutine waits forever |
| done     | Exit signal             |

---

## Golden Rule (IMPORTANT)

> Any goroutine running in a loop and touching channels **must use `for + select` with a done case**.

---

## One-line exam answer

> The `for–select` pattern prevents goroutine leaks by allowing cancellation through a done channel while waiting on blocking channel operations.

---

## Final takeaway

```text
for + channel = danger
for + select + done = safe
```
