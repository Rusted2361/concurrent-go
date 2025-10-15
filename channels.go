package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func countWords(filename string, content string) int {
	time.Sleep(100 * time.Millisecond)
	return len(strings.Fields(content))
}

type Job struct {
	filename string
	content  string
}

func main() {
	const testContent = "this is a sample text file that has eleven words here"

	files := make(map[string]string)
	for i := 1; i <= 20; i++ {
		files[fmt.Sprintf("file%d.txt", i)] = testContent
	}

	// TODO: Step 1 - Create individual input channels for each worker (fan-out)
	// Create a slice of job channels: workerChannels := make([]chan Job, 20)
	// Initialize each channel in the slice
	jobs := make(chan Job)

	// TODO: Step 2 - Create individual output channels for each worker (fan-in)
	// Create a slice of result channels: resultChannels := make([]chan int, 20)
	// Initialize each channel in the slice
	results := make(chan int)
	var wg sync.WaitGroup
	// total := 0
	// var mu sync.Mutex

	// TODO: Step 3 - Modify workers to use individual channels
	// Each worker should read from workerChannels[i] and write to resultChannels[i]
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- countWords(job.filename, job.content)
				// mu.Lock()
				// total += count
				// mu.Unlock()
			}

		}()
	}

	start := time.Now()

	// TODO: Step 4 - Implement fan-out logic
	// Distribute jobs from the files map to individual worker channels
	// Use round-robin or other distribution strategy
	// Close each worker channel when done
	// Send all jobs
	go func() {
		for filename, content := range files {
			jobs <- Job{filename, content}
		}
		close(jobs)
	}()

	// TODO: Step 5 - Implement fan-in logic
	// Merge all resultChannels into a single results channel
	// Create a goroutine for each result channel that forwards to results
	// Use WaitGroup to know when all results are collected
	// Close results once all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// TODO: Step 6 - This part stays the same
	// Collect from the merged results channel
	total := 0
	for c := range results {
		total += c
	}
	fmt.Printf("Worker total: %d\n", total)
	fmt.Printf("Time taken: %v\n", time.Since(start))

}
