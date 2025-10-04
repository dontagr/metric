package store

import (
	"reflect"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"

	"github.com/dontagr/metric/internal/server/service/interfaces"
	"github.com/dontagr/metric/models"
)

func TestMemStorage_SaveMetric(t *testing.T) {
	type fields struct {
		collection map[string]*models.Metrics
	}
	type args struct {
		metrics *models.Metrics
	}
	testFloat := 1.18
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			name:   "save test",
			fields: fields{collection: make(map[string]*models.Metrics)},
			args:   args{metrics: &models.Metrics{ID: "test", MType: models.Gauge, Value: &testFloat}},
			want: fields{collection: map[string]*models.Metrics{
				"gauge_test": {ID: "test", MType: models.Gauge, Value: &testFloat},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				collection: tt.fields.collection,
			}
			err := m.SaveMetric(tt.args.metrics)
			if err != nil {
				return
			}
			assert.Equal(t, m.collection, tt.want.collection)
		})
	}
}

func TestMemStorage_LoadMetric(t *testing.T) {
	type fields struct {
		collection map[string]*models.Metrics
	}
	type args struct {
		id    string
		mType string
	}
	testFloat := 1.18
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *models.Metrics
	}{
		{
			name: "load test",
			fields: fields{
				collection: map[string]*models.Metrics{
					"gauge_test": {ID: "test", MType: models.Gauge, Value: &testFloat},
				},
			},
			args: args{
				id:    "test",
				mType: models.Gauge,
			},
			want: &models.Metrics{ID: "test", MType: models.Gauge, Value: &testFloat},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				collection: tt.fields.collection,
			}
			if got, _ := m.LoadMetric(tt.args.id, tt.args.mType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		{
			name: "test storage",
			want: &MemStorage{
				collection: make(map[string]*models.Metrics),
				name:       models.StoreMem,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newMemStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterStoreMem(t *testing.T) {
	ms := NewMockStoreFactory()

	expectedStore := newMemStorage()
	ms.On("GetStore", models.StoreMem).Return(expectedStore, nil)

	RegisterStoreMem(ms)

	store, err := ms.GetStore(models.StoreMem)
	if err != nil {
		t.Errorf("Expected store to be registered, but got error: %v", err)
	}

	memStore, ok := store.(*MemStorage)
	if !ok {
		t.Errorf("Expected store to be of type *MemStorage, but got %T", store)
	}

	if memStore.GetName() != models.StoreMem {
		t.Errorf("Expected store name to be %s, but got %s", models.StoreMem, memStore.GetName())
	}

	ms.AssertCalled(t, "GetStore", models.StoreMem)
}

type MockStoreFactory struct {
	mock.Mock
}

func NewMockStoreFactory() *MockStoreFactory {
	return &MockStoreFactory{}
}

func (m *MockStoreFactory) GetStore(name string) (interfaces.Store, error) {
	args := m.Called(name)
	return args.Get(0).(interfaces.Store), args.Error(1)
}

func (m *MockStoreFactory) SetStory(s interfaces.Store) {
	//m.Called(s.GetName())
}
