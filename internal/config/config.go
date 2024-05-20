package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

const (
	dotenvPath = ".env"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	GRPC     GRPCConfig    `yaml:"grpc"`
	Storage  StorageConfig `yaml:"storage"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type StorageConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Dbname   string `yaml:"dbname"`
	SslMode  string `yaml:"sslMode"`
	Password string `yaml:"password"`
}

func MustLoad() *Config {
	if err := godotenv.Load(dotenvPath); err != nil {
		panic(err)
	}
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}
	cfg := new(Config)
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return cfg
}

func fetchConfigPath() string {
	var res string
	/**
	Function to get config path from flag(or if there is no flag from .env)
	Example: --config="path/to/config.yaml"
	**/
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("LOCAL_CONFIG")
	}
	return res
}
