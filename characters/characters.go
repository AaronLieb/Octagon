// Package characters provides character name to ID mapping and fuzzy matching for Super Smash Bros Ultimate.
package characters

import (
	"regexp"
	"strings"

	"github.com/agnivade/levenshtein"
)

var characterMap = map[string]int{
	"bayonetta":         1271,
	"bayo":              1271,
	"bowser jr":         1272,
	"jr":                1272,
	"bj":                1272,
	"bowser":            1273,
	"captain falcon":    1274,
	"falcon":            1274,
	"cloud":             1275,
	"corrin":            1276,
	"daisy":             1277,
	"dark pit":          1278,
	"diddy kong":        1279,
	"diddy":             1279,
	"donkey kong":       1280,
	"dk":                1280,
	"dr mario":          1282,
	"doc":               1282,
	"duck hunt":         1283,
	"falco":             1285,
	"fox":               1286,
	"ganondorf":         1287,
	"ganon":             1287,
	"greninja":          1289,
	"ice climbers":      1290,
	"icies":             1290,
	"ike":               1291,
	"inkling":           1292,
	"jigglypuff":        1293,
	"jiggs":             1293,
	"puff":              1293,
	"king dedede":       1294,
	"ddd":               1294,
	"kirby":             1295,
	"link":              1296,
	"little mac":        1297,
	"mac":               1297,
	"lucario":           1298,
	"lucas":             1299,
	"lucina":            1300,
	"luigi":             1301,
	"mario":             1302,
	"marth":             1304,
	"mega man":          1305,
	"meta knight":       1307,
	"mk":                1307,
	"mewtwo":            1310,
	"mii brawler":       1311,
	"brawler":           1311,
	"ness":              1313,
	"olimar":            1314,
	"pacman":            1315,
	"pac":               1315,
	"palutena":          1316,
	"palu":              1316,
	"peach":             1317,
	"pichu":             1318,
	"pikachu":           1319,
	"pika":              1319,
	"pit":               1320,
	"pokemon trainer":   1321,
	"pt":                1321,
	"ridley":            1322,
	"rob":               1323,
	"robin":             1324,
	"rosalina":          1325,
	"rosa":              1325,
	"roy":               1326,
	"ryu":               1327,
	"samus":             1328,
	"sheik":             1329,
	"shulk":             1330,
	"snake":             1331,
	"sonic":             1332,
	"toon link":         1333,
	"tink":              1333,
	"villager":          1334,
	"wario":             1335,
	"wii fit trainer":   1336,
	"wiifit":            1336,
	"wft":               1336,
	"wolf":              1337,
	"yoshi":             1338,
	"young link":        1339,
	"yink":              1339,
	"zelda":             1340,
	"zero suit samus":   1341,
	"zero suit":         1341,
	"zss":               1341,
	"mr game and watch": 1405,
	"gnw":               1405,
	"game&watch":        1405,
	"incineroar":        1406,
	"incin":             1406,
	"king k rool":       1407,
	"krool":             1407,
	"kkr":               1407,
	"dark samus":        1408,
	"chrom":             1409,
	"ken":               1410,
	"simon belmont":     1411,
	"richter":           1412,
	"isabelle":          1413,
	"mii swordfighter":  1414,
	"swordfighter":      1414,
	"mii gunner":        1415,
	"gunner":            1415,
	"piranha plant":     1441,
	"plant":             1441,
	"pp":                1441,
	"joker":             1453,
	"hero":              1526,
	"banjo kazooie":     1530,
	"banjo":             1530,
	"terry":             1532,
	"byleth":            1539,
	"random":            1746,
	"min min":           1747,
	"steve":             1766,
	"sephiroth":         1777,
	"seph":              1777,
	"pyra mythra":       1795,
	"pyra":              1795,
	"pythra":            1795,
	"mythra":            1795,
	"aegis":             1795,
	"kazuya":            1846,
	"kaz":               1846,
	"sora":              1897,
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
