package startgg

import (
	"encoding/json"
	"strconv"
)

type ID any

type CachedPlayer struct {
	Name string `json:"name"`
	ID   ID     `json:"id"`
}

func ToString(value any) string {
	switch v := value.(type) {
	case float64:
		return strconv.FormatFloat(v, 'f', 0, 64)
	case int64:
		return strconv.FormatInt(v, 10)
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	default:
		return ""
	}
}

func ToID(value string) ID {
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}
	return value
}

func UnmarshalJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
