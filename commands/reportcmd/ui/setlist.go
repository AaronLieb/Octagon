package ui

import (
	"fmt"

	"github.com/AaronLieb/octagon/commands/reportcmd/startgg"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
}

func (i item) Title() string       { return fmt.Sprintf("%s vs %s", i.player1, i.player2) }
func (i item) FilterValue() string { return i.Title() }

// TODO: Update description to include event, time elapsed
func (i item) Description() string { return i.round }

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func InitialModel(sets []*startgg.Set) model {
	var items []list.Item
	for _, set := range sets {
		items = append(items, item{
			player1: set.Player1.Name,
			player2: set.Player2.Name,
			round:   set.Round,
		})
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
				m.choice = string(i.Title())
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
