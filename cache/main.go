/* Package cache uses Badger to manage caching */
package cache

import (
	"time"

	"github.com/charmbracelet/log"
	badger "github.com/dgraph-io/badger/v4"
)

const PATH = "/tmp/octagon-cache"

var db *badger.DB

func Open() *badger.DB {
	log.Debug("Opening cache", "path", PATH)

	var err error
	db, err = badger.Open(badger.DefaultOptions(PATH).WithLoggingLevel(badger.WARNING))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Set(key []byte, value []byte) error {
	log.Debug("Setting key-pair in cache", "key", string(key), "value", value)
	err := db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, value).WithTTL(time.Hour * 24)
		err := txn.SetEntry(entry)
		return err
	})

	return err
}

func Get(key []byte) ([]byte, error) {
	var valCopy []byte

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			log.Debug("Fetching value from cache", "key", string(key), "value", val)

			valCopy = append([]byte{}, val...)
			return nil
		})

		return err
	})
	if err != nil {
		return nil, err
	}

	return valCopy, nil
}

func Clear() error {
	err := db.DropAll()
	return err
}
