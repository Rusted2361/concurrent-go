package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// defined a custom type named Job
type Job struct {
	Filename string
	Content  string
}

func worker(jobs chan Job, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		count := len(strings.Fields(job.Content))
		time.Sleep(100 * time.Millisecond)
		results <- count

	}
}

func main() {
	// Each file has exactly 11 words
	const testContent = "this is a sample text file that has eleven words here"

	files := make(map[string]string)
	for i := 1; i <= 20; i++ {
		files[fmt.Sprintf("file%d.txt", i)] = testContent
	}

	numWorkers := 3
	numJobs := len(files)

	jobs := make(chan Job, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup
	// worker goroutine that will perform the work
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// main goroutine that will send the jobs to the worker goroutine
	start := time.Now()
	for filename, content := range files {
		jobs <- Job{filename, content}
	}
	// close the jobs channel to signal to the worker goroutine that there are no more jobs
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	for count := range results { // This will exit when results channel is closed
		total += count
	}
	fmt.Printf("Total words: %d, Time taken: %v\n", total, time.Since(start))
}
