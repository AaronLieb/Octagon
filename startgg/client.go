/* Package startgg provides an interface to the graphql api,
* as well as various helper methods and models */
package startgg

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/charmbracelet/log"
)

var (
	client graphql.Client
	once   sync.Once
)

const (
	APIVersion = "alpha"
	BaseURL    = "https://api.start.gg/gql/"
	URL        = BaseURL + APIVersion
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
	var err error
	once.Do(func() {
		key := os.Getenv("STARTGG_API_KEY")
		if key == "" {
			err = fmt.Errorf("must set STARTGG_API_KEY=<startgg token>")
			return
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

		client = graphql.NewClient(URL, &httpClient)
	})

	return client, err
}
