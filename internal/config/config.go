package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env               string `yaml:"env" env-required:"true"`
	ConnectionStrings `yaml:"connection_strings"`
	HTTPServer        `yaml:"http_server"`
	App               `yaml:"app"`
}

type ConnectionStrings struct {
	UrlShortenerPg `yaml:"url_shortener_pg"`
}

type UrlShortenerPg struct {
	User     string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	HOST     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	PORT     uint   `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	DB       string `yaml:"db" env:"POSTGRES_DB" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env:"POSTGRES_SSLMODE" env-required:"true"`
}

type HTTPServer struct {
	Address           string        `yaml:"address" env-required:"true"`
	Timeout           time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout       time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env-default:"3s"`
}

type App struct {
	ExpiredURLPath          string        `yaml:"expired_url_path" env-default:"/error/expired"`
	GracefulShutdownTimeout time.Duration `yaml:"graceful_shutdown_timeout" env-default:"15s"`
	HashIDConfiguration     `yaml:"hash_id_configuration"`
}

type HashIDConfiguration struct {
	Salt          string `yaml:"salt" env-default:"UrlShortener"`
	Alphabet      string `yaml:"alphabet" env-default:"abcdefghijkABCDEFGHIJK12345"`
	MinHashLength int    `yaml:"min_hash_length" env-default:"5"`
}

func MustLoad() (cfg Config) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("CONFIG_PATH does not exist")
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	return
}
