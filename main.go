package main

import (
	"context"
	"os"

	"github.com/AaronLieb/octagon/cache"
	"github.com/AaronLieb/octagon/commands"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("DEBUG") == "1" {
		log.SetLevel(log.DebugLevel)
	}

	log.Default().SetReportTimestamp(false)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cmd := commands.Command()

	db := cache.Open()
	defer db.Close()

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
