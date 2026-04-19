package store

import (
	"fmt"
	"time"
)

type SessionRow struct {
	ID          int64
	SystemID    string
	RemoteAddr  string
	BindType    string
	State       string
	ConnectedAt time.Time
	UnboundAt   *time.Time
}

func (s *Store) InsertSession(systemID, remoteAddr, bindType string) (int64, error) {
	res, err := s.db.Exec(
		`INSERT INTO sessions (system_id, remote_addr, bind_type, state) VALUES (?, ?, ?, 'bound')`,
		systemID, remoteAddr, bindType,
	)
	if err != nil {
		return 0, fmt.Errorf("insert session: %w", err)
	}
	return res.LastInsertId()
}

func (s *Store) CloseSession(id int64) error {
	_, err := s.db.Exec(
		`UPDATE sessions SET state = 'unbound', unbound_at = CURRENT_TIMESTAMP WHERE id = ?`,
		id,
	)
	return err
}

func (s *Store) GetActiveSessions() ([]SessionRow, error) {
	rows, err := s.db.Query(
		`SELECT id, system_id, remote_addr, bind_type, state, connected_at
         FROM sessions WHERE unbound_at IS NULL ORDER BY connected_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []SessionRow
	for rows.Next() {
		var r SessionRow
		if err := rows.Scan(&r.ID, &r.SystemID, &r.RemoteAddr, &r.BindType, &r.State, &r.ConnectedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, r)
	}
	return sessions, rows.Err()
}
