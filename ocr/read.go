package ocr

import (
	"bytes"
	"image"
	"image/png"

	"github.com/charmbracelet/log"
)

const (
	Numbers = "0123456789"
)

func Read(img image.Image, whitelist string, blacklist string) (string, error) {
	client := GetClient()

	if len(whitelist) > 0 {
		err := client.SetWhitelist(whitelist)
		if err != nil {
			log.Warn("Failed to set tesseract whitelist")
		}
	}
	if len(blacklist) > 0 {
		err := client.SetBlacklist(blacklist)
		if err != nil {
			log.Warn("Failed to set tesseract blacklist")
		}
	}

	buf := new(bytes.Buffer)

	err := png.Encode(buf, img)
	if err != nil {
		return "", err
	}

	err = client.SetImageFromBytes(buf.Bytes())
	if err != nil {
		return "", err
	}

	return client.Text()
}

func ReadText(img image.Image) (string, error) {
	return Read(img, "", "")
}

func ReadNumber(img image.Image) (string, error) {
	return Read(img, Numbers, "")
}
