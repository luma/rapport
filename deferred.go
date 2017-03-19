package rapport

import "github.com/luma/pith/rapport/causality"

func MakeDeferredSet() *DeferredSet {
	return &DeferredSet{
		Members: make(map[string]bool),
	}
}

func UnmarshalDeferredSet(data []byte) (*DeferredSet, error) {
	set := MakeDeferredSet()
	if err := set.Unmarshal(data); err != nil {
		return nil, err
	}

	return set, nil
}

func (d *DeferredSet) Clone() *DeferredSet {
	deferred := MakeDeferredSet()
	for value := range d.Members {
		deferred.Members[value] = true
	}

	return deferred
}

func (d *DeferredSet) Values() []string {
	values := make([]string, 0, len(d.Members))
	for value := range d.Members {
		values = append(values, value)
	}

	return values
}

type DeferredMap map[*causality.VersionVector]*DeferredSet

func (d *DeferredMap) Clone() *DeferredMap {
	deferredMap := make(DeferredMap)

	for version, oldDeferred := range *d {
		deferredMap[version] = oldDeferred.Clone()
	}

	return &deferredMap
}
