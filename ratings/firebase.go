package ratings

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/log"
	"google.golang.org/api/option"
)

var database *db.Client

func Init(ctx context.Context) {
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
}

func Get(ctx context.Context, userId startgg.ID) (float64, error) {
	// temporary alias for subleaf
	// TODO: Separate aliases into a config file
	if userId == 4522139 {
		userId = 2966086
	} else if userId == 4620434 {
		userId = 31584
	}

	log.Debug("Fetching rating", "userId", userId)

	if database == nil {
		Init(ctx)
	}

	cachedRating, found := checkCache(userId)
	if found {
		return cachedRating, nil
	}

	path := fmt.Sprintf("/players/%s/rating", startgg.ToString(userId))
	ratingRef := database.NewRef(path)
	var rating float64
	err := ratingRef.Get(ctx, &rating)
	// TODO: add some better error handling, custom message
	log.Debug("rating", "value", rating)
	if err != nil {
		return 0, err
	}

	updateCache(userId, rating)

	return rating, nil
}
