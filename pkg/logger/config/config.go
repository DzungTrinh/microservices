package config

type Config struct {
	LogLevel string `env:"LOGGER_LEVEL" env-required:"true"`
	DevMode  bool   `env:"LOGGER_DEV_MODE" env-required:"true"`
	Encoder  string `env:"LOGGER_ENCODER" env-required:"true"`
	Path     string `env:"LOGGER_PATH"`
}
