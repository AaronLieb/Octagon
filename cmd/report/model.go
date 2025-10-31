// Package report provides interactive tournament set reporting functionality using bubbletea TUI.
package report

import (
	"context"
	"fmt"
	"strings"

	"github.com/AaronLieb/octagon/characters"
	"github.com/AaronLieb/octagon/tournament"
	"github.com/AaronLieb/octagon/validation"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateSetList state = iota
	stateReportGames
	stateCharacterInput
	stateComplete
)

type (
	Set        = tournament.Set
	Player     = tournament.Player
	GameResult = tournament.GameResult
)

type setItem struct {
	set Set
}

func (i setItem) Title() string {
	return fmt.Sprintf("%s vs %s", i.set.Player1.Name, i.set.Player2.Name)
}
func (i setItem) FilterValue() string { return i.Title() }
func (i setItem) Description() string { return i.set.Round }

type Model struct {
	ctx         context.Context
	state       state
	list        list.Model
	selectedSet Set
	charInput   textinput.Model
	gameResults []GameResult
	currentGame int
	charStep    int // 0 = P1 char, 1 = P2 char
	lastP1Char  string
	lastP2Char  string
	error       string
	success     bool
}

var (
	formStyle  = lipgloss.NewStyle().Margin(1, 2)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	quitStyle  = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	winStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
)

func NewModel(ctx context.Context, sets []tournament.Set) Model {
	items := make([]list.Item, len(sets))
	for i, set := range sets {
		items[i] = setItem{set: set}
	}

	l := list.New(items, list.NewDefaultDelegate(), 80, 20)
	l.Title = "Select Set to Report"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	charInput := textinput.New()
	charInput.Placeholder = "character name"

	return Model{
		ctx:       ctx,
		state:     stateSetList,
		list:      l,
		charInput: charInput,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateSetList:
		return m.updateSetList(msg)
	case stateReportGames:
		return m.updateReportGames(msg)
	case stateCharacterInput:
		return m.updateCharacterInput(msg)
	case stateComplete:
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updateSetList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(setItem); ok {
				m.selectedSet = item.set
				m.gameResults = make([]GameResult, 5)
				m.currentGame = 0
				m.state = stateReportGames
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) updateReportGames(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = stateSetList
			m.error = ""
			return m, nil
		case "enter":
			if err := m.validateAndReportSet(); err != nil {
				m.error = err.Error()
			} else {
				m.success = true
				m.state = stateComplete
			}
		case "j":
			if m.currentGame < len(m.gameResults) {
				m.gameResults[m.currentGame].Winner = 1
				m.gameResults[m.currentGame].P1Char = m.lastP1Char
				m.gameResults[m.currentGame].P2Char = m.lastP2Char
				m.gameResults[m.currentGame].P1CharID, _ = characters.GetCharacterID(m.lastP1Char)
				m.gameResults[m.currentGame].P2CharID, _ = characters.GetCharacterID(m.lastP2Char)
				m.currentGame++
			}
		case "k":
			if m.currentGame < len(m.gameResults) {
				m.gameResults[m.currentGame].Winner = 2
				m.gameResults[m.currentGame].P1Char = m.lastP1Char
				m.gameResults[m.currentGame].P2Char = m.lastP2Char
				m.gameResults[m.currentGame].P1CharID, _ = characters.GetCharacterID(m.lastP1Char)
				m.gameResults[m.currentGame].P2CharID, _ = characters.GetCharacterID(m.lastP2Char)
				m.currentGame++
			}
		case "J": // Shift+J - change P1 character
			if m.currentGame < len(m.gameResults) {
				m.charStep = 0
				m.state = stateCharacterInput
				m.charInput.SetValue(m.lastP1Char)
				m.charInput.Focus()
			}
		case "K": // Shift+K - change P2 character
			if m.currentGame < len(m.gameResults) {
				m.charStep = 1
				m.state = stateCharacterInput
				m.charInput.SetValue(m.lastP2Char)
				m.charInput.Focus()
			}
		case "backspace":
			if m.currentGame > 0 {
				m.currentGame--
				m.gameResults[m.currentGame] = GameResult{}
			}
		}
	}
	return m, nil
}

