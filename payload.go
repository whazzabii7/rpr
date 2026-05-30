package rpr

type Payload interface {
	Get1() any
	Get2() (any, any)
	Get3() (any, any, any)
	Release()
}

type rawPayload struct {
	kind int // 1, 2 oder 3
	v1   any
	v2   any
	v3   any
}

func Pack(t any) Payload {
	raw := payloadPool.Get().(*rawPayload)
	raw.kind = 1
	raw.v1 = t
	return raw
}

func Pack2(t1, t2 any) Payload {
	raw := payloadPool.Get().(*rawPayload)
	raw.kind = 2
	raw.v1 = t1
	raw.v2 = t2
	return raw
}

func Pack3(t1, t2, t3 any) Payload {
	raw := payloadPool.Get().(*rawPayload)
	raw.kind = 3
	raw.v1 = t1
	raw.v2 = t2
	raw.v3 = t3
	return raw
}

func (p *rawPayload) Get1() any {
	if p == nil || p.kind != 1 { return nil }
	return p.v1
}

func (p *rawPayload) Get2() (any, any) {
	if p == nil || p.kind != 2 { return nil, nil }
	return p.v1, p.v2
}

func (p *rawPayload) Get3() (any, any, any) {
	if p == nil || p.kind != 3 { return nil, nil, nil }
	return p.v1, p.v2, p.v3
}

func (p *rawPayload) Release() {
	if p == nil { return }
	p.v1 = nil
	p.v2 = nil
	p.v3 = nil
	p.kind = 0
	payloadPool.Put(p)
}
