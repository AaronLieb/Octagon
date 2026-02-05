// Package characters provides character name to ID mapping and fuzzy matching for Super Smash Bros Ultimate.
package characters

import (
	"regexp"
	"strings"

	"github.com/agnivade/levenshtein"
)

type Character struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Aliases []string `json:"aliases"`
}

var Characters = []Character{
	{ID: 1271, Name: "Bayonetta", Aliases: []string{"bayo"}},
	{ID: 1272, Name: "Bowser Jr", Aliases: []string{"jr", "bj"}},
	{ID: 1273, Name: "Bowser", Aliases: []string{}},
	{ID: 1274, Name: "Captain Falcon", Aliases: []string{"falcon"}},
	{ID: 1275, Name: "Cloud", Aliases: []string{}},
	{ID: 1276, Name: "Corrin", Aliases: []string{}},
	{ID: 1277, Name: "Daisy", Aliases: []string{}},
	{ID: 1278, Name: "Dark Pit", Aliases: []string{}},
	{ID: 1279, Name: "Diddy Kong", Aliases: []string{"diddy"}},
	{ID: 1280, Name: "Donkey Kong", Aliases: []string{"dk"}},
	{ID: 1282, Name: "Dr Mario", Aliases: []string{"doc"}},
	{ID: 1283, Name: "Duck Hunt", Aliases: []string{"dh"}},
	{ID: 1285, Name: "Falco", Aliases: []string{}},
	{ID: 1286, Name: "Fox", Aliases: []string{}},
	{ID: 1287, Name: "Ganondorf", Aliases: []string{"ganon"}},
	{ID: 1289, Name: "Greninja", Aliases: []string{"gren"}},
	{ID: 1290, Name: "Ice Climbers", Aliases: []string{"icies"}},
	{ID: 1291, Name: "Ike", Aliases: []string{}},
	{ID: 1292, Name: "Inkling", Aliases: []string{}},
	{ID: 1293, Name: "Jigglypuff", Aliases: []string{"jiggs", "puff"}},
	{ID: 1294, Name: "King Dedede", Aliases: []string{"ddd"}},
	{ID: 1295, Name: "Kirby", Aliases: []string{}},
	{ID: 1296, Name: "Link", Aliases: []string{}},
	{ID: 1297, Name: "Little Mac", Aliases: []string{"mac"}},
	{ID: 1298, Name: "Lucario", Aliases: []string{}},
	{ID: 1299, Name: "Lucas", Aliases: []string{}},
	{ID: 1300, Name: "Lucina", Aliases: []string{}},
	{ID: 1301, Name: "Luigi", Aliases: []string{}},
	{ID: 1302, Name: "Mario", Aliases: []string{}},
	{ID: 1304, Name: "Marth", Aliases: []string{}},
	{ID: 1305, Name: "Mega Man", Aliases: []string{}},
	{ID: 1307, Name: "Meta Knight", Aliases: []string{"mk"}},
	{ID: 1310, Name: "Mewtwo", Aliases: []string{}},
	{ID: 1311, Name: "Mii Brawler", Aliases: []string{"brawler"}},
	{ID: 1313, Name: "Ness", Aliases: []string{}},
	{ID: 1314, Name: "Olimar", Aliases: []string{}},
	{ID: 1315, Name: "Pac-Man", Aliases: []string{"pac", "pacman"}},
	{ID: 1316, Name: "Palutena", Aliases: []string{"palu"}},
	{ID: 1317, Name: "Peach", Aliases: []string{}},
	{ID: 1318, Name: "Pichu", Aliases: []string{}},
	{ID: 1319, Name: "Pikachu", Aliases: []string{"pika"}},
	{ID: 1320, Name: "Pit", Aliases: []string{}},
	{ID: 1321, Name: "Pokemon Trainer", Aliases: []string{"pt"}},
	{ID: 1322, Name: "Ridley", Aliases: []string{}},
	{ID: 1323, Name: "R.O.B.", Aliases: []string{"rob"}},
	{ID: 1324, Name: "Robin", Aliases: []string{}},
	{ID: 1325, Name: "Rosalina", Aliases: []string{"rosa"}},
	{ID: 1326, Name: "Roy", Aliases: []string{}},
	{ID: 1327, Name: "Ryu", Aliases: []string{}},
	{ID: 1328, Name: "Samus", Aliases: []string{}},
	{ID: 1329, Name: "Sheik", Aliases: []string{}},
	{ID: 1330, Name: "Shulk", Aliases: []string{}},
	{ID: 1331, Name: "Snake", Aliases: []string{}},
	{ID: 1332, Name: "Sonic", Aliases: []string{}},
	{ID: 1333, Name: "Toon Link", Aliases: []string{"tink"}},
	{ID: 1334, Name: "Villager", Aliases: []string{}},
	{ID: 1335, Name: "Wario", Aliases: []string{}},
	{ID: 1336, Name: "Wii Fit Trainer", Aliases: []string{"wiifit", "wft"}},
	{ID: 1337, Name: "Wolf", Aliases: []string{}},
	{ID: 1338, Name: "Yoshi", Aliases: []string{}},
	{ID: 1339, Name: "Young Link", Aliases: []string{"yink"}},
	{ID: 1340, Name: "Zelda", Aliases: []string{}},
	{ID: 1341, Name: "Zero Suit Samus", Aliases: []string{"zss", "zero suit"}},
	{ID: 1405, Name: "Mr. Game & Watch", Aliases: []string{"gnw", "game&watch"}},
	{ID: 1406, Name: "Incineroar", Aliases: []string{"incin"}},
	{ID: 1407, Name: "King K. Rool", Aliases: []string{"krool", "kkr"}},
	{ID: 1408, Name: "Dark Samus", Aliases: []string{}},
	{ID: 1409, Name: "Chrom", Aliases: []string{}},
	{ID: 1410, Name: "Ken", Aliases: []string{}},
	{ID: 1411, Name: "Simon", Aliases: []string{}},
	{ID: 1412, Name: "Richter", Aliases: []string{}},
	{ID: 1413, Name: "Isabelle", Aliases: []string{}},
	{ID: 1414, Name: "Mii Swordfighter", Aliases: []string{"swordfighter"}},
	{ID: 1415, Name: "Mii Gunner", Aliases: []string{"gunner"}},
	{ID: 1441, Name: "Piranha Plant", Aliases: []string{"plant", "pp"}},
	{ID: 1453, Name: "Joker", Aliases: []string{}},
	{ID: 1526, Name: "Hero", Aliases: []string{}},
	{ID: 1530, Name: "Banjo & Kazooie", Aliases: []string{"banjo", "banjo kazooie"}},
	{ID: 1532, Name: "Terry", Aliases: []string{}},
	{ID: 1539, Name: "Byleth", Aliases: []string{}},
	{ID: 1747, Name: "Min Min", Aliases: []string{}},
	{ID: 1766, Name: "Steve", Aliases: []string{}},
	{ID: 1777, Name: "Sephiroth", Aliases: []string{"seph"}},
	{ID: 1795, Name: "Pyra/Mythra", Aliases: []string{"pyra", "mythra", "pythra", "aegis", "pyra mythra"}},
	{ID: 1846, Name: "Kazuya", Aliases: []string{"kaz"}},
	{ID: 1897, Name: "Sora", Aliases: []string{}},
}

var characterMap map[string]int

func init() {
	characterMap = make(map[string]int)
	for _, char := range Characters {
		characterMap[normalize(char.Name)] = char.ID
		for _, alias := range char.Aliases {
			characterMap[normalize(alias)] = char.ID
		}
	}
}

var punctRegex = regexp.MustCompile(`[^\w\s]`)

func normalize(s string) string {
	s = strings.ToLower(s)
	s = punctRegex.ReplaceAllString(s, "")
	return strings.TrimSpace(s)
}

func GetCharacterID(name string) (int, bool) {
	normalized := normalize(name)
	if normalized == "" {
		return 0, false
	}

	bestMatch := ""
	bestDistance := len(normalized) + 1

	for charName := range characterMap {
		distance := levenshtein.ComputeDistance(normalized, normalize(charName))
		if distance < bestDistance {
			bestDistance = distance
			bestMatch = charName
		}
	}

	if bestMatch != "" {
		return characterMap[bestMatch], true
	}
	return 0, false
}
