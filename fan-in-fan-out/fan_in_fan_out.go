package main

import (
	"fmt"
	"sync"
	"time"
)

// STEP 1: Generator (creates input stream)
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// STEP 2: Worker (FAN-OUT)
// Multiple workers will read from the same input channel
func worker(id int, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			fmt.Printf("worker %d processing %d\n", id, n)
			time.Sleep(500 * time.Millisecond) // simulate work
			out <- n * n
		}
	}()
	return out
}

// STEP 3: Fan-In (merge multiple worker outputs)
func fanIn(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Function to copy values from one channel to out
	output := func(ch <-chan int) {
		defer wg.Done()
		for v := range ch {
			out <- v
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go output(ch)
	}

	// Close out once all inputs are drained
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// STEP 4: Main
func main() {
	// Create input
	input := generator(1, 2, 3, 4, 5, 6)

	// FAN-OUT: multiple workers reading from same input
	w1 := worker(1, input)
	w2 := worker(2, input)
	w3 := worker(3, input)

	// FAN-IN: merge worker outputs
	results := fanIn(w1, w2, w3)

	// Consume final output
	for result := range results {
		fmt.Println("result:", result)
	}
}
