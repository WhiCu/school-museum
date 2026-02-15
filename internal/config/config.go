package config

import (
	"net"
	"time"
)

type ServerConfig struct {
	Host            string        `yaml:"host" env:"SERVER_HOST" env-default:"localhost" koanf:"host"`
	Port            string        `yaml:"port" env:"SERVER_PORT" env-default:"8080" koanf:"port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SERVER_SHUTDOWN_TIMEOUT" env-default:"30s" koanf:"shutdown_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT" env-default:"10s" koanf:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" env-default:"30s" koanf:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env:"SERVER_IDLE_TIMEOUT" env-default:"30s" koanf:"idle_timeout"`
}

type StorageConfig struct {
	Host string `yaml:"host" env:"DB_HOST" env-default:"localhost" koanf:"host"`
	Port string `yaml:"port" env:"DB_PORT" env-default:"5432" koanf:"port"`
	User string `yaml:"user" env:"DB_USER" env-default:"user" koanf:"user"`
	Pass string `yaml:"pass" env:"DB_PASS" env-default:"password" koanf:"pass"`
	Name string `yaml:"name" env:"DB_NAME" env-default:"school_museum" koanf:"name"`
}

func (s StorageConfig) DSN() string {
	return "postgres://" + s.User + ":" + s.Pass + "@" + s.Host + ":" + s.Port + "/" + s.Name + "?sslmode=disable"
}

type LoggerConfig struct {
	Level    string `yaml:"level" env:"LOG_LEVEL" env-default:"info" koanf:"level"`
	Path     string `yaml:"path" env:"LOG_PATH" env-default:"" koanf:"path"`
	Size     int    `yaml:"size" env:"LOG_FILE_SIZE" env-default:"128" koanf:"size"`
	Age      int    `yaml:"age" env:"LOG_FILE_AGE" env-default:"7" koanf:"age"`
	Backups  int    `yaml:"backups" env:"LOG_FILE_BACKUPS" env-default:"16" koanf:"backups"`
	Compress bool   `yaml:"compress" env:"LOG_COMPRESS" env-default:"true" koanf:"compress"`
}

type Config struct {
	Server    ServerConfig    `yaml:"server" env:"SERVER" koanf:"server"`
	Storage   StorageConfig   `yaml:"storage" env:"STORAGE" koanf:"storage"`
	Logger    LoggerConfig    `yaml:"logger" env:"LOGGER" koanf:"logger"`
	Telemetry TelemetryConfig `yaml:"telemetry" env:"TELEMETRY" koanf:"telemetry"`
}

type TelemetryConfig struct {
	Enabled      bool        `yaml:"enabled" koanf:"enabled"`
	OTLPEndpoint string      `yaml:"otlp_endpoint" koanf:"otlp_endpoint"`
	ServiceName  string      `yaml:"service_name" koanf:"service_name"`
	Umami        UmamiConfig `yaml:"umami" koanf:"umami"`
}

type UmamiConfig struct {
	Enabled   bool   `yaml:"enabled" koanf:"enabled"`
	URL       string `yaml:"url" koanf:"url"`
	WebsiteID string `yaml:"website_id" koanf:"website_id"`
	Username  string `yaml:"username" koanf:"username"`
	Password  string `yaml:"password" koanf:"password"`
	Domain    string `yaml:"domain" koanf:"domain"`
}

func (srv *ServerConfig) ServerAddr() string {
	return net.JoinHostPort(srv.Host, srv.Port)
}
