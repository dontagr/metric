package store

import (
	"fmt"

	"github.com/dontagr/metric/internal/server/service/intersaces"
)

type StoreFactory struct {
	Collection map[string]intersaces.Store
}

func NewStoreFactory() *StoreFactory {
	return &StoreFactory{
		Collection: make(map[string]intersaces.Store),
	}
}

func (f *StoreFactory) GetStore(name string) (intersaces.Store, error) {
	if val, ok := f.Collection[name]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("story with name %s not found", name)
}

func (f *StoreFactory) SetStory(s intersaces.Store) {
	f.Collection[s.GetName()] = s
}
