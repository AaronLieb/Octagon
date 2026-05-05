package obs

import (
	"testing"

	"github.com/AaronLieb/octagon/config"
)

func TestGetScreenshot(t *testing.T) {
	config.Load()

	img, err := GetScreenshot()
	if err != nil {
		t.Fatalf("Failed to get screenshot: %v", err)
	}

	if img == nil {
		t.Fatal("Screenshot image is nil")
	}

	bounds := img.Bounds()
	t.Logf("Successfully captured screenshot: %dx%d", bounds.Dx(), bounds.Dy())

	SaveImageToDisk(img)
	t.Log("Screenshot saved to vod_screenshot.png")
}
