package ratings

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/AaronLieb/octagon/cache"
	"github.com/charmbracelet/log"
	"github.com/dgraph-io/badger/v4"
)

const KEY_FMT = "rating-%d"

func checkCache(userId int) (float64, bool) {
	key := fmt.Sprintf(KEY_FMT, userId)

	log.Debugf("checking cache for '%s'", key)

	bytes, err := cache.Get([]byte(key))
	if err != nil {
		if err != badger.ErrKeyNotFound {
			log.Errorf("error while trying to read rating for '%d' from cache: %v", userId, err)
		}
		return 0.0, false
	}
	bits := binary.LittleEndian.Uint64(bytes)
	val := math.Float64frombits(bits)

	log.Debugf("found cached value for '%s': %f", key, val)

	return val, true
}

func updateCache(userId int, rating float64) {
	key := fmt.Sprintf(KEY_FMT, userId)

	log.Debugf("updating cache for '%s'", key)

	ratingBits := math.Float64bits(rating)
	ratingBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(ratingBytes, ratingBits)
	err := cache.Set([]byte(key), []byte(ratingBytes))
	if err != nil {
		log.Errorf("error while trying to update cache for '%d': %v", userId, err)
	}
}
