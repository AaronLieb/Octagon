package ratings

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
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

func Get(ctx context.Context, userId string) {
	if database == nil {
		Init(ctx)
	}

	path := fmt.Sprintf("/players/%s/rating", userId)
	ratingRef := database.NewRef(path)
	var rating float32
	err := ratingRef.Get(ctx, &rating)
	if err != nil {
		log.Fatalf("error reading value: %v\n", err)
	}

	fmt.Println(rating)
}
