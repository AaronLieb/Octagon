package conflicts

import (
	"context"
	"fmt"
	"time"

	"github.com/AaronLieb/octagon/conflicts"
	"github.com/urfave/cli/v3"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "lists all conflicts",
		Aliases: []string{"l"},
		Action:  listConflict,
	}
}

func listConflict(ctx context.Context, cmd *cli.Command) error {
	cons := conflicts.GetConflicts([]string{})
	
	if len(cons) == 0 {
		fmt.Println("No conflicts found")
		return nil
	}

	fmt.Printf("Found %d conflicts:\n\n", len(cons))
	
	for i, con := range cons {
		// Format players
		players := ""
		if len(con.Players) >= 2 {
			players = fmt.Sprintf("%-20s vs %-20s", con.Players[0].Name, con.Players[1].Name)
		}
		
		// Format expiration
		expiration := "Never"
		if con.Expiration != nil {
			if con.Expiration.Before(time.Now()) {
				expiration = "EXPIRED"
			} else {
				expiration = con.Expiration.Format("2006-01-02 15:04")
			}
		}
		
		// Print conflict
		fmt.Printf("%2d. %s | P%d | %s | %s\n", 
			i+1, players, con.Priority, expiration, con.Reason)
	}
	
	return nil
}
