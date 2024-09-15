package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	NetworkConfig NetworkConfig `mapstructure:"networkConfig"`
	APIConfig     APIConfig     `mapstructure:"apiConfig"`
	// Add other configuration structs as needed
}

type NetworkConfig struct {
	ListenAddr     string   `json:"listen_addr"`
	BootstrapPeers []string `json:"bootstrap_peers"`
	Port     int    `mapstructure:"port"`
	Protocol string `mapstructure:"protocol"`
}

type APIConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

var globalConfig *Config

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")       // Name of config file (without extension)
	viper.SetConfigType("yaml")         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")            // Look for config in the working directory
	viper.AddConfigPath("./config")     // Look for config in a subdirectory named "config"
	viper.AddConfigPath("../config")    // Look for config in the parent directory's config folder

	err := viper.ReadInConfig()         // Find and read the config file
	if err != nil {                     // Handle errors reading the config file
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetNetworkConfig() NetworkConfig {
	if globalConfig == nil {
		var err error
		globalConfig, err = LoadConfig()
		if err != nil {
			// Handle error, perhaps log it and return a default config
			return NetworkConfig{}
		}
	}
	return globalConfig.NetworkConfig
}
