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

// BenchmarkMetrics holds benchmark metrics
type BenchmarkMetrics struct {
	Name        string  `json:"name"`
	Iterations  int     `json:"iterations"`
	NsPerOp     int64   `json:"ns_per_op"`
	AllocsPerOp int64   `json:"allocs_per_op"`
	BytesPerOp  int64   `json:"bytes_per_op"`
	Throughput  float64 `json:"throughput_ops_sec"`
}

// BenchmarkDelegationRouter benchmarks the delegation router
func BenchmarkDelegationRouter(b *testing.B) {
	router := delegation.NewDelegationRouter()

	router.RegisterAgent(&delegation.AgentCapability{
		Name:         "sisyphus",
		AgentID:      "agent-001",
		Capabilities: []string{"code", "testing"},
		Load:         0,
		Healthy:      true,
	})

	router.RegisterAgent(&delegation.AgentCapability{
		Name:         "prometheus",
		AgentID:      "agent-002",
		Capabilities: []string{"architecture", "planning"},
		Load:         0,
		Healthy:      true,
	})

	task := delegation.NewTask("bench-task", delegation.TaskTypeCode, delegation.PriorityHigh, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := router.Route(task)
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

// BenchmarkTaskSerialization benchmarks JSON serialization
func BenchmarkTaskSerialization(b *testing.B) {
	task := delegation.NewTask("test-task", delegation.TaskTypeCode, delegation.PriorityHigh, map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	})

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
	task := delegation.NewTask("test-task", delegation.TaskTypeCode, delegation.PriorityHigh, map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	})

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
	router := delegation.NewDelegationRouter()
	engine := delegation.NewWorkerPool(10, router)
	defer engine.Shutdown()

	router.RegisterAgent(&delegation.AgentCapability{
		Name:         "sisyphus",
		AgentID:      "agent-001",
		Capabilities: []string{"code"},
		Load:         0,
		Healthy:      true,
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			task := delegation.NewTask(
				fmt.Sprintf("concurrent-task-%d", i),
				delegation.TaskTypeCode,
				delegation.PriorityNormal,
				nil,
			)
			engine.Submit(task)
			i++
		}
	})
}

// BenchmarkQueueOperations benchmarks priority queue operations
func BenchmarkQueueOperations(b *testing.B) {
	b.Run("Enqueue", func(b *testing.B) {
		pq := delegation.NewPriorityQueue()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			task := delegation.NewTask(
				fmt.Sprintf("queue-task-%d", i),
				delegation.TaskTypeCode,
				delegation.PriorityLow,
				nil,
			)
			pq.Enqueue(task)
		}
	})

	b.Run("Dequeue", func(b *testing.B) {
		pq := delegation.NewPriorityQueue()

		for i := 0; i < b.N; i++ {
			task := delegation.NewTask(
				fmt.Sprintf("dequeue-task-%d", i),
				delegation.TaskTypeCode,
				delegation.PriorityLow,
				nil,
			)
			pq.Enqueue(task)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = pq.Dequeue()
		}
	})
}

// BenchmarkAggregator benchmarks result aggregation
func BenchmarkAggregator(b *testing.B) {
	config := delegation.AggregatorConfig{
		Strategy:   delegation.MergeStrategyConcat,
		Timeout:    5 * time.Second,
		MinResults: 1,
	}
	agg := delegation.NewResultAggregator(config)

	results := make([]*delegation.TaskResult, 10)
	for i := range results {
		results[i] = &delegation.TaskResult{
			TaskID:   fmt.Sprintf("result-task-%d", i),
			Success:  true,
			Data:     fmt.Sprintf("Output %d", i),
			Duration: 100 * time.Millisecond,
			AgentID:  "agent-001",
		}
	}

	agg.SetTotalTasks("batch-001", len(results))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, result := range results {
			agg.Collect("batch-001", result)
		}
		_, _ = agg.Wait("batch-001", 5*time.Second)
	}
}

// BenchmarkMemoryAllocation benchmarks memory usage
func BenchmarkMemoryAllocation(b *testing.B) {
	b.ReportAllocs()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := delegation.NewTask(
			fmt.Sprintf("mem-task-%d", i),
			delegation.TaskTypeCode,
			delegation.PriorityNormal,
			make(map[string]interface{}),
		)

		for j := 0; j < 10; j++ {
			task.SetContext(fmt.Sprintf("key-%d", j), fmt.Sprintf("value-%d", j))
		}
	}
}

