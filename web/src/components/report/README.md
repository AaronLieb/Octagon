# Set Reporter

React-based interface for reporting tournament sets with optimized keyboard shortcuts.

## Features

### Set List View
- Displays all reportable sets (state 1 or 2)
- Press `Space` or `/` to enter filter mode
- Fuzzy search by player name
- Arrow keys to navigate, Enter to select

### Set Reporter View
- Better seeded player on the left
- Grid showing up to 5 games (best of 3 or best of 5)
- Color-coded wins (green) and losses (red)

### Keyboard Shortcuts
- `A` - Report left player (P1) win
- `D` - Report right player (P2) win
- `Shift+A` or `Q` - Select character for P1
- Click game cell - Select character for that game
- `Enter` - Submit set (validates Bo3/Bo5)
- `Esc` - Go back

### Character Selection
- Searchable character list with Levenshtein distance sorting
- Selected character becomes default for future games
- Retroactively applies to games without characters

## Components

- `Report.tsx` - Main container
- `SetList.tsx` - List of reportable sets with filtering
- `SetReporter.tsx` - Set reporting interface
- `ListSearch.tsx` - Reusable search component

## API Endpoints

- `GET /api/sets?tournament=octagon&redemption=false` - Fetch reportable sets
- `POST /api/sets/report` - Report set with game results
- `GET /api/characters` - Get character list

## Validation

Validates that sets are either:
- Best of 3: 2-0 or 2-1
- Best of 5: 3-0, 3-1, or 3-2

Shows error message if validation fails.
