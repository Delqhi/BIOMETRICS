package benchmarks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"biometrics-cli/pkg/delegation"
)

// BenchmarkResult holds benchmark metrics
type BenchmarkResult struct {
	Name        string  `json:"name"`
	Iterations  int     `json:"iterations"`
	NsPerOp     int64   `json:"ns_per_op"`
	AllocsPerOp int64   `json:"allocs_per_op"`
	BytesPerOp  int64   `json:"bytes_per_op"`
	Throughput  float64 `json:"throughput_ops_sec"`
}

// BenchmarkDelegationEngine benchmarks the delegation engine
func BenchmarkDelegationEngine(b *testing.B) {
	router := delegation.NewDelegationRouter()
	engine := delegation.NewWorkerPool(10, router)
	defer engine.Shutdown()
	
	ctx := context.Background()

	task := delegation.NewTask("bench-task-1", delegation.TaskTypeCode, delegation.PriorityHigh, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Submit(task)
		resultChan := engine.Results()
		select {
		case result := <-resultChan:
			if !result.Success {
				b.Fatal(result.Error)
			}
		case <-time.After(5 * time.Second):
			b.Fatal("timeout")
		}
	}
}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Delegate(ctx, task)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTaskCreation benchmarks task creation performance
func BenchmarkTaskCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := delegation.NewTask(
			fmt.Sprintf("task-%d", i),
			delegation.TaskTypeCode,
			delegation.PriorityNormal,
			nil,
		)
		_ = task
	}
}
	_ = task
}

// BenchmarkTaskSerialization benchmarks JSON serialization
func BenchmarkTaskSerialization(b *testing.B) {
	task := delegation.NewTask("test-task", delegation.TaskTypeCode, delegation.PriorityHigh, nil)
	task.SetContext("key1", "value1")
	task.SetContext("key2", 123)
	task.SetContext("key3", true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, err := json.Marshal(task)
		if err != nil {
			b.Fatal(err)
		}
		_ = data
	}
}

// BenchmarkTaskDeserialization benchmarks JSON deserialization
func BenchmarkTaskDeserialization(b *testing.B) {
	task := &delegation.Task{
		ID:          "test-task",
		Title:       "Test Task",
		Description: "Testing deserialization performance",
		Priority:    delegation.PriorityHigh,
		Metadata: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
			"key3": true,
		},
	}

	data, _ := json.Marshal(task)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var t delegation.Task
		err := json.Unmarshal(data, &t)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConcurrentTasks benchmarks concurrent task processing
func BenchmarkConcurrentTasks(b *testing.B) {
	engine := delegation.NewEngine()
	ctx := context.Background()
	concurrency := runtime.NumCPU()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			task := &delegation.Task{
				ID:          fmt.Sprintf("concurrent-task-%d", i),
				Title:       fmt.Sprintf("Concurrent Task %d", i),
				Description: "Testing concurrent processing",
				Priority:    delegation.PriorityMedium,
			}
			engine.Delegate(ctx, task)
			i++
		}
	})
}

// BenchmarkQueueOperations benchmarks queue operations
func BenchmarkQueueOperations(b *testing.B) {
	queue := delegation.NewTaskQueue()
	ctx := context.Background()

	b.Run("Enqueue", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			task := &delegation.Task{
				ID:          fmt.Sprintf("queue-task-%d", i),
				Title:       fmt.Sprintf("Queue Task %d", i),
				Description: "Queue benchmark",
				Priority:    delegation.PriorityLow,
			}
			queue.Enqueue(ctx, task)
		}
	})

	b.Run("Dequeue", func(b *testing.B) {
		// Pre-populate queue
		for i := 0; i < b.N; i++ {
			task := &delegation.Task{
				ID:          fmt.Sprintf("dequeue-task-%d", i),
				Title:       fmt.Sprintf("Dequeue Task %d", i),
				Description: "Dequeue benchmark",
				Priority:    delegation.PriorityLow,
			}
			queue.Enqueue(ctx, task)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = queue.Dequeue(ctx)
		}
	})
}

// BenchmarkAggregator benchmarks result aggregation
func BenchmarkAggregator(b *testing.B) {
	agg := delegation.NewResultAggregator()
	ctx := context.Background()

	results := make([]*delegation.TaskResult, 10)
	for i := range results {
		results[i] = &delegation.TaskResult{
			TaskID:      fmt.Sprintf("result-task-%d", i),
			Status:      delegation.StatusCompleted,
			Output:      fmt.Sprintf("Output %d", i),
			CompletedAt: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, result := range results {
			agg.Add(result)
		}
		_ = agg.Aggregate(ctx)
	}
}

// BenchmarkMemoryAllocation benchmarks memory usage
func BenchmarkMemoryAllocation(b *testing.B) {
	b.ReportAllocs()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := &delegation.Task{
			ID:          fmt.Sprintf("mem-task-%d", i),
			Title:       fmt.Sprintf("Memory Task %d", i),
			Description: "Memory allocation benchmark",
			Priority:    delegation.PriorityMedium,
			Metadata:    make(map[string]interface{}),
		}

		// Add metadata to trigger allocations
		for j := 0; j < 10; j++ {
			task.Metadata[fmt.Sprintf("key-%d", j)] = fmt.Sprintf("value-%d", j)
		}
	}
}

