package config

const EnvStoreInterval = "STORE_INTERVAL"
const EnvFileStoragePath = "FILE_STORAGE_PATH"
const EnvRestore = "RESTORE"
const DatabaseDsn = "DATABASE_DSN"

type Config struct {
	Log        Logging    `json:"Logging"`
	HTTPServer HTTPServer `json:"HttpServing"`
	Store      Store      `json:"Store"`
	DataBase   DataBase   `json:"DataBase"`
}

type DataBase struct {
	DatabaseDsn string `json:"DatabaseDsn" env:"DATABASE_DSN" flag:"d"`
	Init        bool   `env-default:"false"`
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
