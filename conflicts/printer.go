package conflicts

import (
	"github.com/AaronLieb/octagon/brackets"
	"github.com/charmbracelet/log"
)

// printConflicts displays conflict information
func printConflicts(conflicts []Conflict) {
	for _, con := range conflicts {
		log.Printf("%-15s %15s  |  p%d - %s", con.Players[0].Name, con.Players[1].Name, con.Priority, con.Reason)
	}
}

// printSeeds displays seed changes between original and final seeding
func printSeeds(before, after []brackets.Player) {
	log.Printf("%-5s %-6s %25s %6s %-7s", "Seed", "Rating", "Name", "Change", "ID")
	log.Print("---------------------------------------------------------")
	
	for i, p := range after {
		for j, q := range before {
			if p == q {
				diff := j - i
				seed := i + 1
				
				switch {
				case diff > 0:
					log.Printf("%-5d %-6.1f %25s %1s%-6d%s %d", seed, p.Rating, p.Name, "\033[32m↑", diff, "\033[0m", p.ID)
				case diff < 0:
					log.Printf("%-5d %-6.1f %25s %1s%-6d%s %d", seed, p.Rating, p.Name, "\033[31m↓", -diff, "\033[0m", p.ID)
				default:
					log.Printf("%-5d %-6.1f %25s  %-6s %d", seed, p.Rating, p.Name, "", p.ID)
				}
			}
		}
	}
}
