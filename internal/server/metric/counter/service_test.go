package counter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dontagr/metric/internal/server/metric/factory"
	"github.com/dontagr/metric/models"
)

func TestRegisterMetric(t *testing.T) {
	type args struct {
		mf *factory.MetricFactory
	}
	tests := []struct {
		args args
		name string
	}{
		{
			name: "еще один тест ((",
			args: args{mf: factory.NewMetricFactory()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterMetric(tt.args.mf)

			_, err := tt.args.mf.GetMetric(models.Counter)
			assert.NoError(t, err)
		})
	}
}

func TestMetric_GetName(t *testing.T) {
	type fields struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test name",
			fields: fields{
				name: models.Counter,
			},
			want: models.Counter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metric{
				name: tt.fields.name,
			}
			assert.Equalf(t, tt.want, m.GetName(), "GetName()")
		})
	}
}

func TestMetric_ConvertToMetrics(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		id    string
		value string
	}
	testInt := int64(118)
	tests := []struct {
		want    *models.Metrics
		wantErr assert.ErrorAssertionFunc
		args    args
		name    string
		fields  fields
	}{
		{
			name: "test convert",
			fields: fields{
				name: models.Counter,
			},
			args: args{
				id:    "test",
				value: "118",
			},
			want: &models.Metrics{
				ID:    "test",
				MType: models.Counter,
				Delta: &testInt,
			},
			wantErr: assert.NoError,
		},
		{
			name: "test convert",
			fields: fields{
				name: models.Gauge,
			},
			args: args{
				id:    "test",
				value: "",
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "test invalid id",
			fields: fields{
				name: models.Counter,
			},
			args: args{
				id:    "",
				value: "118",
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metric{
				name: tt.fields.name,
			}
			got, err := m.ConvertToMetrics(tt.args.id, tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("ConvertToMetrics(%v, %v)", tt.args.id, tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ConvertToMetrics(%v, %v)", tt.args.id, tt.args.value)
		})
	}
}

func TestMetric_Process(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		in0 *models.Metrics
		in1 *models.Metrics
	}
	testInt := int64(118)
	tests := []struct {
		args    args
		wantErr assert.ErrorAssertionFunc
		name    string
		fields  fields
		want    int64
	}{
		{
			name: "test processing",
			fields: fields{
				name: models.Counter,
			},
			args: args{
				in0: &models.Metrics{
					ID:    "test",
					MType: models.Counter,
					Delta: &testInt,
				},
				in1: &models.Metrics{
					ID:    "test",
					MType: models.Counter,
					Delta: &testInt,
				},
			},
			wantErr: assert.NoError,
			want:    int64(236),
		},
		{
			name: "test nil oldMetrics",
			fields: fields{
				name: models.Counter,
			},
			args: args{
				in0: nil,
				in1: &models.Metrics{
					ID:    "test",
					MType: models.Counter,
					Delta: &testInt,
				},
			},
			wantErr: assert.NoError,
			want:    int64(236),
		},
		{
			name: "test nil newMetrics",
			fields: fields{
				name: models.Counter,
			},
			args: args{
				in0: &models.Metrics{
					ID:    "test",
					MType: models.Counter,
					Delta: &testInt,
				},
				in1: nil,
			},
			wantErr: assert.Error,
			want:    int64(236),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metric{
				name: tt.fields.name,
			}

			err := m.Process(tt.args.in0, tt.args.in1)
			if tt.wantErr(t, err, fmt.Sprintf("Process(%v, %v)", tt.args.in0, tt.args.in1)) {
				return
			}

			assert.Equal(t, *tt.args.in1.Delta, int64(236))
		})
	}
}

func TestMetric_ReturnValue(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		metrics *models.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Counter metric with delta",
			fields: fields{
				name: models.Counter,
			},
			args: args{metrics: &models.Metrics{
				ID:    "test",
				MType: models.Counter,
				Delta: ptrInt64(100),
			}},
			want: "100",
		},
		{
			name: "Counter metric with nil delta",
			fields: fields{
				name: models.Counter,
			},
			args: args{metrics: &models.Metrics{
				ID:    "test",
				MType: models.Counter,
				Delta: nil,
			}},
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metric{
				name: tt.fields.name,
			}
			got := m.ReturnValue(tt.args.metrics)
			assert.Equalf(t, tt.want, got, "ReturnValue(%v)", tt.args.metrics)
		})
	}
}
func ptrInt64(i int64) *int64 { return &i }

func TestMetric_GetMetricsByData(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		id    string
		value any
	}
	testInt := int64(150)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Metrics
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid Counter with int64 value",
			fields: fields{
				name: models.Counter,
			},
			args: args{
				id:    "testCounter",
				value: testInt,
			},
			want: &models.Metrics{
				ID:    "testCounter",
				MType: models.Counter,
				Delta: &testInt,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid value type for Counter",
			fields: fields{
				name: models.Counter,
			},
			args: args{
				id:    "testCounter",
				value: "notAnInt",
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metric{
				name: tt.fields.name,
			}
			got, err := m.GetMetricsByData(tt.args.id, tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMetricsByData(%v, %v)", tt.args.id, tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMetricsByData(%v, %v)", tt.args.id, tt.args.value)
		})
	}
}
