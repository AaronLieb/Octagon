package stream

import (
	"fmt"

	"github.com/AaronLieb/octagon/cache"
	"github.com/AaronLieb/octagon/parrygg/pb"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type matchItem struct {
	match       *pb.Match
	player1Name string
	player2Name string
}

func (i matchItem) Title() string {
	return fmt.Sprintf("%s vs %s", i.player1Name, i.player2Name)
}

func (i matchItem) FilterValue() string { return i.Title() }

func (i matchItem) Description() string {
	return fmt.Sprintf("%s - State: %s", i.match.Identifier, i.match.State.String())
}

type SelectModel struct {
	list    list.Model
	matches []*pb.Match
	seeds   []*pb.Seed
	success bool
}

var successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))

func NewSelectModel(matches []*pb.Match, seeds []*pb.Seed) SelectModel {
	items := make([]list.Item, len(matches))
	for i, match := range matches {
		p1, p2 := getPlayerNames(match, seeds)
		items[i] = matchItem{
			match:       match,
			player1Name: p1,
			player2Name: p2,
		}
	}

	l := list.New(items, list.NewDefaultDelegate(), 80, 20)
	l.Title = "Select Stream Match"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	return SelectModel{
		list:    l,
		matches: matches,
		seeds:   seeds,
	}
}

func getPlayerNames(match *pb.Match, seeds []*pb.Seed) (string, string) {
	p1, p2 := "Unknown", "Unknown"
	if len(match.Slots) >= 2 {
		for _, seed := range seeds {
			if seed.Id == match.Slots[0].SeedId && seed.EventEntrant != nil && len(seed.EventEntrant.Entrant.Users) > 0 {
				p1 = seed.EventEntrant.Entrant.Users[0].GamerTag
			}
			if seed.Id == match.Slots[1].SeedId && seed.EventEntrant != nil && len(seed.EventEntrant.Entrant.Users) > 0 {
				p2 = seed.EventEntrant.Entrant.Users[0].GamerTag
			}
		}
	}
	return p1, p2
}

func (m SelectModel) Init() tea.Cmd { return nil }

func (m SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.success {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(matchItem); ok {
				if err := cache.SetStreamMatch(item.match.Id); err != nil {
					log.Error("Failed to set stream match", "error", err)
					return m, tea.Quit
				}
				log.Info("Stream match set", "match", item.match.Identifier, "id", item.match.Id)
				m.success = true
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SelectModel) View() string {
	if m.success {
		return successStyle.Render("✓ Stream match set successfully\n")
	}
	return m.list.View() + "\n\nPress Enter to select, Q to quit"
}
