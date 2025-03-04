package cache

import (
	"github.com/urfave/cli/v3"
)

// TODO: Create disk cache for tournament, events, attendees, and seeding data
func Command() *cli.Command {
	return &cli.Command{
		Name:     "cache",
		Usage:    "cache commands ",
		Aliases:  []string{"c"},
		Commands: []*cli.Command{},
	}
}