// BenchmarkHTTPEndpoints benchmarks HTTP endpoint performance
func BenchmarkHTTPEndpoints(b *testing.B) {
	// Create test server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task := &delegation.Task{
			ID:          "http-task",
			Title:       "HTTP Task",
			Description: "HTTP benchmark",
			Priority:    delegation.PriorityMedium,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.Get(server.URL)
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

// BenchmarkFanOutPattern benchmarks fan-out delegation pattern
func BenchmarkFanOutPattern(b *testing.B) {
	engine := delegation.NewEngine()
	ctx := context.Background()

	parentTask := &delegation.Task{
		ID:          "parent-task",
		Title:       "Parent Task",
		Description: "Fan-out benchmark",
		Priority:    delegation.PriorityHigh,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create 10 subtasks
		var wg sync.WaitGroup
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				subtask := &delegation.Task{
					ID:          fmt.Sprintf("subtask-%d", idx),
					Title:       fmt.Sprintf("Subtask %d", idx),
					Description: "Fan-out subtask",
					Priority:    delegation.PriorityMedium,
					ParentID:    parentTask.ID,
				}
				engine.Delegate(ctx, subtask)
			}(j)
		}
		wg.Wait()
	}
}

// BenchmarkFanInPattern benchmarks fan-in aggregation pattern
func BenchmarkFanInPattern(b *testing.B) {
	agg := delegation.NewResultAggregator()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		// Simulate 10 parallel tasks completing
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				result := &delegation.TaskResult{
					TaskID:      fmt.Sprintf("fanin-task-%d", idx),
					Status:      delegation.StatusCompleted,
					Output:      fmt.Sprintf("Result %d", idx),
					CompletedAt: time.Now(),
				}
				agg.Add(result)
			}(j)
		}
		wg.Wait()

		_ = agg.Aggregate(ctx)
	}
}

// BenchmarkChainPattern benchmarks chain delegation pattern
func BenchmarkChainPattern(b *testing.B) {
	engine := delegation.NewEngine()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create chain of 5 tasks
		var prevTask *delegation.Task
		for j := 0; j < 5; j++ {
			task := &delegation.Task{
				ID:          fmt.Sprintf("chain-task-%d", j),
				Title:       fmt.Sprintf("Chain Task %d", j),
				Description: "Chain delegation",
				Priority:    delegation.PriorityMedium,
			}

			if prevTask != nil {
				task.ParentID = prevTask.ID
			}

			engine.Delegate(ctx, task)
			prevTask = task
		}
	}
}

// BenchmarkRetryPattern benchmarks retry logic performance
func BenchmarkRetryPattern(b *testing.B) {
	engine := delegation.NewEngine()
	ctx := context.Background()

	task := &delegation.Task{
		ID:          "retry-task",
		Title:       "Retry Task",
		Description: "Retry benchmark",
		Priority:    delegation.PriorityHigh,
		RetryConfig: &delegation.RetryConfig{
			MaxRetries: 3,
			Delay:      100 * time.Millisecond,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Delegate(ctx, task)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSystemResources benchmarks system resource usage
func BenchmarkSystemResources(b *testing.B) {
	b.Run("CPU", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// CPU-intensive operation
			sum := 0
			for j := 0; j < 1000; j++ {
				sum += j * j
			}
			_ = sum
		}
	})

	b.Run("Memory", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			data := make([]byte, 1024) // 1KB allocation
			for j := range data {
				data[j] = byte(j % 256)
			}
			_ = data
		}
	})

	b.Run("Goroutines", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for j := 0; j < 100; j++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					time.Sleep(1 * time.Millisecond)
				}()
			}
			wg.Wait()
		}
	})
}

// BenchmarkFileOperations benchmarks file I/O performance
func BenchmarkFileOperations(b *testing.B) {
	tmpDir := b.TempDir()

	b.Run("WriteFile", func(b *testing.B) {
		data := []byte("Benchmark test data for file write operations")

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			path := filepath.Join(tmpDir, fmt.Sprintf("bench-write-%d.txt", i))
			err := os.WriteFile(path, data, 0644)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ReadFile", func(b *testing.B) {
		// Create test file
		data := []byte("Benchmark test data for file read operations")
		path := filepath.Join(tmpDir, "bench-read.txt")
		os.WriteFile(path, data, 0644)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := os.ReadFile(path)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkCommandExecution benchmarks command execution
func BenchmarkCommandExecution(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("echo", "test")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// GenerateBenchmarkReport generates a comprehensive benchmark report
func GenerateBenchmarkReport(results []BenchmarkResult) error {
	report := struct {
		Timestamp     string            `json:"timestamp"`
		GoVersion     string            `json:"go_version"`
		NumCPU        int               `json:"num_cpu"`
		Benchmarks    []BenchmarkResult `json:"benchmarks"`
		TotalOps      int               `json:"total_ops"`
		AvgThroughput float64           `json:"avg_throughput"`
	}{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		GoVersion:  runtime.Version(),
		NumCPU:     runtime.NumCPU(),
		Benchmarks: results,
	}

	// Calculate totals
	for _, r := range results {
		report.TotalOps += r.Iterations
		report.AvgThroughput += r.Throughput
	}

	if len(results) > 0 {
		report.AvgThroughput /= float64(len(results))
	}

	// Write report
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("benchmarks/report.json", data, 0644)
}
