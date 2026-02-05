// Package ocr handles the reading of text from images
package ocr

import (
	"sync"

	"github.com/charmbracelet/log"
	"github.com/otiai10/gosseract/v2"
)

var (
	client *gosseract.Client
	once   sync.Once
)

func GetClient() *gosseract.Client {
	once.Do(func() {
		client = gosseract.NewClient()
		// single word mode
		err := client.SetVariable("tessedit_pageseg_mode", "8")
		if err != nil {
			log.Warn("failed to set tesseract pageseg mode")
		}
	})
	return client
}
