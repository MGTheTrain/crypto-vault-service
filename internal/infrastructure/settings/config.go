package settings

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config struct holds the overall configuration with separate settings for Blob, Key, Logger, and PKCS#11
type Config struct {
	BlobConnector BlobConnectorSettings `mapstructure:"blob_connector"`
	KeyConnector  KeyConnectorSettings  `mapstructure:"key_connector"`
	Logger        LoggerSettings        `mapstructure:"logger"`
	PKCS11        PKCS11Settings        `mapstructure:"pkcs11"`
}

// InitializeConfig function to read the config YAML file and unmarshal it into the struct using Viper
func InitializeConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("unable to read config file, %v", err)
	}

	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return &config, nil
}
