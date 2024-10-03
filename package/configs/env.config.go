package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB_HOST            string `mapstructure:"DB_HOST"`
	DB_PORT            string `mapstructure:"DB_PORT"`
	DB_USER            string `mapstructure:"DB_USER"`
	DB_PASSWORD        string `mapstructure:"DB_PASSWORD"`
	DB_NAME            string `mapstructure:"DB_NAME"`
	SERVER_PORT        string `mapstructure:"SERVER_PORT"`
	TELEGRAM_BOT_TOKEN string `mapstructure:"TELEGRAM_BOT_TOKEN"`
	CHANNEL_ID         int64  `mapstructure:"CHANNEL_ID"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
