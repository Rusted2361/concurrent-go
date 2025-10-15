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

	jobs := make(chan Job)
	results := make(chan int)
	var wg sync.WaitGroup
	// total := 0
	// var mu sync.Mutex

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

	// Send all jobs
	go func() {
		for filename, content := range files {
			jobs <- Job{filename, content}
		}
		close(jobs)
	}()

	// Close results once all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	for c := range results {
		total += c
	}
	fmt.Printf("Worker total: %d\n", total)
	fmt.Printf("Time taken: %v\n", time.Since(start))

}
