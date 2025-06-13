package config

import (
	"clientgo/models"
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"

	"os"
)

var AppConfig models.AppConfig // Global variable

func GetConfigurations() {
	var configPath string

	// Get the directory path of the current file
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	// Determine config path based on environment variable GO_ENV
	if os.Getenv("GO_ENV") == "local" {
		// Load config from local data directory relative to current file
		configPath = filepath.Join(basepath, "data")
	} else {
		// Load config from production path
		configPath = "/etc/config"
	}

	// Add config path for viper to look for the config file
	viper.AddConfigPath(configPath)

	// Set the config file name and type (expecting config.json)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Config file error: %v", err)
	}

	// Unmarshal config JSON into the AppConfig struct
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("Error unmarshalling config json: %v", err)
	}
}
