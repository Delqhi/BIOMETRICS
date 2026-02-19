package benchmarks

import "testing"

func BenchmarkAll(b *testing.B) {
	b.Run("DelegationRouter", BenchmarkDelegationRouter)
	b.Run("TaskCreation", BenchmarkTaskCreation)
	b.Run("TaskSerialization", BenchmarkTaskSerialization)
	b.Run("TaskDeserialization", BenchmarkTaskDeserialization)
	b.Run("ConcurrentTasks", BenchmarkConcurrentTasks)
	b.Run("QueueOperations", BenchmarkQueueOperations)
	b.Run("Aggregator", BenchmarkAggregator)
	b.Run("MemoryAllocation", BenchmarkMemoryAllocation)
	b.Run("HTTPEndpoints", BenchmarkHTTPEndpoints)
	b.Run("FanOutPattern", BenchmarkFanOutPattern)
	b.Run("FanInPattern", BenchmarkFanInPattern)
	b.Run("ChainPattern", BenchmarkChainPattern)
	b.Run("CircuitBreaker", BenchmarkCircuitBreaker)
	b.Run("SystemResources", BenchmarkSystemResources)
	b.Run("FileOperations", BenchmarkFileOperations)
	b.Run("CommandExecution", BenchmarkCommandExecution)
}
