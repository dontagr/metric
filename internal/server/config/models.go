package config

const EnvStoreInterval = "STORE_INTERVAL"
const EnvFileStoragePath = "FILE_STORAGE_PATH"
const EnvRestore = "RESTORE"
const DatabaseDsn = "DATABASE_DSN"
const EnvConfig = "CONFIG"
const KEY = "KEY"
const CryptoKey = "CRYPTO_KEY"
const TrustedSubnet = "TRUSTED_SUBNET"

type Config struct {
	Log           Logging  `json:"Logging"`
	HTTPServer    string   `json:"address" env:"ADDRESS" flag:"a" validate:"required"`
	GRPCServer    string   `json:"address_grpc" env:"ADDRESS_GRPC"`
	Security      Security `json:"Security"`
	DataBase      DataBase `json:"DataBase"`
	Store         Store    `json:"Store"`
	TrustedSubnet string   `json:"trusted_subnet" env:"TRUSTED_SUBNET" flag:"t"`
	Restore       bool     `json:"restore" env:"RESTORE" flag:"r"`
	FilePath      string   `json:"store_file" env:"FILE_STORAGE_PATH" flag:"f"`
	FileName      string   `json:"FileName"`
	Interval      int      `json:"store_interval" env:"STORE_INTERVAL" flag:"i"`
	CryptoKey     string   `json:"crypto_key" env:"CRYPTO_KEY"`
	DatabaseDsn   string   `json:"database_dsn" env:"DATABASE_DSN" flag:"d"`
}

type Security struct {
	Key string `json:"HashKey" env:"KEY"`
}

type DataBase struct {
	Init bool `env-default:"false"`
}

type Logging struct {
	LogLevel string `json:"LogLevel" validate:"required"`
}

type Store struct {
	FilePerm uint32 `env-default:"0644"`
}
