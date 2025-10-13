package main

import (
	"fmt"
	"strings"
	"time"
)

func countWords(filename string, content string) int {
	time.Sleep(100 * time.Millisecond) // Simulate I/O
	return len(strings.Fields(content))
}

// defined a custom type named Job
type Job struct {
	Filename string
	Content  string
}

func main() {
	// Each file has exactly 11 words
	const testContent = "this is a sample text file that has eleven words here"

	files := make(map[string]string)
	for i := 1; i <= 20; i++ {
		files[fmt.Sprintf("file%d.txt", i)] = testContent
	}

	// jobs channel is responsible for sending Job structs to the worker goroutine
	// done channel is responsible for signaling the main goroutine that the worker has finished processing all jobs
	jobs := make(chan Job)
	done := make(chan bool)

	total := 0

	// worker goroutine that will perform the work
	go func() {
		for job := range jobs {
			count := countWords(job.Filename, job.Content)
			total += count
		}
		done <- true
	}()

	// main goroutine that will send the jobs to the worker goroutine
	start := time.Now()
	for filename, content := range files {
		jobs <- Job{filename, content}
	}
	// close the jobs channel to signal to the worker goroutine that there are no more jobs
	close(jobs)
	// wait for the worker goroutine to finish processing all jobs by blocking on the done channel
	<-done
	fmt.Printf("Total words: %d, Time taken: %v\n", total, time.Since(start))
}
