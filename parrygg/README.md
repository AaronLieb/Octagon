# Parry.gg Go Client

Go client for the [Parry.gg](https://parry.gg) gRPC API.

## Setup

1. Add your API key to `.env`:
```bash
PARRYGG_API_KEY=your-api-key-here
```

2. Get your API key from [parry.gg](https://parry.gg)

## Usage

```go
import (
    "context"
    "github.com/AaronLieb/octagon/config"
    "github.com/AaronLieb/octagon/parrygg"
    "github.com/AaronLieb/octagon/parrygg/pb"
)

func main() {
    config.Load()
    apiKey := config.GetParryGGAPIKey()
    
    client, err := parrygg.NewClient(apiKey)
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    // Use authenticated context for all requests
    ctx := client.WithAuth(context.Background())
    
    // Example: Get tournaments
    req := &pb.GetTournamentsRequest{}
    resp, err := client.TournamentService.GetTournaments(ctx, req)
    if err != nil {
        panic(err)
    }
    
    // Process response...
}
```

## Available Services

- `TournamentService` - Tournament management
- `UserService` - User management
- `EventService` - Event management
- `EntrantService` - Entrant management
- `BracketService` - Bracket management
- `PhaseService` - Phase management
- `MatchService` - Match management
- `MatchGameService` - Match game management
- `GameService` - Game management
- `NotificationService` - Notification management

## Documentation

- [Parry.gg API Documentation](https://developer.parry.gg/docs/intro)
- [API Reference](https://developer.parry.gg/protodocs/services/tournament_service.proto)
- [Protocol Buffers](https://github.com/parry-gg/protos)
