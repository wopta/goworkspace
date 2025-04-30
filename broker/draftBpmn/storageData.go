package draftbpmn

import (
	"errors"
	"fmt"
	"maps"
)

type InitializeDataType func() DataBpnm
type StorageData interface {
	ResetLocal()
	ResetGlobal()
	AddLocal(string, DataBpnm) error
	AddGlobal(string, DataBpnm) error
	GetLocal(string) (DataBpnm, error)
	GetGlobal(string) (DataBpnm, error)

	getAllLocal() map[string]any
	getAllGlobal() map[string]any
	// setHigherStorage sets a higher-level storage,
	// which will be used as a fallback when a local/global key is not found in the current storage.
	setHigherStorage(StorageData) error
	// It merges two unique storage
	// If both storage contain the same key, return error
	mergeUnique(StorageData) error
	//Mark what local resources keep when clean is called
	markWhatNeeded([]typeData)
	//Delete the resources that aren't needed(aren't marked)
	clean() error
}

type StorageBpnm struct {
	local       map[string]any
	global      map[string]any
	touched     []string
	higherStore StorageData
}

// The storage manages his own resources:
// At each cycles:
// It cleans itself leaving only the output resources DECLARED
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
	for _, t := range toTouch {
		p.touched = append(p.touched, t.Name)
	}
}

func (p *StorageBpnm) clean() error {
	backup := maps.Clone(p.local)
	p.ResetLocal()
	for i, toSave := range p.touched {
		if err := p.AddLocal(p.touched[i], backup[toSave].(DataBpnm)); err != nil {
			return err
		}
	}
	p.touched = nil
	return nil
}

func (p *StorageBpnm) getAllLocal() map[string]any {
	var res map[string]any = make(map[string]any)
	res = p.local
	if p.higherStore == nil {
		return res
	}
	res = mergeMaps(p.higherStore.getAllLocal(), res)
	return res
}

func (p *StorageBpnm) getAllGlobal() map[string]any {
	var res map[string]any = make(map[string]any)
	res = p.global
	if p.higherStore == nil {
		return res
	}
	res = mergeMaps(p.higherStore.getAllGlobal(), res)
	return res
}

func (p *StorageBpnm) AddLocal(name string, data DataBpnm) error {
	if p.local == nil {
		return errors.New("error initialization local storage")
	}
	if _, ok := p.local[name]; ok {
		return fmt.Errorf("storage has already data with name %v", name)
	}
	p.local[name] = data
	return nil
}

func (p *StorageBpnm) AddGlobal(name string, data DataBpnm) error {
	if p.global == nil {
		return errors.New("error initialization storage storage")
	}
	if _, ok := p.global[name]; ok {
		return fmt.Errorf("storage has already data with name %v", name)
	}
	p.global[name] = data
	return nil
}

func (p *StorageBpnm) GetGlobal(name string) (DataBpnm, error) {
	if p.global == nil {
		return nil, errors.New("error initialization storage storage")
	}
	if data, ok := p.global[name]; ok {
		return data.(DataBpnm), nil
	}
	if p.higherStore != nil {
		return p.higherStore.GetGlobal(name)
	}
	return nil, fmt.Errorf("no data founded %v", name)
}

func (p *StorageBpnm) GetLocal(name string) (DataBpnm, error) {
	if p.global == nil {
		return nil, errors.New("error initialization storage storage")
	}
	if data, ok := p.local[name]; ok {
		return data.(DataBpnm), nil
	}
	if p.higherStore != nil {
		return p.higherStore.GetLocal(name)
	}
	return nil, fmt.Errorf("no data founded %v", name)
}

func (base *StorageBpnm) setHigherStorage(higher StorageData) error {
	if base.higherStore != nil {
		return fmt.Errorf("Higher storage has been already set")
	}
	base.higherStore = higher
	return nil
}

func (base *StorageBpnm) mergeUnique(source StorageData) error {
	var err error
	if source == nil {
		return nil
	}
	base.global, err = mergeUniqueMaps(base.global, source.getAllGlobal())
	if err != nil {
		return err
	}
	base.local, err = mergeUniqueMaps(base.local, source.getAllLocal())
	if err != nil {
		return err
	}
	return nil
}
