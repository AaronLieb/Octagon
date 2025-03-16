package commands

import (
	"github.com/AaronLieb/octagon/commands/attendees"
	"github.com/AaronLieb/octagon/commands/bracket"
	"github.com/AaronLieb/octagon/commands/cache"
	"github.com/AaronLieb/octagon/commands/conflicts"
	"github.com/AaronLieb/octagon/commands/rating"
	"github.com/AaronLieb/octagon/commands/report"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:      "octagon",
		UsageText: "octagon [subcommand]",
		Usage:     "TO automation tool for running SSBU tournaments",
		Authors: []any{
			"alieb",
		},
		Description: `A CLI tool to help automate tournament organizing for 'The Octagon',
an SSBU tournament hosted in Seattle every Tuesday at 7:00pm at the Octopus Bar.
https://start.gg/octagon`,
		Version: "v0.01",
		// EnableShellCompletion:  true,
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			attendees.Command(),
			report.Command(),
			rating.Command(),
			cache.Command(),
			conflicts.Command(),
			bracket.Command(),
		},
	}
}
