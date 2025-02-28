package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Khan/genqlient/graphql"
	"github.com/joho/godotenv"
)

const (
	API_VERSION = "alpha"
	BASE_URL    = "https://api.start.gg/gql/"
	URL         = BASE_URL + API_VERSION
)

type authedTransport struct {
	key     string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.key)
	return t.wrapped.RoundTrip(req)
}

func main() {
	var err error

	defer func() {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	err = godotenv.Load()
	if err != nil {
		return
	}

	key := os.Getenv("API_KEY")
	if key == "" {
		err = fmt.Errorf("must set API_KEY=<startgg token>")
		return
	}

	httpClient := http.Client{
		Transport: &authedTransport{
			key:     key,
			wrapped: http.DefaultTransport,
		},
	}

	ctx := context.Background()
	client := graphql.NewClient(URL, &httpClient)
	currentUserResp, err := getCurrentUser(ctx, client)
	if err != nil {
		return
	}
	currentUser := currentUserResp.CurrentUser

	fmt.Println(currentUser)

	eventResp, err := getEvent(ctx, client, "tournament/octagon-101/event/ultimate-singles")
	if err != nil {
		return
	}
	event := eventResp.Event

	fmt.Println(event)

	reg, err := generateRegistrationToken(ctx, client, event.Id, currentUser.Id)
	if err != nil {
		return
	}

	fmt.Println(reg)
}
