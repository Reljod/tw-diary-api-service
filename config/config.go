package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var DefaultConfig ConfigSchema = ConfigSchema{
	Database: DatabaseConfigSchema{
		Engine:   "postgres",
		Host:     "localhost",
		Db:       "postgres",
		Port:     5432,
		User:     "dummy",
		Password: "dummy",
	},
}

type ConfigSchema struct {
	Database DatabaseConfigSchema
}

type DatabaseConfigSchema struct {
	Engine   string
	Host     string
	Port     int32
	Db       string
	User     string
	Password string
}

func ReadConfig() *ConfigSchema {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		return &DefaultConfig
	}

	var config ConfigSchema
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		return &DefaultConfig
	}

	return &config
}

var Config ConfigSchema = *ReadConfig()
