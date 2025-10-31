package startgg

import "fmt"

// PrintPlayerCSV prints player information in CSV format: ID\tGamerTag
func PrintPlayerCSV(id any, gamerTag string) {
	fmt.Printf("%s\t%s\n", ToString(id), gamerTag)
}

// PrintPlayerTable prints player information in table format with contact info
func PrintPlayerTable(gamerTag string, id any, firstName, lastName string) {
	fmt.Printf("%-25s %-8s %-15s %-15s\n", gamerTag, ToString(id), firstName, lastName)
}

// PrintPlayerSimple prints player information in simple format: GamerTag ID
func PrintPlayerSimple(gamerTag string, id any) {
	fmt.Printf("%-30s %s\n", gamerTag, ToString(id))
}

// PrintSeed prints seed information: GamerTag\tSeedNum.
func PrintSeed(gamerTag string, seedNum int) {
	fmt.Printf("%s\t%d.\n", gamerTag, seedNum)
}
