package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	App_env           string `env:"app_env,required" envDefault:"DEVELOPMENT"`
	Token             string `env:"token"`
	Debug_addr        string `env:"debug_addr,required,notEmpty" envDefault:":8080"`
	Http_addr         string `env:"http_addr,required,notEmpty" envDefault:":8081"`
	Channel           string `env:"gokit-films"`
	Postgres_database string `env:"postgres_database,required,notEmpty" envDefault:"dvdrental"`
	Postgres_dbuser   string `env:"postgres_dbuser,required,notEmpty" envDefault:"postgres"`
	Postgres_host     string `env:"postgres_host,required,notEmpty" envDefault:"localhost"`
	Postgres_port     string `env:"postgres_port,required,notEmpty" envDefault:"5432"`
	Postgres_password string `env:"postgres_password,required,notEmpty,unset" envDefault:"secret"`
	Postgres_sslmode  string `env:"postgres_sslmode,required,notEmpty" envDefault:"disable"`
	Postgres_timezone string `env:"postgres_timezone,required,notEmpty" envDefault:"Europe/Prague"`
}

func InitConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
