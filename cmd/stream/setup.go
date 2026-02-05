package stream

import (
	"context"
	"image"

	"github.com/AaronLieb/octagon/obs"
	"github.com/AaronLieb/octagon/ocr"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func SetupCommand() *cli.Command {
	return &cli.Command{
		Name:      "setup",
		Aliases:   []string{"s"},
		Usage:     "sets up tsh and obs",
		UsageText: "octagon stream setup",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tournament",
				Aliases: []string{"t"},
				Usage:   "Tournament slug",
				Value:   "octagon",
			},
			&cli.StringFlag{
				Name:    "event",
				Aliases: []string{"e"},
				Usage:   "Event name",
				Value:   "ultimate-singles",
			},
		},
		Action: setup,
	}
}

func setup(ctx context.Context, cmd *cli.Command) error {
	img, err := obs.GetScreenshot()
	if err != nil {
		return err
	}

	imgCropped := obs.Crop(img, image.Rect(610, 945, 640, 985))
	log.Info(ocr.GetAverageNonBlackColor(imgCropped))

	// imgProcessed := ocr.PreprocessPercent(imgCropped)
	obs.SaveImageToDisk(imgCropped)
	// fmt.Println(ocr.ReadNumber(imgCropped))

	return nil
}
