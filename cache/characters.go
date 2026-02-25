package cache

import (
	"fmt"
)

func SetPlayerCharacter(playerID int, characterName string) error {
	key := fmt.Appendf(nil, "player_char:%d", playerID)
	return Set(key, []byte(characterName))
}

func GetPlayerCharacter(playerID int) (string, error) {
	key := fmt.Appendf(nil, "player_char:%d", playerID)
	val, err := Get(key)
	if err != nil {
		return "", err
	}
	return string(val), nil
}
