package commands

import (
	"github.com/AaronLieb/octagon/commands/attendees"
	"github.com/AaronLieb/octagon/commands/cache"
	"github.com/AaronLieb/octagon/commands/rating"
	"github.com/AaronLieb/octagon/commands/report"
	"github.com/AaronLieb/octagon/commands/seed"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:    "octagon",
		Usage:   "A CLI tool for running The Octagon, an SSBU tournament",
		Version: "v0.01",
		// EnableShellCompletion:  true,
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			attendees.Command(),
			report.Command(),
			rating.Command(),
			seed.Command(),
			cache.Command(),
		},
	}
}
