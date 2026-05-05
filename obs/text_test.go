package obs

import (
	"testing"

	"github.com/AaronLieb/octagon/config"
)

func TestUpdateText(t *testing.T) {
	config.Load()

	err := UpdateText("Test", "John Doe")
	if err != nil {
		t.Fatalf("Failed to update text: %v", err)
	}
}
