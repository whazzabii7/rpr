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


func MakeRequestChan[T RequestConstraint](size int) chan *Request[T] {
	if size < 1 {
		size = 1
	}
	return make(chan *Request[T], size)
}

func MakeResponseChan() chan *Response {
	return make(chan *Response, 1)
}

func Ptr[T any](v T) *T {
	return &v
}

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

func Recover(responseCh chan *Response) {
	if r := recover(); r != nil {
		err := fmt.Errorf("critical panic recovered: %v", r)
		NewResponseErr(err).Submit(responseCh)
	}
}
