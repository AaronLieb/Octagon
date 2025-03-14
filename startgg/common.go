package startgg

import "context"

// TODO: Cache response
func GetTournamentSlug(ctx context.Context, tournamentShortSlug string) (string, error) {
	tournamentResp, err := GetTournament(ctx, tournamentShortSlug)
	if err != nil {
		return "", err
	}

	tournamentSlug := tournamentResp.Tournament.Slug
	return tournamentSlug, nil
}
