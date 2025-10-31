package report

import (
	"context"

	"github.com/AaronLieb/octagon/tournament"
)

func FetchReportableSets(ctx context.Context, eventSlug string) ([]tournament.Set, error) {
	return tournament.FetchReportableSets(ctx, eventSlug)
}
