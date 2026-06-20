package config

import "github.com/spf13/viper"

type Config struct {
	DBName     string `mapstructure:"db_name"`
	ServerPort string `mapstructure:"server_port"`
}

var AppConfig *Config

func InitConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("db_name", "todo.db")
	viper.SetDefault("server_port", "8080")

	if err := viper.ReadInConfig(); err != nil {
		viper.AutomaticEnv()
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic("Failed to unmarshal config")
	}

	AppConfig = &config
	return &config
}
