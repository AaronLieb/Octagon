package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
)

const BiasFile = "bias.json"

type Bias struct {
	PlayerID   string     `json:"playerId"`
	Ratio      float64    `json:"ratio"`
	Reason     string     `json:"reason"`
	Expiration *time.Time `json:"expiration,omitempty"`
}

func GetBiases() []Bias {
	return readBiasFile()
}

func SaveBias(newBias Bias) error {
	biases := readBiasFile()
	biases = removeExpiredBiases(biases)
	biases = append(biases, newBias)
	return writeBiasFile(biases)
}

func GetBiasForPlayer(userID startgg.ID) float64 {
	biases := readBiasFile()
	biases = removeExpiredBiases(biases)

	userIDStr := startgg.ToString(userID)
	for _, bias := range biases {
		if bias.PlayerID == userIDStr {
			return bias.Ratio
		}
	}
	return 1.0
}

func readBiasFile() []Bias {
	var biases []Bias
	path := GetConfigPath() + BiasFile

	file, err := os.ReadFile(path)
	if err != nil {
		log.Warn("unable to read bias file: %v", err)
		return biases
	}

	if err := json.Unmarshal(file, &biases); err != nil {
		log.Errorf("unable to unmarshal bias file: %v", err)
		return biases
	}

	log.Debug("Reading biases", "path", path, "n", len(biases))
	return biases
}

func writeBiasFile(biases []Bias) error {
	path := GetConfigPath() + BiasFile

	data, err := json.MarshalIndent(biases, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func removeExpiredBiases(biases []Bias) []Bias {
	now := time.Now()
	var active []Bias

	for _, bias := range biases {
		if bias.Expiration == nil || bias.Expiration.After(now) {
			active = append(active, bias)
		}
	}

	return active
}
