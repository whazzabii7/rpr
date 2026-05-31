package rpr

import (
	"fmt"
	"sync"
)

// Global internal sync.Pools used for automatic memory recycling.
// These are unexported (lowercase) and do not need public IDE documentation.
var requestPool = sync.Pool{
	New: func() any { return new(rawRequest) },
}

var responsePool = sync.Pool{
	New: func() any { return new(Response) },
}

var payloadPool = sync.Pool{
	New: func() any { return new(rawPayload) },
}

// MakeRequestChan instantiates a buffered channel for strongly-typed Requests.
// If the provided size is less than 1, it automatically defaults to a buffer size of 1.
func MakeRequestChan[T RequestConstraint](size int) chan *Request[T] {
	if size < 1 {
		size = 1
	}
	return make(chan *Request[T], size)
}

// MakeResponseChan instantiates a standardized response channel with a fixed
// buffer size of 1 to prevent blocking during worker routing.
func MakeResponseChan() chan *Response {
	return make(chan *Response, 1)
}

// Ptr is a generic helper that takes any value and returns a pointer to it.
// Highly useful for inline structural initializations.
func Ptr[T any](v T) *T {
	return &v
}

// Assign performs a reflection-free, type-safe pointer type assertion.
// It maps the content of src into dst and returns true if successful.
// Returns false if either src or dst is nil, or if the underlying types do not match.
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

// Recover acts as a high-performance panic recovery mechanism for worker goroutines.
// It intercept panics, formats them into an error, and automatically submits an
// error Response back to the sender via responseCh.
// An optional callback (onPanic) can be provided for custom logging or metrics.
// Must be called via 'defer rpr.Recover(ch)'.
func Recover(responseCh chan *Response, onPanic ...func(panicErr any)) {
	r := recover()
	if r == nil {
		return
	}
	err := fmt.Errorf("critical panic recovered: %v", r)
	NewResponseErr(err).Submit(responseCh)
	if len(onPanic) > 0 && onPanic[0] != nil {
		onPanic[0](r)
	}
}
