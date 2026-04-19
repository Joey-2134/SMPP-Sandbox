package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const DefaultPath = "./smsc.db"

type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite %q: %w", path, err)
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	db.Exec("PRAGMA journal_mode=WAL;")
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS sessions (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id    TEXT    NOT NULL,
    remote_addr  TEXT    NOT NULL,
    bind_type    TEXT    NOT NULL,
    state        TEXT    NOT NULL,
    connected_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    unbound_at   DATETIME
);

CREATE TABLE IF NOT EXISTS messages (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id    TEXT    NOT NULL UNIQUE,
    session_id    INTEGER NOT NULL REFERENCES sessions(id),
    system_id     TEXT    NOT NULL,
    source_addr   TEXT    NOT NULL,
    dest_addr     TEXT    NOT NULL,
    short_message TEXT    NOT NULL,
    dr_requested  INTEGER NOT NULL DEFAULT 0,
    submitted_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    delivered_at  DATETIME
);
`)
	return err
}
