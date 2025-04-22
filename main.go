package main

import (
	"context"
	"os"

	"github.com/AaronLieb/octagon/cache"
	"github.com/AaronLieb/octagon/cmd"
	"github.com/AaronLieb/octagon/config"
	"github.com/charmbracelet/log"
)

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		log.SetLevel(log.DebugLevel)
	}

	log.Default().SetReportTimestamp(false)

	config.Load()

	cmd := cmd.Command()

	db := cache.Open()
	defer db.Close()

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
