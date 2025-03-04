package commands

import (
	"github.com/AaronLieb/octagon/commands/attendees"
	"github.com/AaronLieb/octagon/commands/report"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:    "octagon",
		Usage:   "A CLI tool for running The Octagon, an SSBU tournament",
		Version: "v0.01",
		// EnableShellCompletion:  true,
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Enabled DEBUG log level",
			},
		},
		Commands: []*cli.Command{
			attendees.Command(),
			report.Command(),
		},
	}
}
