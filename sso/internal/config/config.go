package config

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string `yaml:"env" env-default:"local"`
	GRPC     GRPCConfig
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
	Database DatabaseConfig
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname"`
	Port     string `yaml:"port"`
	SSLmode  string `yaml:"sslmode"`
	TimeZone string `yaml:"time_zone"`
}

func MustLoad() *Config {
	configPath := fethConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	path, err := filepath.Abs(configPath)
	if err != nil {
		panic("no such file exists")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}

func fethConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to yaml config")
	flag.Parse()

	if res == "" {
		err := godotenv.Load()
		if err != nil {
			panic("err loading .env files")
		}
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
