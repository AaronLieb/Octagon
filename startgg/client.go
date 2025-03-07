package startgg

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/charmbracelet/log"
)

var client graphql.Client

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

func GetClient(ctx context.Context) (graphql.Client, error) {
	if client != nil {
		return client, nil
	}

	key := os.Getenv("STARTGG_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("must set STARTGG_API_KEY=<startgg token>")
	}
	log.Debug("Found API Key")

	httpClient := http.Client{
		Transport: &authedTransport{
			key: key,
			wrapped: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},
		Timeout: 10 * time.Second,
	}

	client := graphql.NewClient(URL, &httpClient)
	return client, nil
}
