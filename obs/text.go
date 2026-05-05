package obs

import (
	"github.com/andreykaipov/goobs/api/requests/inputs"
)

// UpdateText updates the text content of a text source in OBS.
// inputName is the name of the text source in OBS (e.g., "Player 1 Name").
// text is the new text content to display.
func UpdateText(inputName, text string) error {
	client, err := Client()
	if err != nil {
		return err
	}

	_, err = client.Inputs.SetInputSettings(&inputs.SetInputSettingsParams{
		InputName: &inputName,
		InputSettings: map[string]interface{}{
			"text": text,
		},
	})
	return err
}
