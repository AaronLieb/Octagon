package commands

import (
	"github.com/AaronLieb/octagon/commands/attendeescmd"
	"github.com/AaronLieb/octagon/commands/bracketcmd"
	"github.com/AaronLieb/octagon/commands/cachecmd"
	"github.com/AaronLieb/octagon/commands/conflictscmd"
	"github.com/AaronLieb/octagon/commands/ratingcmd"
	"github.com/AaronLieb/octagon/commands/reportcmd"
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
			attendeescmd.Command(),
			reportcmd.Command(),
			ratingcmd.Command(),
			cachecmd.Command(),
			conflictscmd.Command(),
			bracketcmd.Command(),
		},
	}
}
