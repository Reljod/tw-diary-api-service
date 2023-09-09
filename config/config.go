package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

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
	viper.SetConfigType("yaml")
	configCliName := flag.String("config", "./config/config.yml", "File configuration")
	flag.Parse()

	configPath := filepath.Dir(*configCliName)
	configName := strings.TrimSuffix(filepath.Base(*configCliName), filepath.Ext(configPath))
	fmt.Printf("%v, %v", configName, configPath)
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	var config ConfigSchema
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	return &config
}

var Config ConfigSchema = *ReadConfig()
