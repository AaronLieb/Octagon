package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	badger "github.com/dgraph-io/badger/v4"
)

func PopulatePlayerCache(ctx context.Context, tournament string, history int) error {
	// Extract tournament number from slug like "octagon-136"
	re := regexp.MustCompile(`(\w+)-(\d+)`)
	matches := re.FindStringSubmatch(tournament)
	if len(matches) != 3 {
		return fmt.Errorf("invalid tournament format: %s", tournament)
	}

	baseName := matches[1]
	endNum, err := strconv.Atoi(matches[2])
	if err != nil {
		return fmt.Errorf("invalid tournament number: %s", matches[2])
	}

	startNum := endNum - history + 1
	if startNum < 1 {
		startNum = 1
	}

	log.Infof("Populating player cache from %s-%d to %s-%d", baseName, startNum, baseName, endNum)

	for i := startNum; i <= endNum; i++ {
		tournamentSlug := fmt.Sprintf("%s-%d", baseName, i)
		err := cachePlayersFromTournament(ctx, tournamentSlug)
		if err != nil {
			log.Errorf("Failed to cache players from %s: %v", tournamentSlug, err)
			continue
		}
		log.Infof("Cached players from %s", tournamentSlug)
	}

	return nil
}

func cachePlayersFromTournament(ctx context.Context, tournamentSlug string) error {
	resp, err := startgg.GetParticipants(ctx, tournamentSlug)
	if err != nil {
		return err
	}

	participants := resp.Tournament.Participants.GetNodes()
	for _, participant := range participants {
		player := startgg.CachedPlayer{
			Name: participant.GamerTag,
			ID:   participant.Player.Id,
		}

		data, err := json.Marshal(player)
		if err != nil {
			log.Errorf("Failed to marshal player %s: %v", player.Name, err)
			continue
		}

		cacheKey := []byte("player_name:" + strings.ToLower(player.Name))
		err = Set(cacheKey, data)
		if err != nil {
			log.Errorf("Failed to cache player %s: %v", player.Name, err)
		}
	}

	return nil
}
func GetAllCachedPlayers() ([]startgg.CachedPlayer, error) {
	var players []startgg.CachedPlayer
	
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		
		prefix := []byte("player_name:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var player startgg.CachedPlayer
				if err := json.Unmarshal(val, &player); err == nil {
					players = append(players, player)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	
	return players, err
}
