package config

import (
	"encoding/json"
	"os"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
)

const AliasFile = "aliases.json"

var aliasConfig map[string]string

func LoadAliases() error {
	if aliasConfig != nil {
		return nil
	}

	path := GetConfigPath() + AliasFile
	data, err := os.ReadFile(path)
	if err != nil {
		log.Warn("unable to read aliases file: %v", err)
		aliasConfig = make(map[string]string)
		return nil
	}

	aliasConfig = make(map[string]string)
	err = json.Unmarshal(data, &aliasConfig)
	if err != nil {
		log.Errorf("unable to unmarshal aliases file: %v", err)
		aliasConfig = make(map[string]string)
		return err
	}

	log.Debug("Reading aliases", "path", path, "n", len(aliasConfig))
	return nil
}

func ResolveAlias(userID startgg.ID) startgg.ID {
	if aliasConfig == nil {
		err := LoadAliases()
		if err != nil {
			log.Error("unable to load aliases", "err", err)
			return userID
		}
	}

	if aliasConfig == nil {
		return userID
	}

	userIDStr := startgg.ToString(userID)
	if realID, exists := aliasConfig[userIDStr]; exists {
		return startgg.ToID(realID)
	}

	return userID
}
