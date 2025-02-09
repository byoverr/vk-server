package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	HTTPServer HTTPServer `yaml:"http_server"`
	Storage    Storage    `yaml:"storage"`
	Pinger     Pinger     `yaml:"pinger"`
}

type HTTPServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Storage struct {
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	DbName   string `yaml:"dbname" env:"NAME" env-default:"postgres"`
	User     string `yaml:"user" env:"USER" env-default:"user"`
	SSLMode  string `yaml:"sslmode" env:"SSL_MODE" env-default:"disable"`
	Password string `yaml:"password" env:"PASSWORD"`
}

type Pinger struct {
	Interval int `yaml:"interval" env:"INTERVAL" env-default:"10"`
	Timeout  int `yaml:"timeout" env:"TIMEOUT" env-default:"2"`
}

func Load() *Config {
	var cfg Config

	err := cleanenv.ReadConfig("../backend/config.yml", &cfg)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
