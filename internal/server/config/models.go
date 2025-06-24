package config

const EnvStoreInterval = "STORE_INTERVAL"
const EnvFileStoragePath = "FILE_STORAGE_PATH"
const EnvRestore = "RESTORE"

type Config struct {
	Log        Logging    `json:"Logging"`
	HTTPServer HTTPServer `json:"HttpServing"`
	Store      Store      `json:"Store"`
}

type HTTPServer struct {
	BindAddress string `json:"BindAddress" env:"ADDRESS" flag:"a" validate:"required"`
}

type Logging struct {
	LogLevel string `json:"LogLevel" validate:"required"`
}

type Store struct {
	Interval int    `json:"Interval" env:"STORE_INTERVAL" flag:"i"`
	FilePath string `json:"FilePath" env:"FILE_STORAGE_PATH" flag:"f"`
	FileName string `json:"FileName"`
	FilePerm uint32 `env-default:"0644"`
	Restore  bool   `json:"Restore" env:"RESTORE" flag:"r"`
}
