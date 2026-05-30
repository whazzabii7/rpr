package rpr

import (
	"fmt"
	"sync"
)


var requestPool = sync.Pool{
	New: func() any { return new(rawRequest) },
}

var responsePool = sync.Pool{
	New: func() any { return new(Response) },
}

var payloadPool = sync.Pool{
	New: func() any { return new(rawPayload) },
}


// MakeRequestChan erzwingt einen Mindest-Puffer von 1 für asynchrone Pipelines.
func MakeRequestChan[T RequestConstraint](size int) chan *Request[T] {
	if size < 1 {
		size = 1
	}
	return make(chan *Request[T], size)
}

// MakeResponseChan erzwingt exakt einen Puffer von 1. 
// Da eine Response pro Request erwartet wird, eliminiert das Goroutine-Leaks zu 100 %.
func MakeResponseChan() chan *Response {
	return make(chan *Response, 1)
}

// Ptr ist eine generische Hilfsfunktion, um Werte direkt als Pointer zu übergeben.
// Nützlich für Pack(), z.B.: rpr.Pack(rpr.Ptr("mein_string"))
func Ptr[T any](v T) *T {
	return &v
}

// Assign castet einen rohen Interface-Wert (any) typsicher in den Zielpointer.
// Bietet 0 Allokationen und benötigt keine langsame Reflection.
func Assign[T any](src any, dst *T) bool {
	if src == nil || dst == nil {
		return false
	}
	if val, ok := src.(*T); ok {
		*dst = *val
		return true
	}
	return false
}

// Recover sichert eine laufende Goroutine gegen unvorhergesehene Abstürze (Panics) ab.
// Schreibt die Fehlermeldung direkt in den Response-Channel des Callers.
func Recover(responseCh chan *Response) {
	if r := recover(); r != nil {
		err := fmt.Errorf("critical panic recovered: %v", r)
		NewResponseErr(err).Submit(responseCh)
	}
}
