package store

import (
	"fmt"
	"time"
)

type MessageRow struct {
	ID           int64
	MessageID    string
	SessionID    int64
	SystemID     string
	SourceAddr   string
	DestAddr     string
	ShortMessage string
	DRRequested  bool
	SubmittedAt  time.Time
	DeliveredAt  *time.Time
}

func (s *Store) InsertMessage(sessionID int64, messageID, systemID, src, dst, body string, drRequested bool) error {
	dr := 0
	if drRequested {
		dr = 1
	}
	_, err := s.db.Exec(
		`INSERT INTO messages (message_id, session_id, system_id, source_addr, dest_addr, short_message, dr_requested)
         VALUES (?, ?, ?, ?, ?, ?, ?)`,
		messageID, sessionID, systemID, src, dst, body, dr,
	)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}
	return nil
}

func (s *Store) MarkDelivered(messageID string) error {
	_, err := s.db.Exec(
		`UPDATE messages SET delivered_at = CURRENT_TIMESTAMP WHERE message_id = ?`,
		messageID,
	)
	return err
}

func (s *Store) GetRecentMessages(limit int) ([]MessageRow, error) {
	rows, err := s.db.Query(
		`SELECT id, message_id, session_id, system_id, source_addr, dest_addr, short_message,
                dr_requested, submitted_at, delivered_at
         FROM messages ORDER BY submitted_at DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []MessageRow
	for rows.Next() {
		var r MessageRow
		var dr int
		if err := rows.Scan(
			&r.ID, &r.MessageID, &r.SessionID, &r.SystemID,
			&r.SourceAddr, &r.DestAddr, &r.ShortMessage,
			&dr, &r.SubmittedAt, &r.DeliveredAt,
		); err != nil {
			return nil, err
		}
		r.DRRequested = dr == 1
		messages = append(messages, r)
	}
	return messages, rows.Err()
}
