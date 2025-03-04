package main

import (
	"context"
	"os"

	"github.com/AaronLieb/octagon/commands"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cmd := commands.Command()

	log.Default().SetReportTimestamp(false)

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
