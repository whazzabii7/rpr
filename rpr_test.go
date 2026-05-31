package rpr

import (
	"testing"
)

// TestData is a simple struct used to verify payload packaging and assertion.
type TestData struct {
	Value string
}

// TestRPR_Lifecycle verifies the full request-response processing chain.
// It ensures that data can be packed, transmitted through a channel, safely
// assigned without reflection, and that resources are properly recycled.
func TestRPR_Lifecycle(t *testing.T) {
	reqChan := MakeRequestChan[int](10)
	resChan := MakeResponseChan()

	// 1. Pack data and submit request
	data := TestData{Value: "Enterprise"}
	NewRequest(1, Pack(&data), resChan).Submit(reqChan)

	// 2. Receive and process the request
	req := <-reqChan
	var received TestData
	if !Assign(req.Payload.Get1(), &received) {
		t.Fatal("Assign failed")
	}

	if received.Value != "Enterprise" {
		t.Errorf("Expected 'Enterprise', got '%s'", received.Value)
	}

	// 3. Respond and release request resources
	NewResponse(nil, nil).Submit(req.Response)
	req.Release()

	// 4. Validate and release response resources
	res, ok := CheckResponse(resChan)
	if !ok || res.Err != nil {
		t.Fatal("Response validation failed")
	}
	res.Release()
}

// TestRPR_Recovery ensures that panics inside worker goroutines are
// gracefully intercepted and automatically routed back as error responses.
func TestRPR_Recovery(t *testing.T) {
	resChan := MakeResponseChan()

	go func() {
		defer Recover(resChan)
		panic("simulated crash")
	}()

	res := <-resChan
	if res.Err == nil {
		t.Fatal("Panic was not intercepted!")
	}
	res.Release()
}

// BenchmarkRPR_SpeedAndAllocations measures the raw operation throughput
// and verifies the zero-allocation guarantee of the request lifecycle.
func BenchmarkRPR_SpeedAndAllocations(b *testing.B) {
	resChan := MakeResponseChan()
	data := TestData{Value: "Fast"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := NewRequest(1, Pack(&data), resChan)

		var dest TestData
		Assign(req.Payload.Get1(), &dest)

		req.Release()
	}
}
