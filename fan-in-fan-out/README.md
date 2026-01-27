# Fan-In / Fan-Out in Go

This document explains **Fan-Out** and **Fan-In** concurrency patterns in Go, starting from simple intuition and moving toward practical usage.

---

## What is Fan-Out?

**Fan-Out** means:

> One input channel → multiple worker goroutines

* Many goroutines read from the **same input channel**
* Each value is processed by **only one worker**
* Work is distributed automatically by Go’s scheduler

### Why Fan-Out?

* To use multiple CPU cores
* To speed up slow or expensive work
* To process items in parallel

---

## What is Fan-In?

**Fan-In** means:

> Multiple worker output channels → one output channel

* Results from many goroutines are merged
* The consumer reads from **one single channel**

### Why Fan-In?

* Simplifies consumption of results
* Keeps downstream pipeline stages clean
* Centralizes output handling

---

## Mental Picture

```
            worker 1 ──┐
input ───► worker 2 ──┼──► output
            worker 3 ──┘
```

---

## Fan-Out Example (Theory)

```go
// Multiple workers read from the same input channel
input := make(chan int)

go worker(input)
go worker(input)
go worker(input)
```

* All workers listen on `input`
* Each value sent to `input` goes to **exactly one worker**
* Order of processing is **not guaranteed**

---

## Fan-In Example (Theory)

```go
// Multiple channels merged into one
out := fanIn(worker1Out, worker2Out, worker3Out)
```

* `fanIn` waits for values from all worker channels
* Emits values on a single output channel
* Output channel closes when all workers finish

---

## Basic Worker Pattern

```go
func worker(in <-chan int, out chan<- int) {
    for v := range in {
        out <- v * v
    }
}
```

* Reads input
* Processes data
* Sends result

---

## Correct Fan-In Pattern (Important)

```go
func fanIn(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    output := func(c <-chan int) {
        defer wg.Done()
        for v := range c {
            out <- v
        }
    }

    wg.Add(len(channels))
    for _, ch := range channels {
        go output(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

### Why this works

* Uses `sync.WaitGroup`
* Ensures output channel is closed exactly once
* Prevents goroutine leaks

---

## Important Characteristics

### 1. Order is NOT guaranteed

* Faster workers finish earlier
* Output order may differ from input order

### 2. Work distribution is automatic

* Go runtime schedules receives
* No manual load balancing needed

### 3. Fan-Out and Fan-In are usually used together

* Fan-Out for parallel work
* Fan-In to collect results

---

## Cancellation (Best Practice)

In real applications:

* Add `done` channel or `context.Context`
* Workers must stop when canceled
* Fan-In must also respect cancellation

---

## When to Use Fan-In / Fan-Out

* CPU-bound tasks
* I/O-bound tasks
* Worker pools
* Parallel stages in pipelines

---

## When NOT to Use

* Very small workloads
* Order must be strictly preserved (without extra logic)
* Simple sequential logic

---

## One-Line Summary

> Fan-Out distributes work across multiple goroutines, and Fan-In merges their results into a single channel to enable safe and efficient parallel processing.