// BenchmarkHTTPEndpoints benchmarks HTTP endpoint performance
func BenchmarkHTTPEndpoints(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task := delegation.NewTask("http-task", delegation.TaskTypeCode, delegation.PriorityNormal, nil)

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
	router := delegation.NewDelegationRouter()
	engine := delegation.NewWorkerPool(10, router)
	defer engine.Shutdown()

	for i := 0; i < 5; i++ {
		router.RegisterAgent(&delegation.AgentCapability{
			Name:         fmt.Sprintf("agent-%d", i),
			AgentID:      fmt.Sprintf("agent-%03d", i),
			Capabilities: []string{"code"},
			Load:         0,
			Healthy:      true,
		})
	}

	parentTask := delegation.NewTask("parent-task", delegation.TaskTypeCode, delegation.PriorityHigh, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				subtask := delegation.NewTask(
					fmt.Sprintf("subtask-%d", idx),
					delegation.TaskTypeCode,
					delegation.PriorityNormal,
					nil,
				)
				engine.Submit(subtask)
			}(j)
		}
		wg.Wait()
		_ = parentTask
	}
}

// BenchmarkFanInPattern benchmarks fan-in aggregation pattern
func BenchmarkFanInPattern(b *testing.B) {
	config := delegation.AggregatorConfig{
		Strategy:   delegation.MergeStrategyConcat,
		Timeout:    5 * time.Second,
		MinResults: 1,
	}
	agg := delegation.NewResultAggregator(config)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				result := &delegation.TaskResult{
					TaskID:   fmt.Sprintf("fanin-task-%d", idx),
					Success:  true,
					Data:     fmt.Sprintf("Result %d", idx),
					Duration: 50 * time.Millisecond,
					AgentID:  "agent-001",
				}
				agg.Collect("batch-001", result)
			}(j)
		}
		wg.Wait()

		_, _ = agg.Wait("batch-001", 5*time.Second)
		_ = ctx
	}
}

// BenchmarkChainPattern benchmarks chain delegation pattern
func BenchmarkChainPattern(b *testing.B) {
	router := delegation.NewDelegationRouter()
	engine := delegation.NewWorkerPool(10, router)
	defer engine.Shutdown()

	router.RegisterAgent(&delegation.AgentCapability{
		Name:         "sisyphus",
		AgentID:      "agent-001",
		Capabilities: []string{"code", "testing", "deploy"},
		Load:         0,
		Healthy:      true,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var prevTask *delegation.Task
		for j := 0; j < 5; j++ {
			task := delegation.NewTask(
				fmt.Sprintf("chain-task-%d", j),
				delegation.TaskTypeCode,
				delegation.PriorityNormal,
				nil,
			)

			if prevTask != nil {
				task.SetContext("parent_id", prevTask.ID)
			}

			engine.Submit(task)
			prevTask = task
		}
	}
}

// BenchmarkCircuitBreaker benchmarks circuit breaker performance
func BenchmarkCircuitBreaker(b *testing.B) {
	cb := delegation.NewCircuitBreaker(3, 30*time.Second)

	b.Run("Success", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if !cb.CanExecute() {
				b.Fatal("Circuit should be closed")
			}
			cb.RecordSuccess()
		}
	})

	b.Run("Failure", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cb.RecordFailure()
			_ = cb.CanExecute()
		}
	})
}

// BenchmarkSystemResources benchmarks system resource usage
func BenchmarkSystemResources(b *testing.B) {
	b.Run("CPU", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
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
			data := make([]byte, 1024)
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
func GenerateBenchmarkReport(results []BenchmarkMetrics) error {
	report := struct {
		Timestamp     string             `json:"timestamp"`
		GoVersion     string             `json:"go_version"`
		NumCPU        int                `json:"num_cpu"`
		Benchmarks    []BenchmarkMetrics `json:"benchmarks"`
		TotalOps      int                `json:"total_ops"`
		AvgThroughput float64            `json:"avg_throughput"`
	}{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		GoVersion:  runtime.Version(),
		NumCPU:     runtime.NumCPU(),
		Benchmarks: results,
	}

	for _, r := range results {
		report.TotalOps += r.Iterations
		report.AvgThroughput += r.Throughput
	}

	if len(results) > 0 {
		report.AvgThroughput /= float64(len(results))
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("benchmarks/report.json", data, 0644)
}
