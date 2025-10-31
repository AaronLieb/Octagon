/* Package config handles configuration of the `octagon` command */
package config

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

const (
	DefaultConfigPath = "/.config/octagon/"
	DefaultEnvName    = "octagonrc"
)

func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Warnf("unable to find user home directory")
	}
	return homeDir + DefaultConfigPath
}

func Load() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Warnf("unable to find user home directory")
	}

	err = godotenv.Load(homeDir + DefaultConfigPath + DefaultEnvName)
	if err != nil {
		log.Warnf("unable to load default configuration: %v", err)
	} else {
		log.Debugf("successfully loaded config from %s", DefaultConfigPath)
		return
	}
	err = godotenv.Load()
	if err != nil {
		log.Warnf("unable to load configuration from local .env: %v", err)
		log.Fatal("unable to load a configuration file")
	}
	log.Debug("successfully loaded config from .env")
}
