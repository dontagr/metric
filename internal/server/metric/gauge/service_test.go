package gauge

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

			_, err := tt.args.mf.GetMetric(models.Gauge)
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
				name: models.Gauge,
			},
			want: models.Gauge,
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
	testFloat := 1.18
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
				name: models.Gauge,
			},
			args: args{
				id:    "test",
				value: "1.18",
			},
			want: &models.Metrics{
				ID:    "test",
				MType: models.Gauge,
				Value: &testFloat,
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
			name: "test convert",
			fields: fields{
				name: models.Gauge,
			},
			args: args{
				id:    "",
				value: "",
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
	tests := []struct {
		args    args
		wantErr assert.ErrorAssertionFunc
		name    string
		fields  fields
	}{
		{
			name: "test processing",
			fields: fields{
				name: models.Gauge,
			},
			args: args{
				in0: nil,
				in1: nil,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metric{
				name: tt.fields.name,
			}
			tt.wantErr(t, m.Process(tt.args.in0, tt.args.in1), fmt.Sprintf("Process(%v, %v)", tt.args.in0, tt.args.in1))
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
			name: "Gauge metric with value",
			fields: fields{
				name: models.Gauge,
			},
			args: args{metrics: &models.Metrics{
				ID:    "test",
				MType: models.Gauge,
				Value: ptrFloat64(25.5),
			}},
			want: "25.5",
		},
		{
			name: "Gauge metric with nil value",
			fields: fields{
				name: models.Gauge,
			},
			args: args{metrics: &models.Metrics{
				ID:    "test",
				MType: models.Gauge,
				Value: nil,
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
func ptrFloat64(i float64) *float64 { return &i }

func TestMetric_GetMetricsByData(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		id    string
		value any
	}
	testInt := float64(25.5)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Metrics
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid Gauge with float64 value",
			fields: fields{
				name: models.Gauge,
			},
			args: args{
				id:    "testGauge",
				value: testInt,
			},
			want: &models.Metrics{
				ID:    "testGauge",
				MType: models.Gauge,
				Value: &testInt,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid value type for Gauge",
			fields: fields{
				name: models.Gauge,
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
