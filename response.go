package rpr

type Response struct {
	Payload Payload `json:"payload"`
	Err     error   `json:"err"`
}

func NewResponse(payload Payload, err error) *Response {
	res := responsePool.Get().(*Response)

	res.Payload = payload
	res.Err = err

	return res
}

func NewResponseErr(err error) *Response {
	res := responsePool.Get().(*Response)

	res.Payload = nil
	res.Err = err

	return res
}

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

func (r *Response) Submit(responseCh chan *Response) {
	responseCh <- r
}

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
