package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func NewModel(s *store.Store) Model {
	return Model{store: s}
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
	if m.width == 0 {
		return "Loading..."
	}

	header := renderHeader(m.width)

	// middle row: sessions on left (60%), stats on right (40%)
	// panelStyle has padding(0,1) so content width = outer - 2 (border) - 2 (padding) = outer - 4
	sessOuter := m.width * 6 / 10
	statsOuter := m.width - sessOuter
	middleContentH := 8

	sessPanel := renderSessions(m.sessions, sessOuter-4, middleContentH)
	statsPanel := renderStats(m.stats, statsOuter-4, middleContentH)
	middle := lipgloss.JoinHorizontal(lipgloss.Top, sessPanel, statsPanel)

	// messages panel fills remaining height
	headerH := lipgloss.Height(header)
	middleH := lipgloss.Height(middle)
	msgContentH := m.height - headerH - middleH - 4 // 4 = border(2) + padding(0) + breathing room
	if msgContentH < 3 {
		msgContentH = 3
	}

	msgPanel := renderMessages(m.messages, m.width-4, msgContentH)

	return lipgloss.JoinVertical(lipgloss.Left, header, middle, msgPanel)
}

func renderHeader(width int) string {
	left := lipgloss.NewStyle().Bold(true).Foreground(colorBright).Render("SMPP SMSC") +
		"  " + greenStyle.Render("●") +
		"  " + dimStyle.Render(":2775")
	right := dimStyle.Render(time.Now().Format("2006-01-02  15:04:05"))

	gap := width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}

	bar := left + strings.Repeat(" ", gap) + right
	return lipgloss.NewStyle().
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorBorder).
		Width(width).
		Render(bar)
}
