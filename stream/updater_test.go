package stream

import (
	"testing"

	"github.com/AaronLieb/octagon/config"
	"github.com/AaronLieb/octagon/obs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
)

func TestUpdateFromParryGG(t *testing.T) {
	config.Load()

	err := UpdateFromParryGG()
	if err != nil {
		t.Fatalf("Failed to update from Parry.gg: %v", err)
	}

	// Verify OBS text fields were updated
	client, err := obs.Client()
	if err != nil {
		t.Fatalf("Failed to get OBS client: %v", err)
	}

	fields := []string{"Player 1 Name", "Player 2 Name", "Player 1 Score", "Player 2 Score"}
	for _, field := range fields {
		resp, err := client.Inputs.GetInputSettings(&inputs.GetInputSettingsParams{
			InputName: &field,
		})
		if err != nil {
			t.Errorf("Failed to get %s: %v", field, err)
			continue
		}

		text, ok := resp.InputSettings["text"].(string)
		if !ok || text == "" {
			t.Errorf("%s is empty or not a string", field)
		} else {
			t.Logf("%s: %s", field, text)
		}
	}
}
