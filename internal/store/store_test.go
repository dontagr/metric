package store

import (
	"reflect"
	"testing"

	"github.com/go-playground/assert/v2"

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
