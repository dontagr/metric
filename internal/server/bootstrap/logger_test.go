package bootstrap

import (
	"testing"

	"go.uber.org/fx"

	server "github.com/dontagr/metric/internal/server/config"
	mocks "github.com/dontagr/metric/test/mock"
)

func Test_newLogger(t *testing.T) {
	type args struct {
		lc  fx.Lifecycle
		cnf *server.Config
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			name: "Valid log level",
			args: args{
				lc:  mocks.NewMockLifecycle(),
				cnf: &server.Config{Log: server.Logging{LogLevel: "INFO"}},
			},
			wantErr: false,
		},
		{
			name: "Invalid log level",
			args: args{
				lc:  mocks.NewMockLifecycle(),
				cnf: &server.Config{Log: server.Logging{LogLevel: "invalidLevel"}},
			},
			wantErr: true,
		},
		{
			name: "Debug log level",
			args: args{
				lc:  mocks.NewMockLifecycle(),
				cnf: &server.Config{Log: server.Logging{LogLevel: "DEBUG"}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newLogger(tt.args.lc, tt.args.cnf)
			if (err != nil) != tt.wantErr {
				t.Errorf("newLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("newLogger() expected a logger, got nil")
			}
		})
	}
}
