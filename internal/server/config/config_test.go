package config

import (
	"reflect"
	"testing"
)

func Test_loadConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{
			name: "корректная загрузка",
			want: &Config{
				HTTPServer{
					BindAddress: ":8080",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getConfigFile(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "success case",
			want: []string{"../../../configs/server.json", "./configs/server.json"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConfigFile(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConfigFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
