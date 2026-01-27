package main

import (
	"fmt"
	"time"
)

func main() {

	// =====================================================
	// EXAMPLE 1: for-loop RECEIVER blocks forever (LEAK)
	// =====================================================
	fmt.Println("Example 1: for receiver without select (LEAK)")

	ch1 := make(chan int)

	go func() {
		fmt.Println("[Ex1] Goroutine started")
		for {
			<-ch1                               // BLOCKS forever (no sender)
			fmt.Println("[Ex1] Received value") // NEVER runs
		}
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("[Ex1] Main continues")

	// =====================================================
	// EXAMPLE 2: for-loop SENDER blocks forever (LEAK)
	// =====================================================
	fmt.Println("Example 2: for sender without select (LEAK)")

	ch2 := make(chan int)

	go func() {
		fmt.Println("[Ex2] Producer started")
		for {
			ch2 <- 10                       // BLOCKS forever (no receiver)
			fmt.Println("[Ex2] Sent value") // NEVER runs
		}
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("[Ex2] Main continues")

	// =====================================================
	// EXAMPLE 3: for-select RECEIVER with done (FIXED)
	// =====================================================
	fmt.Println("Example 3: for-select receiver with done (FIXED)")

	ch3 := make(chan int)
	done3 := make(chan struct{})

	go func() {
		fmt.Println("[Ex3] Goroutine started")
		for {
			select {
			case v := <-ch3:
				fmt.Println("[Ex3] Received:", v)
			case <-done3:
				fmt.Println("[Ex3] Cancel signal received")
				return
			}
		}
	}()

	time.Sleep(1 * time.Second)
	close(done3)
	time.Sleep(500 * time.Millisecond)
	fmt.Println("[Ex3] Main continues")

	// =====================================================
	// EXAMPLE 4: for-select SENDER with done (FIXED)
	// =====================================================
	fmt.Println("Example 4: for-select sender with done (FIXED)")

	ch4 := make(chan int)
	done4 := make(chan struct{})

	go func() {
		fmt.Println("[Ex4] Producer started")
		for {
			select {
			case ch4 <- 42:
				fmt.Println("[Ex4] Sent value")
			case <-done4:
				fmt.Println("[Ex4] Producer canceled")
				return
			}
		}
	}()

	time.Sleep(1 * time.Second)
	close(done4)
	time.Sleep(500 * time.Millisecond)
	fmt.Println("[Ex4] Main exiting")

	// =====================================================
	// END
	// =====================================================
}
