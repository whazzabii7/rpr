package rpr

// Payload defines a generic, allocation-free container interface for wrapping
// and passing multiple data values across isolated system channels.
type Payload interface {
	// Get1 returns the single encapsulated value if the payload is of kind 1.
	// Returns nil if the payload kind does not match.
	Get1() any

	// Get2 returns the two encapsulated values if the payload is of kind 2.
	// Returns nil, nil if the payload kind does not match.
	Get2() (any, any)

	// Get3 returns the three encapsulated values if the payload is of kind 3.
	// Returns nil, nil, nil if the payload kind does not match.
	Get3() (any, any, any)

	// Release clears all underlying references and returns the payload object
	// back to the internal sync.Pool to completely prevent heap allocations.
	Release()
}

// rawPayload is the internal concrete implementation of the Payload interface.
// Comments here are kept minimal as internal types are hidden from external IDE docs.
type rawPayload struct {
	kind int // Identifies the payload kind (1, 2, or 3)
	v1   any
	v2   any
	v3   any
}

// Pack wraps a single data value into an allocation-free Payload container
// recycled from the internal sync.Pool.
func Pack(t any) Payload {
	raw := payloadPool.Get().(*rawPayload)
	raw.kind = 1
	raw.v1 = t
	return raw
}

// Pack2 wraps two data values into an allocation-free Payload container
// recycled from the internal sync.Pool.
func Pack2(t1, t2 any) Payload {
	raw := payloadPool.Get().(*rawPayload)
	raw.kind = 2
	raw.v1 = t1
	raw.v2 = t2
	return raw
}

// Pack3 wraps three data values into an allocation-free Payload container
// recycled from the internal sync.Pool.
func Pack3(t1, t2, t3 any) Payload {
	raw := payloadPool.Get().(*rawPayload)
	raw.kind = 3
	raw.v1 = t1
	raw.v2 = t2
	raw.v3 = t3
	return raw
}

func (p *rawPayload) Get1() any {
	if p == nil || p.kind != 1 {
		return nil
	}
	return p.v1
}

func (p *rawPayload) Get2() (any, any) {
	if p == nil || p.kind != 2 {
		return nil, nil
	}
	return p.v1, p.v2
}

func (p *rawPayload) Get3() (any, any, any) {
	if p == nil || p.kind != 3 {
		return nil, nil, nil
	}
	return p.v1, p.v2, p.v3
}

func (p *rawPayload) Release() {
	if p == nil {
		return
	}
	p.v1 = nil
	p.v2 = nil
	p.v3 = nil
	p.kind = 0
	payloadPool.Put(p)
}
