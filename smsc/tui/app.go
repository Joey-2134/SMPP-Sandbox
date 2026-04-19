package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joeygalvin/smpp-sandbox/smsc/internal/store"
)

type Model struct {
	store    *store.Store
	sessions []store.SessionRow
	messages []store.MessageRow
	stats    store.Stats
	width    int
	height   int
	err      error
}

func NewModel(s *store.Store) Model {
	return Model{store: s}
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		m.sessions, _ = m.store.GetActiveSessions()
		m.messages, _ = m.store.GetRecentMessages(50)
		m.stats, _ = m.store.GetStats()
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf(
		"SMPP SMSC\n\nActive sessions:  %d\nTotal submitted:  %d\nTotal delivered:  %d\nAll-time sessions: %d\n\nPress q to quit",
		m.stats.ActiveSessions,
		m.stats.TotalSubmitted,
		m.stats.TotalDelivered,
		m.stats.TotalSessions,
	)
}
