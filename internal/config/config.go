package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	configFileName = "config"
	cwdConfigPath  = "./"
	homeConfigPath = "$HOME/.maladapt"
)

var ErrConfigFileNotFound = errors.New("Unable to Find Config file.")

type Config struct {
	BindAddress    string
	QuarantinePath string
	MaxUploadSize  int64
}

func Initialize() *Config {
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(cwdConfigPath)
	viper.AddConfigPath(homeConfigPath)

	viper.SetDefault("maladapt.bind_address", "localhost:3333")
	viper.SetDefault("maladapt.max_upload_size", 24000)
	viper.SetDefault("maladapt.quarantine_path", "$HOME/.maladapt/quarantined")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			log.Fatal(viper.ConfigParseError{})
		}
		log.Fatal(ErrConfigFileNotFound)
	}

	return NewConfig(
		viper.GetString("maladapt.bind_address"),
		viper.GetString("maladapt.quarantine_path"),
		viper.GetInt64("maladapt.max_upload_size")+int64(10<<20), // Reserve an additional 10 MB for non-file parts.
	)
}

func (c *Config) Validate() error {
	return os.MkdirAll(c.QuarantinePath, 0770)
}

func NewConfig(bindAddress string, quarantinePath string, maxUploadSize int64) *Config {
	return &Config{
		BindAddress:    bindAddress,
		QuarantinePath: quarantinePath,
		MaxUploadSize:  maxUploadSize,
	}
}
