package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-default:"true"`
	GRPC        GRPCConfig
}
type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout" env-default:"5h"`
}

func MustLoad() *Config {
	configPath := fethConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
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
