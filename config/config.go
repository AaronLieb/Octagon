package config

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

const (
	DEFAULT_CONFIG_PATH = "/.config/octagon/"
	DEFAULT_ENV_NAME    = "octagonrc"
)

func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Warnf("unable to find user home directory")
	}
	return homeDir + DEFAULT_CONFIG_PATH
}

func Load() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Warnf("unable to find user home directory")
	}

	err = godotenv.Load(homeDir + DEFAULT_CONFIG_PATH + DEFAULT_ENV_NAME)
	if err != nil {
		log.Warnf("unable to load default configuration: %v", err)
	} else {
		log.Debug("successfully loaded config from %s", DEFAULT_CONFIG_PATH)
		return
	}
	err = godotenv.Load()
	if err != nil {
		log.Warnf("unable to load configuration from local .env: %v", err)
		log.Fatal("unable to load a configuration file")
	}
	log.Debug("successfully loaded config from .env")
}
