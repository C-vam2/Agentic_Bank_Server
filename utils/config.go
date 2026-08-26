package utils

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configurations of the application.
// The values are read by Viper from config file or environment variables
type Config struct {
	DBDriver       string        `mapstructure:"DB_DRIVER"`
	DBSource       string        `mapstructure:"DB_SOURCE"`
	ServerAddress  string        `mapstructure:"SERVER_ADDRESS"`
	SymmetricKey   string        `mapstructure:"SYMMETRIC_KEY"`
	AccessDuration time.Duration `mapstructure:"ACCESS_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
