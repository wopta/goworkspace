package draftbpnm

import (
	"errors"
	"fmt"
)

type InitializeDataTipe func() DataBpnm
type StorageData interface {
	ResetLocal()
	ResetGlobal()
	AddLocal(string, DataBpnm) error
	AddGlobal(string, DataBpnm) error
	GetLocal(string) (DataBpnm, error)
	GetGlobal(string) (DataBpnm, error)
	GetAllLocal() map[string]any
	GetAllGlobal() map[string]any
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
