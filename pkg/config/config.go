package config

import (
  "github.com/spf13/viper"
)

type Config struct {
  ServerPort int    `mapstructure:"SERVER_PORT"`
  DBURL      string `mapstructure:"DB_URL"`
  RedisURL   string `mapstructure:"REDIS_URL"`
}

func Load() (*Config, error) {
  viper.SetDefault("SERVER_PORT", 8080)
  viper.SetDefault("DB_URL", "postgres://user:pass@localhost:5432/dbname?sslmode=disable")
  viper.SetDefault("REDIS_URL", "redis://localhost:6379/0")

  viper.AutomaticEnv()

  var cfg Config
  if err := viper.Unmarshal(&cfg); err != nil {
    return nil, err
  }
  return &cfg, nil
}