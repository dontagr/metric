package main

import (
	"reflect"
	"testing"

	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/bootstrap"
)

func TestCreateApp(t *testing.T) {
	tests := []struct {
		name string
		want fx.Option
	}{
		{
			name: "create di",
			want: fx.Options(
				bootstrap.Config,
				bootstrap.Server,
				bootstrap.Route,
				bootstrap.Service,
				bootstrap.Store,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateApp(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateApp() = %v, want %v", got, tt.want)
			}
		})
	}
}
