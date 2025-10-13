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

func main() {
	// Each file has exactly 11 words
	const testContent = "this is a sample text file that has eleven words here"

	files := make(map[string]string)
	for i := 1; i <= 20; i++ {
		files[fmt.Sprintf("file%d.txt", i)] = testContent
	}

	start := time.Now()
	total := 0
	for filename, content := range files {
		count := countWords(filename, content)
		total += count
	}
	fmt.Printf("Total words: %d, Time taken: %v\n", total, time.Since(start))
}
