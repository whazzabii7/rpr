package rpr

import (
	"testing"
)

type TestData struct {
	Value string
}

func TestRPR_Lifecycle(t *testing.T) {
	reqChan := MakeRequestChan[int](10)
	resChan := MakeResponseChan()

	data := TestData{Value: "Enterprise"}
	NewRequest(1, Pack(&data), resChan).Submit(reqChan)

	req := <-reqChan
	var received TestData
	if !Assign(req.Payload.Get1(), &received) {
		t.Fatal("Assign fehlgeschlagen")
	}

	if received.Value != "Enterprise" {
		t.Errorf("Erwartet 'Enterprise', got '%s'", received.Value)
	}

	NewResponse(nil, nil).Submit(req.Response)
	req.Release()

	res, ok := CheckResponse(resChan)
	if !ok || res.Err != nil {
		t.Fatal("Response-Validierung fehlgeschlagen")
	}
	res.Release()
}

func TestRPR_Recovery(t *testing.T) {
	resChan := MakeResponseChan()

	go func() {
		defer Recover(resChan)
		panic("simulierter Absturz")
	}()

	res := <-resChan
	if res.Err == nil {
		t.Fatal("Panic wurde nicht abgefangen!")
	}
	res.Release()
}

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
