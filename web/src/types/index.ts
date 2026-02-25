export interface Set {
  id: number;
  player1: Player;
  player2: Player;
  round: string;
  entrant1: number;
  entrant2: number;
  p1Char?: string;
  p2Char?: string;
  isRedemption: boolean;
}

export interface Player {
  name: string;
  id: number;
}

export interface GameResult {
  winner: number;
  p1Char: string;
  p2Char: string;
}

export interface Character {
  id: number;
  name: string;
  aliases?: string[];
}

export interface Game {
  winner: number;
  p1Char?: Character;
  p2Char?: Character;
}

export interface Attendee {
  id: string;
  gamerTag: string;
  firstName: string;
  lastName: string;
  playerId: string;
}

export interface SeedResult {
  name: string;
  rating: number;
  seed: number;
  id: string;
}

export interface Conflict {
  player1: string;
  player2: string;
  reason: string;
  priority: number;
  expiration: string | null;
}

export interface ConflictsResponse {
  conflicts: Conflict[];
}

export const CHARACTERS: Character[] = [
  { id: 1271, name: 'Bayonetta' },
  { id: 1272, name: 'Bowser Jr' },
  { id: 1273, name: 'Bowser' },
  { id: 1274, name: 'Captain Falcon' },
  { id: 1275, name: 'Cloud' },
  { id: 1276, name: 'Corrin' },
  { id: 1277, name: 'Daisy' },
  { id: 1278, name: 'Dark Pit' },
  { id: 1279, name: 'Diddy Kong' },
  { id: 1280, name: 'Donkey Kong' },
  { id: 1282, name: 'Dr Mario' },
  { id: 1283, name: 'Duck Hunt' },
  { id: 1285, name: 'Falco' },
  { id: 1286, name: 'Fox' },
  { id: 1287, name: 'Ganondorf' },
  { id: 1289, name: 'Greninja' },
  { id: 1290, name: 'Ice Climbers' },
  { id: 1291, name: 'Ike' },
  { id: 1292, name: 'Inkling' },
  { id: 1293, name: 'Jigglypuff' },
  { id: 1294, name: 'King Dedede' },
  { id: 1295, name: 'Kirby' },
  { id: 1296, name: 'Link' },
  { id: 1297, name: 'Little Mac' },
  { id: 1298, name: 'Lucario' },
  { id: 1299, name: 'Lucas' },
  { id: 1300, name: 'Lucina' },
  { id: 1301, name: 'Luigi' },
  { id: 1302, name: 'Mario' },
  { id: 1304, name: 'Marth' },
  { id: 1305, name: 'Mega Man' },
  { id: 1307, name: 'Meta Knight' },
  { id: 1310, name: 'Mewtwo' },
  { id: 1311, name: 'Mii Brawler' },
  { id: 1313, name: 'Ness' },
  { id: 1314, name: 'Olimar' },
  { id: 1315, name: 'Pac-Man' },
  { id: 1316, name: 'Palutena' },
  { id: 1317, name: 'Peach' },
  { id: 1318, name: 'Pichu' },
  { id: 1319, name: 'Pikachu' },
  { id: 1320, name: 'Pit' },
  { id: 1321, name: 'Pokemon Trainer' },
  { id: 1322, name: 'Ridley' },
  { id: 1323, name: 'R.O.B.' },
  { id: 1324, name: 'Robin' },
  { id: 1325, name: 'Rosalina' },
  { id: 1326, name: 'Roy' },
  { id: 1327, name: 'Ryu' },
  { id: 1328, name: 'Samus' },
  { id: 1329, name: 'Sheik' },
  { id: 1330, name: 'Shulk' },
  { id: 1331, name: 'Snake' },
  { id: 1332, name: 'Sonic' },
  { id: 1333, name: 'Toon Link' },
  { id: 1334, name: 'Villager' },
  { id: 1335, name: 'Wario' },
  { id: 1336, name: 'Wii Fit Trainer' },
  { id: 1337, name: 'Wolf' },
  { id: 1338, name: 'Yoshi' },
  { id: 1339, name: 'Young Link' },
  { id: 1340, name: 'Zelda' },
  { id: 1341, name: 'Zero Suit Samus' },
  { id: 1405, name: 'Mr. Game & Watch' },
  { id: 1406, name: 'Incineroar' },
  { id: 1407, name: 'King K. Rool' },
  { id: 1408, name: 'Dark Samus' },
  { id: 1409, name: 'Chrom' },
  { id: 1410, name: 'Ken' },
  { id: 1411, name: 'Simon' },
  { id: 1412, name: 'Richter' },
  { id: 1413, name: 'Isabelle' },
  { id: 1414, name: 'Mii Swordfighter' },
  { id: 1415, name: 'Mii Gunner' },
  { id: 1441, name: 'Piranha Plant' },
  { id: 1453, name: 'Joker' },
  { id: 1526, name: 'Hero' },
  { id: 1530, name: 'Banjo & Kazooie' },
  { id: 1532, name: 'Terry' },
  { id: 1539, name: 'Byleth' },
  { id: 1747, name: 'Min Min' },
  { id: 1766, name: 'Steve' },
  { id: 1777, name: 'Sephiroth' },
  { id: 1795, name: 'Pyra/Mythra' },
  { id: 1846, name: 'Kazuya' },
  { id: 1897, name: 'Sora' },
];
