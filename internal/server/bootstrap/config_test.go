package bootstrap

import (
	"os"
	"reflect"
	"testing"

	server "github.com/dontagr/metric/internal/server/config"
)

type callbackFunc func(t *testing.T, fileName string)

func Test_newConfig(t *testing.T) {
	tests := []struct {
		fnc         callbackFunc
		want        *server.Config
		name        string
		cnfFileName string
		wantErr     bool
	}{
		{
			name:        "Valid Configuration",
			fnc:         setEnv,
			cnfFileName: "server.json",
			want: &server.Config{
				HTTPServer: server.HTTPServer{
					BindAddress: ":8080",
				},
				Log: server.Logging{
					LogLevel: "INFO",
				},
				Store: server.Store{
					Interval: 0,
					FileName: "bochenok_s_medom",
					FilePath: "./",
					Restore:  true,
					FilePerm: 420,
				},
				DataBase: server.DataBase{
					DatabaseDsn: "postgres://postgres:postgres@localhost:5432/metrics",
				},
				Security: server.Security{
					Key: "",
				},
			},
			wantErr: false,
		},
		{
			name:        "Invalid Configuration",
			fnc:         setEnv,
			cnfFileName: "config.json",
			want:        nil,
			wantErr:     true, // Предполагаем, что будет ошибка, например из-за отсутствия обязательного параметра
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fnc(t, tt.cnfFileName)
			got, err := newConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("newConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}

			tt.want.Log.LogLevel = got.Log.LogLevel
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func setEnv(t *testing.T, fileName string) {
	err := os.Setenv("CONFIG_FILE_NAME", fileName)
	if err != nil {
		t.Errorf("setting env CONFIG_FILE_NAME failed: %v", err)
	}
}
