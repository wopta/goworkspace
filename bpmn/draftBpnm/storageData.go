package draftbpnm

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
	GetAllLocal() map[string]any
	GetAllGlobal() map[string]any
	// setHigherStorage sets a higher-level storage,
	// which will be used as a fallback when a local/global key is not found in the current storage.
	setHigherStorage(StorageData) error
	// It merges two unique storage
	// If both storage contain the same key, return error
	MergeUnique(StorageData) error
	//Mark what local resource keep when clean is called
	markWhatNeeded([]TypeData)
	//Delete the resource that aren't needed(aren't marked)
	clean()
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
	return res
}

func (p *StorageBpnm) ResetLocal() {
	p.local = make(map[string]any)
}

func (p *StorageBpnm) ResetGlobal() {
	p.global = make(map[string]any)
}

func (p *StorageBpnm) markWhatNeeded(toTouch []TypeData) {
	for _, t := range toTouch {
		p.touched = append(p.touched, t.Name)
	}
}

func (p *StorageBpnm) clean() {
	backup := maps.Clone(p.local)
	p.ResetLocal()
	for i, toSave := range p.touched {
		p.AddLocal(p.touched[i], backup[toSave].(DataBpnm))
	}
	p.touched = nil
}

func (p *StorageBpnm) GetAllLocal() map[string]any {
	return p.local
}

func (p *StorageBpnm) GetAllGlobal() map[string]any {
	return p.global
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

// copy all data from source -> base, if both have same key return error
func (base *StorageBpnm) MergeUnique(source StorageData) error {
	var err error
	if source == nil {
		return nil
	}
	base.global, err = mergeUniqueMaps(base.global, source.GetAllGlobal())
	if err != nil {
		return err
	}
	base.local, err = mergeUniqueMaps(base.local, source.GetAllLocal())
	if err != nil {
		return err
	}
	return nil
}

func GetData[t DataBpnm](name string, storage StorageData) (t, error) {
	data, err := storage.GetLocal(name)
	var result t
	if err != nil {
		data, err = storage.GetGlobal(name)
	}
	if err != nil {
		return *new(t), err
	}

	result = data.(t)
	if data.GetType() != result.GetType() {
		return *new(t), fmt.Errorf("Data '%v' with type %v founded has a different type than '%v'", name, result.GetType(), data.GetType())
	}
	return result, nil
}
