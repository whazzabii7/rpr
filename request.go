package rpr

import "unsafe"

// RequestType represents the underlying base type for system commands.
type RequestType int

// RequestConstraint defines the type constraint for compilation-time safety,
// ensuring only integer-based custom enums can be used as module commands.
type RequestConstraint interface{ ~int }

// rawRequest is the internal, non-generic memory layout used to store and
// recycle request structures within the sync.Pool without generic type overhead.
type rawRequest struct {
	Type     int
	Payload  Payload
	Response chan *Response
}

// Request represents a strongly-typed, allocation-free message envelope
// carrying a command type, an optional data payload, and a response channel.
type Request[T RequestConstraint] struct {
	Type     T              `json:"type"`
	Payload  Payload        `json:"payload"`
	Response chan *Response `json:"response"`
}

// NewRequest creates a new strongly-typed Request by retrieving an underlying
// memory block from the sync.Pool and safely casting it via unsafe.Pointer.
// This guarantees zero-allocation message creation.
func NewRequest[T RequestConstraint](reqType T, payload Payload, response chan *Response) *Request[T] {
	raw := requestPool.Get().(*rawRequest)
	req := (*Request[T])(unsafe.Pointer(raw))

	req.Type = reqType
	req.Payload = payload
	req.Response = response

	return req
}

// Submit non-blockingly or blockingly pushes the request into the provided
// module-specific processing pipeline channel.
func (r *Request[T]) Submit(requestCh chan *Request[T]) {
	requestCh <- r
}

// Release completely zeroes out the request data, releases the attached Payload
// back to its pool, and returns the request memory back to the internal sync.Pool.
// Must be called by the worker/receiver once the request handling is finished.
func (r *Request[T]) Release() {
	if r == nil {
		return
	}

	if r.Payload != nil {
		r.Payload.Release()
	}
	r.Response = nil

	requestPool.Put((*rawRequest)(unsafe.Pointer(r)))
}
