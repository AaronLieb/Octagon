package stream

import (
	"fmt"
	"time"

	"github.com/AaronLieb/octagon/obs"
	"github.com/AaronLieb/octagon/ocr"
	"github.com/charmbracelet/log"
)

func StartWatcher() {
	for {
		img, err := obs.GetScreenshot()
		if err != nil {
			fmt.Printf("Screenshot error: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		p1PercentImg := obs.Crop(img, P1Percent)
		p1PercentTextColor := ocr.GetAverageNonBlackColor(obs.Crop(img, P1PercentColor))
		processedImg := ocr.PreprocessPercent(p1PercentImg, p1PercentTextColor)
		redCount := ocr.GetRed(processedImg)
		log.Debug("watcher", "redCount", redCount)
		obs.SaveImageToDisk(processedImg)
		p1Percent, err := ocr.ReadNumber(processedImg)
		if err == nil && len(p1Percent) == 0 {
			log.Debug("no readable percent")
		} else if err != nil {
			fmt.Printf("OCR error: %v\n", err)
		} else {
			p1PercentProcessed, err := NormalizePercentage(p1Percent, redCount)
			if err != nil {
				log.Warn("Failed to read p1Percent", "percent", p1PercentProcessed, "err", err)
				continue
			}
			log.Infof("P1 Percent: %d", p1PercentProcessed)
		}

		time.Sleep(1 * time.Second)
	}
}
