package config

import "github.com/spf13/viper"

type Config struct {
	DBHost     string `mapstructure:"db_host"`
	DBPort     string `mapstructure:"db_port"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
	DBName     string `mapstructure:"db_name"`
	ServerPort string `mapstructure:"server_port"`
}

var AppConfig *Config

func InitConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_port", "3306")
	viper.SetDefault("db_user", "root")
	viper.SetDefault("db_password", "password")
	viper.SetDefault("db_name", "todo_db")
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
