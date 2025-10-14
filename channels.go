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
	var wg sync.WaitGroup
	total := 0

	var mu sync.Mutex

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				count := countWords(job.filename, job.content)
				mu.Lock()
				total += count
				mu.Unlock()
			}

		}()
	}

	start := time.Now()

	for filename, content := range files {
		jobs <- Job{filename, content}
	}
	close(jobs)
	wg.Wait()

	fmt.Printf("Worker total: %d\n", total)
	fmt.Printf("Time taken: %v\n", time.Since(start))
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	fmt.Println("pprof running on :6060")
	// 	http.ListenAndServe("localhost:6060", nil)
	// }()
	// wg.Wait()

}
