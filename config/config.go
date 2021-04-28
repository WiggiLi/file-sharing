package config

import (
	"sync"
	"encoding/json"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel            string `envconfig:"LOG_LEVEL"`
	PgURL               string `envconfig:"PG_URL"`
	PgMigrationsPath    string `envconfig:"PG_MIGRATIONS_PATH"`
	HTTPAddr            string `envconfig:"HTTP_ADDR"`
	FilePath            string `envconfig:"FILE_PATH"`
}

var (
	config Config
	once sync.Once
)

//! Fatal from not main
func Get() *Config {
	once.Do(func(){
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatal(err)
		}
		configBytes, err := json.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Configuration:", string(configBytes))
	})
	return &config
}