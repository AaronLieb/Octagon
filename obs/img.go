package obs

import (
	"encoding/base64"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/andreykaipov/goobs/api/requests/sources"
	"github.com/charmbracelet/log"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func Crop(img image.Image, r image.Rectangle) image.Image {
	if source, ok := img.(SubImager); ok {
		return source.SubImage(r)
	}
	return img
}

// GetScreenshot captures a screenshot of the named source and returns the decoded image
func GetScreenshot() (image.Image, error) {
	client, err := Client()
	if err != nil {
		return nil, err
	}

	source := "VOD"
	imageFormat := "png"
	width := 1920.0
	height := 1080.0
	quality := 100.0

	resp, err := client.Sources.GetSourceScreenshot(&sources.GetSourceScreenshotParams{
		ImageCompressionQuality: &quality,
		ImageWidth:              &width,
		ImageHeight:             &height,
		ImageFormat:             &imageFormat,
		SourceName:              &source,
	})
	if err != nil {
		return nil, err
	}

	// Remove data URL prefix (data:image/png;base64,)
	b64Data := strings.Split(resp.ImageData, ",")[1]

	imageData, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		log.Error("failed to decode image data")
		return nil, err
	}

	img, _, err := image.Decode(strings.NewReader(string(imageData)))
	if err != nil {
		log.Error("failed to decode image")
		return nil, err
	}

	return img, err
}

func SaveImageToDisk(img image.Image) {
	// Save image to local file for testing
	file, err := os.Create("vod_screenshot.png")
	if err != nil {
		log.Error("failed to create file")
		return
	}

	err = png.Encode(file, img)
	if err != nil {
		log.Error("failed to encode image")
		return
	}

	err = file.Close()
	if err != nil {
		log.Error("failed to close file")
		return
	}
}
