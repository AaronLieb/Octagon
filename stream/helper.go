package stream

import (
	"strconv"

	"github.com/charmbracelet/log"
)

func NormalizePercentage(percentStr string, redCount float64) (int, error) {
	log.Debug(percentStr)
	percent, err := strconv.Atoi(percentStr)
	if err != nil {
		return 0, err
	}

	if percent > 1000 {
		digit1 := percent / 1000
		digit2 := (percent / 100) % 10
		percent = min(digit1, digit2)*100 + percent%100
	}

	if percent/100 == 7 {
		/* 1 is often confused for 7, and in this case, we know know that
		 * 7xx percent is impossible and is most likely a 1xx percent */
		percent = 100 + percent%100
	}

	percent = max(0, percent)
	return percent, nil
}
