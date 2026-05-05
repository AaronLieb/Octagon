# OBS Module

Go module for interacting with OBS Studio via WebSocket.

## Setup

Configure OBS WebSocket credentials in `.env`:

```bash
OBS_WS_PORT=4455
OBS_WS_PASS=your-password
```

## Usage

### Update Text Sources

```go
import "github.com/AaronLieb/octagon/obs"

// Update a text source
err := obs.UpdateText("Player 1 Name", "John Doe")
if err != nil {
    log.Fatal(err)
}
```

**Required information:**
- `inputName` - The exact name of the text source in OBS (case-sensitive)
- `text` - The new text content to display

### Get the Text Source Name

To find the correct input name:
1. Open OBS Studio
2. Look at your Sources panel
3. Find the text source you want to update
4. Use that exact name (e.g., "Player 1 Name", "Score", "Timer")

## Dependencies

Uses [goobs](https://github.com/andreykaipov/goobs) for OBS WebSocket communication.
