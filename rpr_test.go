package rpr

import (
	"errors"
	"testing"
)

// TestTestData ist eine fiktive Struktur für den Test
type TestData struct {
	Value string
}

// 1. Funktionaler Test: Der Standard-Pipeline-Durchlauf
func TestRPR_Lifecycle(t *testing.T) {
	reqChan := MakeRequestChan[int](10)
	resChan := MakeResponseChan()

	// Sender sendet Daten
	data := TestData{Value: "Enterprise"}
	NewRequest(1, Pack(&data), resChan).Submit(reqChan)

	// Empfänger verarbeitet Daten
	req := <-reqChan
	var received TestData
	if !Assign(req.Payload.Get1(), &received) {
		t.Fatal("Assign fehlgeschlagen")
	}

	if received.Value != "Enterprise" {
		t.Errorf("Erwartet 'Enterprise', got '%s'", received.Value)
	}

	// Antwort senden
	NewResponse(nil, nil).Submit(req.Response)
	req.Release()

	// Sender prüft Antwort
	res, ok := CheckResponse(resChan)
	if !ok || res.Err != nil {
		t.Fatal("Response-Validierung fehlgeschlagen")
	}
	res.Release()
}

// 2. Test der Panic Recovery
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

// 3. BENCHMARK: Beweise die Allokationsfreiheit (0 B/op)
func BenchmarkRPR_SpeedAndAllocations(b *testing.B) {
	resChan := MakeResponseChan()
	data := TestData{Value: "Fast"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Packen & Request bauen
		req := NewRequest(1, Pack(&data), resChan)
		
		// Direktes Entpacken simulieren
		var dest TestData
		Assign(req.Payload.Get1(), &dest)
		
		// Aufräumen
		req.Release()
	}
}
