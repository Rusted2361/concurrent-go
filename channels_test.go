package main

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func fastCounter(_ string, content string) int {
	return len(strings.Fields(content))
}

func makeFiles(n int, content string) map[string]string {
	m := make(map[string]string, n)
	for i := 1; i <= n; i++ {
		m[fmt.Sprintf("file%d.txt", i)] = content
	}
	return m
}

func TestProcessFiles_Deterministic(t *testing.T) {
	const testContent = "this is a sample text file that has eleven words here"
	files := makeFiles(20, testContent)

	got := ProcessFiles(files, 4, fastCounter)
	want := 20 * 11
	if got != want {
		t.Fatalf("ProcessFiles() = %d, want %d", got, want)
	}
}

func TestProcessFiles_HandlesNumWorkersZero(t *testing.T) {
	const testContent = "a b c"
	files := makeFiles(10, testContent)

	got := ProcessFiles(files, 0, fastCounter) // should default to 1
	want := 10 * 3
	if got != want {
		t.Fatalf("ProcessFiles() = %d, want %d", got, want)
	}
}

func BenchmarkProcessFiles_WorkerPool(b *testing.B) {
	const testContent = "this is a sample text file that has eleven words here"
	files := makeFiles(200, testContent)
	workers := runtime.NumCPU() * 2

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ProcessFiles(files, workers, fastCounter)
	}
}

// func BenchmarkProcessFiles_SingleWorker(b *testing.B) {
// 	const testContent = "this is a sample text file that has eleven words here"
// 	files := makeFiles(200, testContent)

// 	b.ReportAllocs()
// 	for i := 0; i < b.N; i++ {
// 		_ = ProcessFiles(files, 1, fastCounter)
// 	}
// }

func BenchmarkProcessFiles_PerGoroutine(b *testing.B) {
	const testContent = "this is a sample text file that has eleven words here"
	files := makeFiles(200, testContent)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ProcessFilesPerGoroutine(files, fastCounter)
	}
}
