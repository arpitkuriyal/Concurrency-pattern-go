package main

import (
	"fmt"
	"time"
)

func main() {

	// =====================================================
	// EXAMPLE 1: Goroutine blocks forever (RECEIVER LEAK)
	// =====================================================
	fmt.Println("Example 1: Receiver blocks forever")

	ch1 := make(chan int)

	go func() {
		fmt.Println("[Ex1] Goroutine started")
		<-ch1                                // BLOCKS forever (no sender)
		fmt.Println("[Ex1] Goroutine ended") // NEVER runs
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("[Ex1] Main continues")

	// =====================================================
	// EXAMPLE 2: Goroutine blocks forever (SENDER LEAK)
	// =====================================================
	fmt.Println("Example 2: Sender blocks forever")

	ch2 := make(chan int)

	go func() {
		fmt.Println("[Ex2] Producer started")
		ch2 <- 10                           // BLOCKS forever (no receiver)
		fmt.Println("[Ex2] Producer ended") // NEVER runs
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("[Ex2] Main continues")

	// =====================================================
	// EXAMPLE 3: FIX receiver leak using done channel
	// =====================================================
	fmt.Println("Example 3: Receiver with done (FIXED)")

	ch3 := make(chan int)
	done3 := make(chan struct{})

	go func() {
		fmt.Println("[Ex3] Goroutine started")
		select {
		case <-ch3:
			fmt.Println("[Ex3] Received value")
		case <-done3:
			fmt.Println("[Ex3] Cancel signal received")
			return
		}
		fmt.Println("[Ex3] Goroutine ended")
	}()

	time.Sleep(1 * time.Second)
	close(done3) // cancel goroutine
	time.Sleep(500 * time.Millisecond)
	fmt.Println("[Ex3] Main continues")

	// =====================================================
	// EXAMPLE 4: FIX sender leak using done channel
	// =====================================================
	fmt.Println("Example 4: Sender with done (FIXED)")

	ch4 := make(chan int)
	done4 := make(chan struct{})

	go func() {
		fmt.Println("[Ex4] Producer started")
		select {
		case ch4 <- 42:
			fmt.Println("[Ex4] Value sent")
		case <-done4:
			fmt.Println("[Ex4] Producer canceled")
			return
		}
		fmt.Println("[Ex4] Producer ended")
	}()

	time.Sleep(1 * time.Second)
	close(done4)
	time.Sleep(500 * time.Millisecond)
	fmt.Println("[Ex4] Main exiting")

	// =====================================================
	// END
	// =====================================================
}
