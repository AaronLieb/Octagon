package cache

const streamMatchKey = "stream_match"

func SetStreamMatch(matchID string) error {
	return Set([]byte(streamMatchKey), []byte(matchID))
}

func GetStreamMatch() (string, error) {
	val, err := Get([]byte(streamMatchKey))
	if err != nil {
		return "", err
	}
	return string(val), nil
}
