# RPR (Remote Process Routing)

High-Performance In-Memory Communication Matrix for Go.

## Description

RPR ist eine minimalistische, hocheffiziente und typsichere In-Memory-Kommunikationsmatrix für die Go-Laufzeitumgebung. Das Framework wurde speziell für eingebettete, ereignisgesteuerte Micro-Engines und nach dem Actor-Modell designte Systeme entwickelt, bei denen Latenzzeiten im Nanosekundenbereich und die absolute Vermeidung von Garbage-Collection-Overhead (GC) kritische Systemanforderungen darstellen.

Durch den konsequenten Einsatz von internem Speicher-Recycling via sync.Pool und den Verzicht auf rechenintensive Laufzeit-Reflection garantiert RPR eine allokationsfreie Datenübertragung zwischen isolierten Systemkomponenten.

## AIM-Goal

Die Kernziele des Frameworks sind streng auf die Anforderungen von Hochverfügbarkeitssystemen (SLA-driven Enterprise Architectures) ausgerichtet:

* Zero-Allocation Lifecycle: Vollständige Eliminierung von Heap-Allokationen während des gesamten Request-Response-Durchlaufs zur Gewährleistung deterministischer Latenzzeiten.
* Compile-Time Type Safety: Physische Isolation von Modul-Pipelines durch den Einsatz von Go-Generics und Typ-Constraints (RequestConstraint). Fehlerhafte Daten- oder Befehlsadressierungen werden bereits zur Kompilierzeit und nicht erst zur Laufzeit abgefangen.
* Resilience & Fault Isolation: Absicherung langlebiger Worker-Goroutinen durch integrierte, performante Panic-Recovery-Mechanismen.
* Low Memory Footprint: Minimierung des Memory-Overheads durch effiziente, unsafe-unterstützte Zeiger-Konvertierungen innerhalb der Objekt-Pools.

## Usage

Das folgende Architektur-Beispiel demonstriert die Implementierung einer typsicheren, isolierten Pipeline für ein fiktives Submodul.

### 1. Definition des Modul-Vertrags (Contract)

Jedes Modul definiert seine eigenen, dedizierten Befehlssätze über Go-Konstanten. Dies erzwingt die Isolation auf Compiler-Ebene.

```go
package main

import "[github.com/whazzabii7/rpr](https://github.com/whazzabii7/rpr)"

// Definition des modulspezifischen Typs
type EngineCommand int

// Definition der eindeutigen Befehle innerhalb des Moduls
const (
	CmdEngineSpawn EngineCommand = iota
	CmdEngineDestroy
)

// Definition der Payload-Struktur
type EnginePayload struct {
	ClusterID string
	NodeCount int
}
```

### 2. Request-Übertragung und Typsichere Zuweisung

Der Sender verpackt die Daten allokationsfrei. Der Empfänger (Worker) entpackt diese ohne Performance-Verlust mittels rpr.Assign.

```go
func handlePipeline(queue chan *rpr.Request[EngineCommand], responseChan chan *rpr.Response) {
	for req := range queue {
		// Allokationsfreie Absicherung gegen unvorhergesehene Panics im Worker
		defer rpr.Recover(req.Response)

		switch req.Type {
		case CmdEngineSpawn:
			var data EnginePayload
			
			// Typsicheres Mapping ohne Reflection
			if rpr.Assign(req.Payload.Get1(), &data) {
				// Verarbeitung der Geschäftslogik
				// ...
				
				// Erfolgreiche Antwort zurückgeben, Payload darf nicht nil sein für CheckResponse
				rpr.NewResponse(req.Payload, nil).Submit(req.Response)
			} else {
				rpr.NewResponseErr(rpr.ErrUnpackPayloadFail).Submit(req.Response)
			}
		}
		
		// Wichtig: Ressourcen zurück in den sync.Pool übergeben
		req.Release()
	}
}
```

### 3. Pipeline-Initialisierung und Aufruf

```go
func main() {
	// Erstellung einer typsicheren Pipeline
	pipeline := rpr.MakeRequestChan[EngineCommand](100)
	resChan := rpr.MakeResponseChan()

	// Initialisierung der Datenstruktur
	payloadData := EnginePayload{ClusterID: "eu-central-1", NodeCount: 12}

	// Request erstellen, Payload packen und in die Pipeline einspeisen
	rpr.NewRequest(CmdEngineSpawn, rpr.Pack(&payloadData), resChan).Submit(pipeline)

	// Validierung der Antwort
	if res, ok := rpr.CheckResponse(resChan); ok {
		// Transaktion erfolgreich abgeschlossen
		res.Release()
	}
}
```

## Benchmarks

Die Performance-Metriken wurden unter realen Bedingungen getestet. Das Framework erzielt über den gesamten Lifecycle hinweg eine konstante Allokationsrate von Null.

```go
go test -bench=. -benchmem
```

Testergebnisse:
```go
BenchmarkRPR_SpeedAndAllocations-16    184203108    6.421 ns/op    0 B/op    0 allocs/op
PASS
ok      [github.com/whazzabii7/rpr](https://github.com/whazzabii7/rpr)       1.243s
```

## Installation

Zur Integration des Moduls in ein bestehendes Go-Projekt muss der Go-Paketmanager aufgerufen werden:

```go
go get [github.com/whazzabii7/rpr@v1.0.0](https://github.com/whazzabii7/rpr@v1.0.0)
```

## License

Dieses Framework ist lizenziert unter der MIT-Lizenz. Siehe die LICENSE Datei für detaillierte Informationen.
