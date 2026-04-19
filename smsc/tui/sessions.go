package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/joeygalvin/smpp-sandbox/smsc/internal/store"
)

func renderSessions(sessions []store.SessionRow, contentWidth, contentHeight int) string {
	lines := []string{titleStyle.Render("CONNECTED SESSIONS"), ""}

	if len(sessions) == 0 {
		lines = append(lines, dimStyle.Render("No active sessions"))
	} else {
		lines = append(lines, dimStyle.Render(fmt.Sprintf("%-20s %-5s %s", "SYSTEM ID", "MODE", "DURATION")))
		for _, s := range sessions {
			lines = append(lines, fmt.Sprintf(
				"%s %s %s",
				cyanStyle.Render(fmt.Sprintf("%-20s", truncate(s.SystemID, 20))),
				greenStyle.Render(fmt.Sprintf("%-5s", s.BindType)),
				dimStyle.Render(formatDuration(time.Since(s.ConnectedAt))),
			))
		}
	}

	return panelStyle.Width(contentWidth).Height(contentHeight).Render(strings.Join(lines, "\n"))
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
