package store

import "fmt"

type Stats struct {
	ActiveSessions int
	TotalSubmitted int
	TotalDelivered int
	TotalSessions  int
}

func (s *Store) GetStats() (Stats, error) {
	var st Stats
	err := s.db.QueryRow(`
        SELECT
            (SELECT COUNT(*) FROM sessions WHERE unbound_at IS NULL),
            (SELECT COUNT(*) FROM messages),
            (SELECT COUNT(*) FROM messages WHERE delivered_at IS NOT NULL),
            (SELECT COUNT(*) FROM sessions)
    `).Scan(&st.ActiveSessions, &st.TotalSubmitted, &st.TotalDelivered, &st.TotalSessions)
	if err != nil {
		return Stats{}, fmt.Errorf("get stats: %w", err)
	}
	return st, nil
}
