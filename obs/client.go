// Package obs handles the obs websocket api and it's functions
package obs

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/andreykaipov/goobs"
	"github.com/charmbracelet/log"
)

const Host = "localhost"

var (
	client *goobs.Client
	once   sync.Once
)

func Client() (*goobs.Client, error) {
	var err error
	once.Do(func() {
		port := os.Getenv("OBS_WS_PORT")
		pass := os.Getenv("OBS_WS_PASS")

		if len(port) == 0 {
			err = errors.New("no obs web socket port provided")
			return
		}

		if len(pass) == 0 {
			err = errors.New("no obs web socket password provided")
			return
		}

		source := fmt.Sprintf("%s:%s", Host, port)
		client, err = goobs.New(source, goobs.WithPassword(pass))
		if err != nil {
			log.Error(err)
		}
	})

	return client, err
}
