package marshalling

func CreatePNCounter(replicaId string) *PNCounterValue {
	p := &PNCounterValue{
		Inc: make(map[string]int64),
		Dec: make(map[string]int64),
	}

	incZero := int64(0)
	decZero := int64(0)

	p.Inc[replicaId] = incZero
	p.Dec[replicaId] = decZero

	return p
}

func (p *PNCounterValue) Incr(replicaId string) int64 {
	value := p.Inc[replicaId] + 1
	p.Inc[replicaId] = value
	return value
}

func (p *PNCounterValue) IncrBy(replicaId string, amount int64) int64 {
	if amount > 0 {
		value := p.Inc[replicaId] + amount
		p.Inc[replicaId] = value
	} else if amount < 0 {
		value := p.Dec[replicaId] - amount
		p.Dec[replicaId] = value
	}

	return p.Value()
}

func (p *PNCounterValue) Decr(replicaId string) (int64, error) {
	value := p.Dec[replicaId] + 1
	p.Dec[replicaId] = value

	return value, nil
}

func (p *PNCounterValue) DecrBy(replicaId string, amount int64) (int64, error) {
	return p.IncrBy(replicaId, -amount), nil
}

func (p *PNCounterValue) Value() (total int64) {
	for id, IncVal := range p.Inc {
		total += IncVal
		total -= p.Dec[id]
	}

	return total
}

// func (p *PNCounterValue) Merge(crdt CRDT) {
// 	other := crdt.(*PNCounterValue)
//
// 	for id, IncVal := range other.Inc {
// 		if localInc, exists := p.Inc[id]; !exists || localInc < IncVal {
// 			p.Inc[id] = IncVal
// 		}
//
// 		if localDec, exists := p.Dec[id]; !exists || localDec < other.Dec[id] {
// 			p.Dec[id] = other.Dec[id]
// 		}
// 	}
// }
