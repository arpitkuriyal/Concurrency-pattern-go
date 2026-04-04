package main

import (
	"fmt"
	"time"
)

// Job represents a unit of work that needs to be processed.
// In real systems, this could be an email task, payment, etc.
type Job struct {
	ID   int
	Data string
}

// worker is a goroutine that continuously listens for jobs from the channel.
//
// WHY channel?
// Channels ensure safe communication between goroutines
// They transfer ownership of data (confinement)
//
// jobs <-chan Job means:
// This function can ONLY READ from the channel (receive-only)
// Prevents accidental writes
func worker(id int, jobs <-chan Job) {

	for job := range jobs {

		// Each worker processes its own job independently
		// No shared state, No race conditions
		fmt.Printf("Worker %d processing job %d: %s\n", id, job.ID, job.Data)

		// Simulate real work (API call, DB write, etc.)
		time.Sleep(1 * time.Second)
	}
}

// StartWorkers starts multiple worker goroutines.
// WHY multiple workers?
// Enables parallel processing
// Improves throughput (faster job handling)
//
// Each worker listens on the SAME channel:
// Go scheduler distributes jobs automatically
// This pattern is called "worker pool"
func StartWorkers(numWorkers int, jobs <-chan Job) {
	for i := 1; i <= numWorkers; i++ {
		go worker(i, jobs) // start worker as a goroutine
	}
}

// DispatchJobs sends jobs into the channel.
//
// jobs chan<- Job means:
// This function can ONLY WRITE to the channel (send-only)
// Prevents accidental reads (clear responsibility) (only one at a time get job)
//
// WHY close(jobs)?
// Signals workers: "no more jobs are coming"
// Without closing, workers would wait forever (goroutine leak)
func DispatchJobs(jobs chan<- Job, totalJobs int) {
	for j := 1; j <= totalJobs; j++ {

		// Sending job to channel:
		// Transfers ownership to a worker
		// Main goroutine should NOT touch it after sending
		jobs <- Job{
			ID:   j,
			Data: "Send Email",
		}
	}

	// Important: close channel after sending all jobs
	// Allows workers to exit gracefully
	close(jobs)
}

// func main() {

// 	// Create a channel to send jobs
// 	// This acts as a queue between producer (main) and consumers (workers)
// 	jobs := make(chan Job)

// 	// Start 3 worker goroutines
// 	// They will listen for jobs from the channel
// 	StartWorkers(3, jobs)

// 	// Send 5 jobs into the channel
// 	// Jobs will be picked up by available workers
// 	DispatchJobs(jobs, 5)

// 	// Wait for workers to finish processing
// 	// We use sleep here just to keep program alive
// 	// (Better approach = sync.WaitGroup, but this is simple)
// 	time.Sleep(4 * time.Second)
// }
