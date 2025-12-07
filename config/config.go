package config

import (
	"github.com/caarlos0/env/v11"
	"sync"
	"time"
)

var (
	instance *Config
	once     sync.Once
	initErr  error
)

type Config struct {
	Server     Server
	Postgresql Postgresql
	JWT        JWT
	AWS        AWS
}

type Server struct {
	Port    string `env:"SERVER_PORT"`
	GinMode string `env:"SERVER_GIN_MODE"`
}

type Postgresql struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Name     string `env:"POSTGRES_NAME"`
	SSLMode  string `env:"POSTGRES_SSLMODE"`
}

type JWT struct {
	Secret              string        `env:"JWT_SECRET"`
	ExpiresIn           time.Duration `env:"JWT_EXPIRES_IN"`
	RefreshTokenExpires time.Duration `env:"JWT_REFRESH_TOKEN_EXPIRES"`
}

type AWS struct {
	Region          string `env:"AWS_REGION"`
	AccessKeyId     string `env:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
	S3Bucket        string `env:"AWS_S3_BUCKET"`
	S3Endpoint      string `env:"AWS_S3_ENDPOINT"`
}

type Upload struct {
	Path        string `env:"UPLOAD_PATH"`
	MaxFileSize int64  `env:"UPLOAD_MAX_FILE_SIZE"`
}

func GetInstance() (*Config, error) {
	once.Do(func() {
		instance = &Config{}
		initErr = env.Parse(instance)
		if initErr != nil {
			instance = nil
		}
	})
	return instance, initErr
}
