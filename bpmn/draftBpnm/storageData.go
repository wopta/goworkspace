package draftbpnm

import (
	"errors"
	"fmt"
)

type InitializeDataTipe func() DataBpnm
type StorageData interface {
	resetLocal()
	AddLocal(string, DataBpnm) error
	AddGlobal(string, DataBpnm) error
	GetLocal(string) (DataBpnm, error)
	GetGlobal(string) (DataBpnm, error)
	GetAllLocal() map[string]any
	GetAllGlobal() map[string]any
}
type StorageBpnm struct {
	Local  map[string]any
	Global map[string]any
}

func NewStorageBpnm() *StorageBpnm {
	res := new(StorageBpnm)
	res.Local = make(map[string]any)
	res.Global = make(map[string]any)
	return res
}
func (p *StorageBpnm) resetLocal() {
	p.Local = make(map[string]any)
}
func (p *StorageBpnm) GetAllLocal() map[string]any {
	return p.Local
}
func (p *StorageBpnm) GetAllGlobal() map[string]any {
	return p.Global
}

func (p *StorageBpnm) AddLocal(name string, data DataBpnm) error {
	if p.Local == nil {
		return errors.New("error initialization local storage")
	}
	if _, ok := p.Local[name]; ok {
		return fmt.Errorf("storage has already data with name %v", name)
	}
	p.Local[name] = data
	return nil
}

func (p *StorageBpnm) AddGlobal(name string, data DataBpnm) error {
	if p.Global == nil {
		return errors.New("error initialization storage storage")
	}
	if _, ok := p.Global[name]; ok {
		return fmt.Errorf("storage has already data with name %v", name)
	}
	p.Global[name] = data
	return nil
}

func (p *StorageBpnm) GetGlobal(name string) (DataBpnm, error) {
	if p.Global == nil {
		return nil, errors.New("error initialization storage storage")
	}
	if data, ok := p.Global[name]; ok {
		return data.(DataBpnm), nil
	}
	return nil, fmt.Errorf("no data founded %v", name)
}

func (p *StorageBpnm) GetLocal(name string) (DataBpnm, error) {
	if p.Global == nil {
		return nil, errors.New("error initialization storage storage")
	}
	if data, ok := p.Local[name]; ok {
		return data.(DataBpnm), nil
	}
	return nil, fmt.Errorf("no data founded %v", name)
}