func (m Model) updateCharacterInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.state = stateReportGames
			return m, nil
		case "enter":
			charName := strings.TrimSpace(m.charInput.Value())

			charID, _ := characters.GetCharacterID(charName)

			// Update character data
			if m.charStep == 0 {
				m.gameResults[m.currentGame].P1Char = charName
				m.gameResults[m.currentGame].P1CharID = charID
				m.lastP1Char = charName
			} else {
				m.gameResults[m.currentGame].P2Char = charName
				m.gameResults[m.currentGame].P2CharID = charID
				m.lastP2Char = charName
			}

			m.state = stateReportGames
			m.error = ""
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.charInput, cmd = m.charInput.Update(msg)
	return m, cmd
}

func (m Model) validateAndReportSet() error {
	if err := validation.ValidateSetScore(m.gameResults); err != nil {
		return err
	}
	return tournament.ReportSet(m.ctx, m.selectedSet, m.gameResults)
}

func (m Model) View() string {
	switch m.state {
	case stateSetList:
		return "\n" + m.list.View()
	case stateReportGames:
		return m.reportGamesView()
	case stateCharacterInput:
		return m.characterInputView()
	case stateComplete:
		if m.success {
			return quitStyle.Render("Set reported successfully!")
		}
		return quitStyle.Render("Cancelled.")
	}
	return ""
}

func (m Model) reportGamesView() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Reporting: %s vs %s (%s)\n\n",
		m.selectedSet.Player1.Name, m.selectedSet.Player2.Name, m.selectedSet.Round))

	// Table header
	p1Name := m.selectedSet.Player1.Name
	p2Name := m.selectedSet.Player2.Name

	// Truncate names to 20 characters
	if len(p1Name) > 20 {
		p1Name = p1Name[:20]
	}
	if len(p2Name) > 20 {
		p2Name = p2Name[:20]
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Underline(true).Width(20)
	b.WriteString(fmt.Sprintf("%-8s %s %s\n",
		"Game",
		headerStyle.Render(p1Name),
		headerStyle.Render(p2Name)))

	// Table rows
	for i := range 5 {
		gameLabel := fmt.Sprintf("Game %d", i+1)

		var p1Cell, p2Cell string

		if i < m.currentGame {
			result := m.gameResults[i]

			// Player 1 cell
			char1 := result.P1Char
			if char1 == "" {
				char1 = "no character"
			}
			if len(char1) > 20 {
				char1 = char1[:20]
			}
			if result.Winner == 1 {
				p1Cell = winStyle.Width(20).Render(char1)
			} else {
				p1Cell = errorStyle.Width(20).Render(char1)
			}

			// Player 2 cell
			char2 := result.P2Char
			if char2 == "" {
				char2 = "no character"
			}
			if len(char2) > 20 {
				char2 = char2[:20]
			}
			if result.Winner == 2 {
				p2Cell = winStyle.Width(20).Render(char2)
			} else {
				p2Cell = errorStyle.Width(20).Render(char2)
			}
		} else if i == m.currentGame {
			p1Cell = lipgloss.NewStyle().Width(20).Render("[ j to win ]")
			p2Cell = lipgloss.NewStyle().Width(20).Render("[ k to win ]")
		} else {
			p1Cell = lipgloss.NewStyle().Width(20).Render("-")
			p2Cell = lipgloss.NewStyle().Width(20).Render("-")
		}

		b.WriteString(fmt.Sprintf("%-8s %s %s\n", gameLabel, p1Cell, p2Cell))
	}

	b.WriteString("\n")

	if m.error != "" {
		b.WriteString(errorStyle.Render("Error: " + m.error))
		b.WriteString("\n\n")
	}

	b.WriteString("j/k: report winner • Enter: finish set • Shift+J/K: change character • backspace: undo • Esc: back")

	return formStyle.Render(b.String())
}

func (m Model) characterInputView() string {
	var b strings.Builder

	playerName := m.selectedSet.Player1.Name
	if m.charStep == 1 {
		playerName = m.selectedSet.Player2.Name
	}

	b.WriteString(fmt.Sprintf("Game %d - %s character: ", m.currentGame+1, playerName))
	b.WriteString(m.charInput.View())
	b.WriteString("\n\n")

	if m.error != "" {
		b.WriteString(errorStyle.Render("Error: " + m.error))
		b.WriteString("\n\n")
	}

	b.WriteString("Enter: confirm (empty = use previous) • Esc: back")

	return formStyle.Render(b.String())
}
