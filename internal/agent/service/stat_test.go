package service

import (
	"reflect"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNewStats(t *testing.T) {
	tests := []struct {
		want *Stats
		name string
	}{
		{
			name: "проверка структуры",
			want: &Stats{
				UpdateWg: NewStatUpdateWg(),
				SendWg:   NewStatSendWg(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStats(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStats_Update(t *testing.T) {
	type fields struct {
		PollCount int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "проверка PollCount",
			fields: struct {
				PollCount int
			}{PollCount: 1},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stats{
				PollCount: tt.fields.PollCount,
				UpdateWg:  NewStatUpdateWg(),
				SendWg:    NewStatSendWg(),
			}
			s.Update()
			s.Update()
			s.Update()
			s.Update()
			s.Update()

			assert.Equal(t, s.PollCount, tt.want)
		})
	}
}
