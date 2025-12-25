package ratings

import (
	"context"
	"fmt"
	"os"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/AaronLieb/octagon/config"
	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"google.golang.org/api/option"
)

var (
	database *db.Client
	once     sync.Once
)

func Init(ctx context.Context) {
	once.Do(func() {
		config := &firebase.Config{}
		opt := option.WithAPIKey(os.Getenv("FIREBASE_API_KEY"))
		app, err := firebase.NewApp(ctx, config, opt)
		if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
		}
		database, err = app.DatabaseWithURL(ctx, os.Getenv("FIREBASE_DATABASE_URL"))
		if err != nil {
			log.Fatalf("error initializing db: %v\n", err)
		}
	})
}

func Get(ctx context.Context, userID startgg.ID) (float64, error) {
	// Resolve alias to real ID
	userID = config.ResolveAlias(userID)

	log.Debug("Fetching rating", "userId", userID)

	if database == nil {
		Init(ctx)
	}

	cachedRating, found := checkCache(userID)
	if found {
		return applyBias(userID, cachedRating), nil
	}

	path := fmt.Sprintf("/players/%s/rating", startgg.ToString(userID))
	ratingRef := database.NewRef(path)
	var rating float64
	err := ratingRef.Get(ctx, &rating)
	// TODO: add some better error handling, custom message
	log.Debug("rating", "value", rating)
	if err != nil {
		return 0, err
	}

	updateCache(userID, rating)

	return applyBias(userID, rating), nil
}

func applyBias(userID startgg.ID, rating float64) float64 {
	ratio := config.GetBiasForPlayer(userID)
	if ratio == 1.0 {
		return rating
	}

	biasedRating := rating * ratio
	log.Debug("Applied bias", "userId", userID, "original", rating, "ratio", ratio, "biased", biasedRating)
	return biasedRating
}
