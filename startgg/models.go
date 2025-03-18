package startgg

import "strconv"

type ID interface{}

func ToString(value interface{}) string {
	str := ""

	switch value := value.(type) {
	case float64:
		str = strconv.FormatFloat(value, 'f', 0, 64)
	case int64:
		str = strconv.FormatInt(value, 10)
	case int:
		str = strconv.Itoa(value)
	case string:
		str = value
	}

	return str
}
