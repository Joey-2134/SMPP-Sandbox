package tui

import (
	"fmt"
	"strings"

	"github.com/joeygalvin/smpp-sandbox/smsc/internal/store"
)

func renderStats(stats store.Stats, contentWidth, contentHeight int) string {
	lines := []string{
		titleStyle.Render("STATS"),
		"",
		statRow("Active sessions", fmt.Sprintf("%d", stats.ActiveSessions)),
		statRow("Total submitted", fmt.Sprintf("%d", stats.TotalSubmitted)),
		statRow("Total delivered", fmt.Sprintf("%d", stats.TotalDelivered)),
		statRow("All-time sessions", fmt.Sprintf("%d", stats.TotalSessions)),
	}
	return panelStyle.Width(contentWidth).Height(contentHeight).Render(strings.Join(lines, "\n"))
}

func statRow(label, value string) string {
	return fmt.Sprintf("%s  %s", dimStyle.Render(label+":"), brightStyle.Render(value))
}
