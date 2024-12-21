package settings

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config struct holds the overall configuration with separate settings for Blob, Key, Logger, and PKCS#11
type Config struct {
	Database      DatabaseSettings      `mapstructure:"database"`
	BlobConnector BlobConnectorSettings `mapstructure:"blob_connector"`
	KeyConnector  KeyConnectorSettings  `mapstructure:"key_connector"`
	Logger        LoggerSettings        `mapstructure:"logger"`
	PKCS11        PKCS11Settings        `mapstructure:"pkcs11"`
	Port          string                `mapstructure:"port"`
}

// InitializeConfig function to read the config, prioritize environment variables and fall back to config file
func InitializeConfig(path string) (*Config, error) {
	viper.AutomaticEnv()
	// env vars example:
	// export BLOB_CONNECTOR_CONNECTION_STRING="your_blob_connection_string"
	// export BLOB_CONNECTOR_CONTAINER_NAME="your_blob_container_name"
	// export LOGGER_LOG_LEVEL="info"
	// export LOGGER_LOG_TYPE="console"

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
