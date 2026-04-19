package tui

import (
	"fmt"
	"strings"

	"github.com/joeygalvin/smpp-sandbox/smsc/internal/store"
)

func renderMessages(messages []store.MessageRow, contentWidth, contentHeight int) string {
	lines := []string{titleStyle.Render("MESSAGE LOG"), ""}

	if len(messages) == 0 {
		lines = append(lines, dimStyle.Render("No messages yet"))
	} else {
		lines = append(lines, dimStyle.Render(fmt.Sprintf(
			"%-8s  %-12s  %-15s  %-20s  %s",
			"TIME", "SYSTEM ID", "DEST", "MESSAGE", "DR",
		)))
		// messages are newest-first from the store; reverse so newest prints at the bottom
		for i := len(messages) - 1; i >= 0; i-- {
			msg := messages[i]
			dr := dimStyle.Render("-")
			if msg.DRRequested {
				if msg.DeliveredAt != nil {
					dr = greenStyle.Render("✓")
				} else {
					dr = yellowStyle.Render("⋯")
				}
			}
			lines = append(lines, fmt.Sprintf(
				"%s  %s  %s  %s  %s",
				dimStyle.Render(msg.SubmittedAt.Format("15:04:05")),
				cyanStyle.Render(fmt.Sprintf("%-12s", truncate(msg.SystemID, 12))),
				dimStyle.Render(fmt.Sprintf("%-15s", truncate(msg.DestAddr, 15))),
				fmt.Sprintf("%-20s", truncate(msg.ShortMessage, 20)),
				dr,
			))
		}
	}

	return panelStyle.Width(contentWidth).Height(contentHeight).Render(strings.Join(lines, "\n"))
}
