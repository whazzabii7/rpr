package rpr

import "unsafe"

type RequestType int
type RequestConstraint interface{ ~int }

type rawRequest struct {
	Type     int
	Payload  Payload
	Response chan *Response
}

type Request[T RequestConstraint] struct {
	Type     T              `json:"type"`
	Payload  Payload        `json:"payload"`
	Response chan *Response `json:"response"`
}

func NewRequest[T RequestConstraint](reqType T, payload Payload, response chan *Response) *Request[T] {
	raw := requestPool.Get().(*rawRequest)
	req := (*Request[T])(unsafe.Pointer(raw))

	req.Type = reqType
	req.Payload = payload
	req.Response = response

	return req
}

func (r *Request[T]) Submit(requestCh chan *Request[T]) {
	requestCh <- r
}

func (r *Request[T]) Release() {
	if r == nil {
		return
	}

	if r.Payload != nil { r.Payload.Release() }
	r.Response = nil

	requestPool.Put((*rawRequest)(unsafe.Pointer(r)))
}
