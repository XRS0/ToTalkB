package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server              ServerConfig              `mapstructure:"server_event-manager"`
	Database            DatabaseConfig            `mapstructure:"database"`
	NotificationService NotificationServiceConfig `mapstructure:"notification_service"`
}

type ServerConfig struct {
	Port     int    `mapstructure:"port"`
	GRPCPort int    `mapstructure:"grpc_port"`
	Host     string `mapstructure:"host"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type NotificationServiceConfig struct {
	Host     string `mapstructure:"host"`
	GRPCPort int    `mapstructure:"grpc_port"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// fmt.Printf("%+v", config)

	return &config, nil
}
