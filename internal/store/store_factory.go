package store

import (
	"fmt"

	"github.com/dontagr/metric/internal/server/service/interfaces"
)

type StoreFactory struct {
	Collection map[string]interfaces.Store
}

func NewStoreFactory() interfaces.IStoreFactory {
	return &StoreFactory{
		Collection: make(map[string]interfaces.Store),
	}
}

func (f *StoreFactory) GetStore(name string) (interfaces.Store, error) {
	if val, ok := f.Collection[name]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("story with name %s not found", name)
}

func (f *StoreFactory) SetStory(s interfaces.Store) {
	f.Collection[s.GetName()] = s
}
