package draftbpnm

import (
	"errors"
	"fmt"
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
	// move/overwrite all data from source -> base
	Merge(StorageData) error
}
type StorageBpnm struct {
	local  map[string]any
	global map[string]any
}

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
	return nil, fmt.Errorf("no data founded %v", name)
}

func (p *StorageBpnm) GetLocal(name string) (DataBpnm, error) {
	if p.global == nil {
		return nil, errors.New("error initialization storage storage")
	}
	if data, ok := p.local[name]; ok {
		return data.(DataBpnm), nil
	}
	return nil, fmt.Errorf("no data founded %v", name)
}

// move/overwrite all data from source -> base
func (base *StorageBpnm) Merge(source StorageData) error {
	if source == nil {
		return nil
	}
	base.global = mergeMaps(base.global, source.GetAllGlobal())
	base.local = mergeMaps(base.local, source.GetAllLocal())
	return nil
}

func GetData[t DataBpnm](name string, storage StorageData) (t, error) {
	data, err := storage.GetGlobal("policy")
	var result t
	if err != nil {
		return *new(t), err
	}

	result = data.(t)
	if data.Type() == result.Type() {
		return *new(t), fmt.Errorf("Data '%v' with type %v founded has a different type than '%v'", name, result.Type(), data.Type())
	}
	return result, nil
}
