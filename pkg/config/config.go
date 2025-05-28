package config

import (
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort int
}

func Load() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetDefault("SERVER_PORT", "8080")

	portStr := viper.GetString("SERVER_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerPort: port,
	}, nil
}
