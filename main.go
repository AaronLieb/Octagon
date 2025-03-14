package main

import (
	"context"
	"os"

	"github.com/AaronLieb/octagon/cache"
	"github.com/AaronLieb/octagon/commands"
	"github.com/AaronLieb/octagon/config"
	"github.com/charmbracelet/log"
)

func main() {
	if os.Getenv("DEBUG") == "1" {
		log.SetLevel(log.DebugLevel)
	}

	log.Default().SetReportTimestamp(false)

	config.Load()

	cmd := commands.Command()

	db := cache.Open()
	defer db.Close()

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
