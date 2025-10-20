package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Counter func(filename, content string) int

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

	// total := 0
	// var mu sync.Mutex
	numWorkers := runtime.NumCPU() * 2
	start := time.Now()
	total := ProcessFiles(files, numWorkers, countWords)

	fmt.Printf("Worker total: %d\n", total)
	fmt.Printf("Time taken: %v\n", time.Since(start))

	start = time.Now()
	total = ProcessFilesPerGoroutine(files, countWords)

	fmt.Printf("Worker total: %d\n", total)
	fmt.Printf("Time taken: %v\n", time.Since(start))

}

func ProcessFiles(files map[string]string, numWorkers int, counter Counter) int {
	if numWorkers <= 0 {
		numWorkers = 1
	}

	jobs := make(chan Job)
	results := make(chan int)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- counter(job.filename, job.content)
				// mu.Lock()
				// total += count
				// mu.Unlock()
			}
		}()
	}

	go func() {
		for filename, content := range files {
			jobs <- Job{filename, content}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	for c := range results {
		total += c
	}
	return total
}

func ProcessFilesPerGoroutine(files map[string]string, counter Counter) int {
	results := make(chan int)
	var wg sync.WaitGroup

	wg.Add(len(files))
	for filename, content := range files {
		fn := filename
		ct := content
		go func() {
			defer wg.Done()
			results <- counter(fn, ct)
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	for c := range results {
		total += c
	}
	return total
}
