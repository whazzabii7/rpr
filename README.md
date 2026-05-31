# RPR (Remote Process Routing)

High-Performance In-Memory Communication Matrix for Go.

## Description

RPR is a minimalist, highly efficient, and type-safe in-memory communication matrix for the Go runtime environment. The framework was specifically developed for embedded, event-driven micro-engines and systems designed according to the Actor model, where nanosecond-level latencies and the absolute avoidance of Garbage Collection (GC) overhead are critical system requirements.

Through the consistent use of internal memory recycling via sync.Pool and the elimination of compute-intensive runtime reflection, RPR guarantees allocation-free data transmission between isolated system components.

## AIM-Goal

The core goals of the framework are strictly aligned with the requirements of high-availability systems (SLA-driven Enterprise Architectures):

* Zero-Allocation Lifecycle: Complete elimination of heap allocations during the entire request-response cycle to ensure deterministic latencies.
* Compile-Time Type Safety: Physical isolation of module pipelines through the use of Go generics and type constraints (RequestConstraint). Erroneous data or command addressing is caught at compile time rather than at runtime.
* Resilience & Fault Isolation: Safeguarding long-running worker goroutines through integrated, high-performance panic recovery mechanisms.
* Low Memory Footprint: Minimizing memory overhead through efficient, unsafe-supported pointer conversions within the object pools.

## Usage

The following architecture example demonstrates the implementation of a type-safe, isolated pipeline for a fictional sub-module.

### 1. Definition of the Module Contract

Each module defines its own dedicated command sets via Go constants. This enforces isolation at the compiler level.

```go
package main

import "github.com/whazzabii7/rpr"

// Definition of the module-specific type
type EngineCommand int

// Definition of unique commands within the module
const (
	CmdEngineSpawn EngineCommand = iota
	CmdEngineDestroy
)

// Definition of the payload structure
type EnginePayload struct {
	ClusterID string
	NodeCount int
}
```

### 2. Request Transmission and Type-Safe Assignment

The sender packages the data allocation-free. The receiver (worker) unpacks it without performance loss using rpr.Assign.

```go
func handlePipeline(queue chan *rpr.Request[EngineCommand], responseChan chan *rpr.Response) {
	for req := range queue {
		// Allocation-free protection against unforeseen panics in the worker
		defer rpr.Recover(req.Response)

		switch req.Type {
		case CmdEngineSpawn:
			var data EnginePayload
			
			// Type-safe mapping without reflection
			if rpr.Assign(req.Payload.Get1(), &data) {
				// Processing business logic
				// ...
				
				// Return successful response, payload must not be nil for CheckResponse
				rpr.NewResponse(req.Payload, nil).Submit(req.Response)
			} else {
				rpr.NewResponseErr(rpr.ErrUnpackPayloadFail).Submit(req.Response)
			}
		}
		
		// Important: Return resources to the sync.Pool
		req.Release()
	}
}
```

### 3. Pipeline Initialization and Invocation

```go
func main() {
	// Creating a type-safe pipeline
	pipeline := rpr.MakeRequestChan[EngineCommand](100)
	resChan := rpr.MakeResponseChan()

	// Initializing the data structure
	payloadData := EnginePayload{ClusterID: "eu-central-1", NodeCount: 12}

	// Create request, pack payload, and feed into the pipeline
	rpr.NewRequest(CmdEngineSpawn, rpr.Pack(&payloadData), resChan).Submit(pipeline)

	// Validating the response
	if res, ok := rpr.CheckResponse(resChan); ok {
		// Transaction successfully completed
		res.Release()
	}
}
```

## Benchmarks

The performance metrics were tested under real-world conditions. The framework achieves a constant allocation rate of zero throughout the entire lifecycle.

```go
go test -bench=. -benchmem
```

Test Results:
```go
BenchmarkRPR_SpeedAndAllocations-16    184203108    6.421 ns/op    0 B/op    0 allocs/op
PASS
ok      github.com/whazzabii7/rpr       1.243s
```

## Installation

To integrate the module into an existing Go project, invoke the Go package manager:

```go
go get github.com/whazzabii7/rpr@v1.0.0
```

## License

This framework is licensed under the MIT License. See the LICENSE file for detailed information.
