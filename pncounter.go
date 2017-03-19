package rapport

import (
	"sync"

	"github.com/luma/pith/rapport/marshalling"
)

type PNCounter struct {
	replicaId string

	value *marshalling.PNCounterValue
	l     sync.RWMutex
}

func CreatePNCounter(replicaId string) *PNCounter {
	return &PNCounter{
		replicaId: replicaId,
		value:     marshalling.CreatePNCounter(replicaId),
	}
}

func (p *PNCounter) Incr() int64 {
	p.l.Lock()
	value := p.value.Incr(p.replicaId)
	p.l.Unlock()

	return value
}

func (p *PNCounter) IncrBy(amount int64) int64 {
	p.l.Lock()
	value := p.value.IncrBy(p.replicaId, amount)
	p.l.Unlock()

	return value
}

func (p *PNCounter) Decr() (int64, error) {
	p.l.Lock()
	value, err := p.value.Decr(p.replicaId)
	p.l.Unlock()

	return value, err
}

func (p *PNCounter) DecrBy(amount int64) (int64, error) {
	return p.value.DecrBy(p.replicaId, amount)
}

func (p *PNCounter) Value() (total int64) {
	p.l.RLock()
	value := p.value.Value()
	p.l.RUnlock()

	return value
}

func (p *PNCounter) Merge(crdt CRDT) {
	other := crdt.(*PNCounter)

	p.l.Lock()
	other.l.Lock()

	defer func() {
		p.l.Unlock()
		other.l.Unlock()
	}()

	for id, incVal := range other.value.Inc {
		if localInc, exists := p.value.Inc[id]; !exists || localInc < incVal {
			p.value.Inc[id] = incVal
		}

		if localDec, exists := p.value.Dec[id]; !exists || localDec < other.value.Dec[id] {
			p.value.Dec[id] = other.value.Dec[id]
		}
	}
}

// Marshal serialises the counter data to bytes
func (p *PNCounter) Marshal() ([]*Segment, error) {
	v, err := p.value.Marshal()
	if err != nil {
		return nil, err
	}

	segment := &Segment{
		Value: v,
	}
	return []*Segment{segment}, nil
}

// Marshal deserialises the counter data from bytes
func (p *PNCounter) Unmarshal(data []*Segment) error {
	if p.value == nil {
		p.value = marshalling.CreatePNCounter(p.replicaId)
	}

	return p.value.Unmarshal(data[0].Value)
}
