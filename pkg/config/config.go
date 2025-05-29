package config

import (
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  int
	RedisDB    int
}

func Load() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "url_shortener")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_DB", "0")

	port, err := strconv.Atoi(viper.GetString("SERVER_PORT"))
	if err != nil {
		return nil, err
	}

	dbPort, err := strconv.Atoi(viper.GetString("DB_PORT"))
	if err != nil {
		return nil, err
	}

	redisPort, err := strconv.Atoi(viper.GetString("REDIS_PORT"))
	if err != nil {
		return nil, err
	}

	redisDB, err := strconv.Atoi(viper.GetString("REDIS_DB"))
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerPort: port,
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     dbPort,
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		RedisHost:  viper.GetString("REDIS_HOST"),
		RedisPort:  redisPort,
		RedisDB:    redisDB,
	}, nil
}
