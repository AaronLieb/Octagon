package startgg

import "context"

const (
	// Event names
	EventUltimateSingles   = "ultimate-singles"
	EventRedemptionBracket = "redemption-bracket"

	// URL format strings
	EventSlugFormat           = "%s/event/%s"
	TournamentEventSlugFormat = "tournament/%s/event/%s"
)

func GetTournamentSlug(ctx context.Context, tournamentShortSlug string) (string, error) {
	tournamentResp, err := GetTournament(ctx, tournamentShortSlug)
	if err != nil {
		return "", err
	}

	tournamentSlug := tournamentResp.Tournament.Slug
	return tournamentSlug, nil
}
