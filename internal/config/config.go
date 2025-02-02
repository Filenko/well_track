package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
	"os"
)

type Config struct {
	Env      string   `yaml:"env" env-default:"local"`
	Telegram TgConfig `yaml:"telegram"`
	Postgres PgConfig `yaml:"postgres"`
	Redis    Redis    `yaml:"redis"`
	RabbitMQ RabbitMQ `yaml:"rabbitmq"`
}

type PgConfig struct {
	Host          string `yaml:"host" env:"PG_HOST" env-default:"localhost"`
	Port          int    `yaml:"port" env:"PG_PORT" env-default:"5432"`
	User          string `yaml:"user" env:"PG_USER" env-default:"welltrack"`
	Password      string `yaml:"password" env:"PG_PASSWORD" env-default:"welltrack"`
	DBName        string `yaml:"database" env:"PG_DB_NAME" env-default:"welltrack"`
	MigrationsDir string `yaml:"migrations_dir" env:"PG_MIGRATIONS_DIR" env-default:"/app/migrations"`
}

type Redis struct {
	Address string `yaml:"address" env:"REDIS_ADDR" env-default:"localhost:6379"`
}

type RabbitMQ struct {
	Host     string `yaml:"host" env:"RABBIT_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"RABBIT_PORT" env-default:"5672"`
	User     string `yaml:"user" env:"RABBIT_USER" env-default:"guest"`
	Password string `yaml:"password" env:"RABBIT_PASS" env-default:"guest"`
}

type TgConfig struct {
	TelegramApiToken string `yaml:"telegram_api_token" env:"TELEGRAM_API_TOKEN"`
}

func MustLoad(log *zerolog.Logger) *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal().Msg("CONFIG_PATH is not set")
	}

	// Проверяем, существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal().Str("path", configPath).Msg("config file does not exist")
	}

	var cfg Config

	// Считываем из переменных окружения
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Error().Msg("cannot read environment variables")
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal().Err(err).Msg("cannot read config file")
	}

	return &cfg
}
