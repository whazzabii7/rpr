package rpr

// Response represents an allocation-free response envelope containing
// the processing results (Payload) or any error encountered during execution.
type Response struct {
	Payload Payload `json:"payload"`
	Err     error   `json:"err"`
}

// NewResponse creates a new Response envelope by retrieving an object from
// the internal sync.Pool, wrapping both the given payload and error.
func NewResponse(payload Payload, err error) *Response {
	res := responsePool.Get().(*Response)

	res.Payload = payload
	res.Err = err

	return res
}

// NewResponseErr creates an error-only Response envelope by retrieving an object
// from the internal sync.Pool, specifically used when no payload data is produced.
func NewResponseErr(err error) *Response {
	res := responsePool.Get().(*Response)

	res.Payload = nil
	res.Err = err

	return res
}

// Release zeroes out the response fields, automatically releases the enclosed
// Payload back to its pool, and returns the Response object to the responsePool.
// Must be called by the client/sender after processing the response to prevent leaks.
func (r *Response) Release() {
	if r == nil {
		return
	}

	if r.Payload != nil {
		r.Payload.Release()
	}
	r.Err = nil

	responsePool.Put(r)
}

// Submit non-blockingly or blockingly pushes the response envelope back
// into the provided execution-specific response channel.
func (r *Response) Submit(responseCh chan *Response) {
	responseCh <- r
}

// CheckResponse blocks on the provided response channel, reads the incoming
// Response envelope, and validates its error state.
// Returns (responseData, true) if the operation was successful.
// Returns (responseData, false) if an error occurred or the channel was closed.
func CheckResponse(response <-chan *Response) (*Response, bool) {
	responseData, ok := <-response
	if !ok {
		responseData = NewResponseErr(ErrNotAResponse)
	}
	if responseData.Err != nil {
		return responseData, false
	}
	return responseData, true
}
