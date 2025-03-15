package test

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/startgg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

func Command() *cli.Command {
	return &cli.Command{
		Name:    "test",
		Usage:   "test command",
		Aliases: []string{"r"},
		Action:  test,
	}
}

func test(ctx context.Context, cmd *cli.Command) error {
	tournamentName := "octagon"
	tournamentSlug, err := startgg.GetTournamentSlug(ctx, tournamentName)
	if err != nil {
		log.Fatalf("unable to find tournament '%s': %v", tournamentName, err)
	}

	eventSlug := fmt.Sprintf("%s/event/ultimate-singles", tournamentSlug)
	setsResp, err := startgg.GetReportableSets(ctx, eventSlug)
	if err != nil {
		log.Fatalf("unable to find sets: %v", err)
	}
	sets := setsResp.Event.Sets.Nodes

	p := tea.NewProgram(initialModel(sets))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func initialModel(sets []startgg.GetReportableSetsEventSetsSetConnectionNodesSet) model {
	var choices []string
	for _, set := range sets {
		parts1 := set.Slots[0].Entrant.Participants
		parts2 := set.Slots[1].Entrant.Participants
		if len(parts1) > 0 && len(parts2) > 0 {
			p1 := set.Slots[0].Entrant.Participants[0].Player.GamerTag
			p2 := set.Slots[1].Entrant.Participants[0].Player.GamerTag
			choices = append(choices, fmt.Sprintf("%s vs %s", p1, p2))
		}
	}
	return model{
		// Our to-do list is a grocery list
		choices: choices,

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
