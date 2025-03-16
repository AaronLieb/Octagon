package report

import (
	"context"
	"fmt"

	"github.com/AaronLieb/octagon/startgg"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

const listHeight = 30

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	player1 string
	player2 string
	round   string
	text    string
}

func (i item) Title() string       { return fmt.Sprintf("%s vs %s", i.player1, i.player2) }
func (i item) FilterValue() string { return i.Title() }
func (i item) Description() string { return i.round }

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func testCommand() *cli.Command {
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

func parseRound(round int) string {
	if round < 0 {
		return fmt.Sprintf("LR%d", -round)
	} else {
		return fmt.Sprintf("WR%d", round)
	}
}

func initialModel(sets []startgg.GetReportableSetsEventSetsSetConnectionNodesSet) model {
	var items []list.Item
	for _, set := range sets {
		parts1 := set.Slots[0].Entrant.Participants
		parts2 := set.Slots[1].Entrant.Participants
		if len(parts1) > 0 && len(parts2) > 0 {
			p1 := set.Slots[0].Entrant.Participants[0].Player.GamerTag
			p2 := set.Slots[1].Entrant.Participants[0].Player.GamerTag
			items = append(items, item{
				player1: p1,
				player2: p2,
				round:   parseRound(set.Round),
			})
		}
	}
	const defaultWidth = 20

	l := list.New(items, list.NewDefaultDelegate(), defaultWidth, listHeight)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Title = "Sets"
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return model{
		list: l,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.text)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("Report set for: %s", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Cancelling...")
	}
	return "\n" + m.list.View()
}
