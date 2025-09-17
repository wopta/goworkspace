package bpmnEngine

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
)

type InitializeDataType func() DataBpnm

type StorageBpnm struct {
	local       map[string]any
	global      map[string]any
	marked      []string
	higherStore *StorageBpnm
}

// The storage manages his own data:
// At each cycles:
// It cleans itself leaving only the marked data (markWhatNeeded)
func NewStorageBpnm() *StorageBpnm {
	res := new(StorageBpnm)
	res.local = make(map[string]any)
	res.global = make(map[string]any)
	res.higherStore = nil
	return res
}

func (p *StorageBpnm) ResetLocal() {
	p.local = make(map[string]any)
}

func (p *StorageBpnm) ResetGlobal() {
	p.global = make(map[string]any)
}

func (p *StorageBpnm) markWhatNeeded(toTouch []typeData) {
	var toHigherStorage []typeData
	for _, t := range toTouch {
		if _, exist := p.local[t.Name]; !exist {
			toHigherStorage = append(toHigherStorage, t)
			continue
		}
		p.marked = append(p.marked, t.Name)
	}
}

func (p *StorageBpnm) cleanNoMarkedResources() error {
	backup := maps.Clone(p.local)
	p.ResetLocal()
	for i, toSave := range p.marked {
		if backup[toSave] == nil {
			continue
		}
		if err := p.AddLocal(p.marked[i], backup[toSave].(DataBpnm)); err != nil {
			return err
		}
	}
	p.marked = nil
	return nil
}

func (p *StorageBpnm) getAllLocals() map[string]any {
	var res map[string]any = make(map[string]any)
	res = p.local
	if p.higherStore == nil {
		return res
	}
	res = mergeMaps(p.higherStore.getAllLocals(), res)
	return res
}
func (p *StorageBpnm) getAllGlobals() map[string]any {
	var res map[string]any = make(map[string]any)
	res = p.global
	if p.higherStore == nil {
		return res
	}
	res = mergeMaps(p.higherStore.getAllGlobals(), res)
	return res
}

func (p *StorageBpnm) AddLocal(name string, data DataBpnm) error {
	if p.local == nil {
		return errors.New("Error in the initialization of the storage")
	}
	if _, ok := p.local[name]; ok {
		return fmt.Errorf("Storage has already data with name %v", name)
	}
	p.local[name] = data
	return nil
}

func (p *StorageBpnm) AddGlobal(name string, data DataBpnm) error {
	if p.global == nil {
		return errors.New("Error in the initialization of the storage")
	}
	if _, ok := p.global[name]; ok {
		return fmt.Errorf("Storage has already data with name %v", name)
	}
	p.global[name] = data
	return nil
}

// Get global data if no try with higher store
func (p *StorageBpnm) GetGlobal(name string) (DataBpnm, error) {
	if p.global == nil {
		return nil, errors.New("Error in the initialization of the storage")
	}
	if data, ok := p.global[name]; ok {
		return data.(DataBpnm), nil
	}
	if p.higherStore != nil {
		return p.higherStore.GetGlobal(name)
	}
	return nil, fmt.Errorf("No data found %v", name)
}

// Get local data if no try with higher store
func (p *StorageBpnm) GetLocal(name string) (DataBpnm, error) {
	if p.local == nil {
		return nil, errors.New("Error in the initialization of the storage")
	}
	if data, ok := p.local[name]; ok {
		return data.(DataBpnm), nil
	}
	if p.higherStore != nil {
		return p.higherStore.GetLocal(name)
	}
	return nil, fmt.Errorf("No data found %v", name)
}

func (p *StorageBpnm) GetMap() map[string]any {
	mapsMerged := mergeMaps(p.getAllGlobals(), p.getAllLocals())
	bytes, err := json.Marshal(mapsMerged)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(bytes, &mapsMerged)
	if err != nil {
		return nil
	}
	return mapsMerged
}
func (base *StorageBpnm) setHigherStorage(higher *StorageBpnm) error {
	if base == higher {
		return fmt.Errorf("A storage can't reference itself as Higher-level storage")
	}
	base.higherStore = higher
	return nil
}

func (base *StorageBpnm) mergeUnique(source *StorageBpnm) error {
	var err error
	if source == nil {
		return nil
	}
	base.global, err = mergeUniqueMaps(base.global, source.getAllGlobals())
	if err != nil {
		return err
	}
	base.local, err = mergeUniqueMaps(base.local, source.getAllLocals())
	if err != nil {
		return err
	}
	return nil
}
