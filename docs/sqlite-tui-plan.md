# SQLite Persistence & TUI Implementation Plan

## Current State

The SMSC server is fully functional at the protocol level:
- PDU parsing/serialisation is complete and unit tested
- Session state machine handles bind, submit, deliver, enquire_link, unbind
- `session.Manager` holds live sessions in memory (keyed by `net.Conn`)
- `session.Session` tracks `State`, `Conn`, `SystemID`, and `SequenceNumber`
- `main.go` wires server callbacks to the manager - no TUI, no persistence

The next two milestones (SQLite + TUI) must be implemented together because the TUI needs a live data source that the store provides.

---

## Part 1: SQLite Persistence (`smsc/internal/store`)

### Goal
Persist every session lifecycle event and every message so the TUI can display live stats and the data survives a restart.

### Dependencies

Add to `go.mod`:
```
modernc.org/sqlite v1.x   # pure Go, no CGO required
```

Run: `go get modernc.org/sqlite`

---

### Schema

**`sessions` table** — one row per SMPP session, updated in place as state changes.

```sql
CREATE TABLE IF NOT EXISTS sessions (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id   TEXT    NOT NULL,
    remote_addr TEXT    NOT NULL,
    bind_type   TEXT    NOT NULL,  -- 'TX', 'RX', 'TRX'
    state       TEXT    NOT NULL,  -- mirrors session.State
    connected_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    unbound_at  DATETIME           -- NULL until session ends
);
```

**`messages` table** — one row per submit_sm received, updated when a DR is sent.

```sql
CREATE TABLE IF NOT EXISTS messages (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id      TEXT    NOT NULL UNIQUE,  -- the ID we generated and returned
    session_id      INTEGER NOT NULL REFERENCES sessions(id),
    source_addr     TEXT    NOT NULL,
    dest_addr       TEXT    NOT NULL,
    short_message   TEXT    NOT NULL,
    dr_requested    INTEGER NOT NULL DEFAULT 0,  -- boolean (0/1)
    submitted_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    delivered_at    DATETIME           -- NULL until DR sent
);
```

---

### `store.go` API

File: `smsc/internal/store/store.go`

```go
type Store struct { db *sql.DB }

func New(path string) (*Store, error)           // opens DB, runs migrations
func (s *Store) Close() error

// Session lifecycle
func (s *Store) InsertSession(systemID, remoteAddr, bindType string) (int64, error)
func (s *Store) UpdateSessionState(id int64, state string) error
func (s *Store) CloseSession(id int64) error    // sets unbound_at = NOW

// Message lifecycle
func (s *Store) InsertMessage(sessionID int64, messageID, src, dst, body string, drRequested bool) error
func (s *Store) MarkDelivered(messageID string) error

// Read queries (used by TUI)
func (s *Store) GetActiveSessions() ([]SessionRow, error)
func (s *Store) GetRecentMessages(limit int) ([]MessageRow, error)
func (s *Store) GetStats() (Stats, error)       // counts for the stats panel
```

`Stats` struct:
```go
type Stats struct {
    ActiveSessions  int
    TotalSubmitted  int
    TotalDelivered  int
    TotalSessions   int  // all time
}
```

`SessionRow` and `MessageRow` are plain structs mirroring the table columns — no ORM.

---

### Wiring the Store into the Session

`session.Session` currently has no reference to the store. Two options:

**Option A (preferred):** Pass the store into `Session` at construction time. The session calls store methods directly as events happen.

Changes:
- Add `Store *store.Store` and `SessionID int64` fields to `session.Session`
- `NewSession` accepts a `*store.Store` parameter; calls `InsertSession` and stores the returned ID
- `handleBind` calls `store.UpdateSessionState` after the state transition
- `handleSubmitSM` calls `store.InsertMessage` after generating the message ID, and `store.MarkDelivered` after sending the DR
- `handleUnbind` calls `store.CloseSession`

**Option B:** Keep `Session` pure and have the `Manager` callbacks do all store writes by observing the session after each `Handle` call. This is messier because the manager doesn't know *what* changed inside Handle.

Use **Option A**.

---

### Avoiding Import Cycles

The dependency graph must be:

```
main → session → store → (sqlite driver only)
main → server   (no store dependency)
tui  → store    (reads only)
```

`store` must not import `session` or `server`. `session` imports `store` but not `tui`.

---

## Part 2: TUI Dashboard (`smsc/tui`)

### Goal
A Bubbletea terminal dashboard that shows the live state of the SMSC. It polls the store on a ticker and re-renders.

### Dependencies

Add to `go.mod`:
```
github.com/charmbracelet/bubbletea  v1.x
github.com/charmbracelet/lipgloss   v1.x
```

Run: `go get github.com/charmbracelet/bubbletea github.com/charmbracelet/lipgloss`

---

### Layout

Three panels, stacked vertically (or side-by-side if terminal is wide enough):

```
┌─────────────────────────────────────────────────────────┐
│  SMPP SMSC  ●  :2775  ─────────────────── 14:32:01      │  ← header bar
├────────────────────────┬────────────────────────────────┤
│  CONNECTED SESSIONS    │  STATS                         │
│  system_id  mode  dur  │  Active sessions:   2          │
│  esme01     TRX   00:42│  Total submitted:  17          │
│  esme02     TX    00:08│  Total delivered:  12          │
│                        │  All-time sessions: 5          │
├────────────────────────┴────────────────────────────────┤
│  PDU LOG                                                 │
│  14:31:58  esme01  →  submit_sm    id=14  to=447700…    │
│  14:31:58  esme01  ←  submit_sm_resp  id=14  ESME_ROK   │
│  14:31:59  esme01  ←  deliver_sm   id=14  (DR)          │
│  14:31:59  esme01  →  deliver_sm_resp                   │
│  14:32:00  esme02  →  enquire_link                      │
│  14:32:00  esme02  ←  enquire_link_resp                 │
└─────────────────────────────────────────────────────────┘
```

---

### File Structure

```
smsc/tui/
├── app.go        # Root Bubbletea Model, layout, tick handling
├── sessions.go   # Sessions panel view
├── messages.go   # PDU log panel view
└── stats.go      # Stats panel view
```

---

### `app.go` — Root Model

```go
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
```

`Init()` returns `tickCmd()`.

`Update()` handles:
- `tea.WindowSizeMsg` — store width/height for layout
- `tickMsg` — call `store.GetActiveSessions()`, `store.GetRecentMessages(50)`, `store.GetStats()`, then return `tickCmd()` to schedule the next tick
- `tea.KeyMsg` for `q` / `ctrl+c` → `tea.Quit`

`View()` calls the three panel renderers and joins them with lipgloss.

---

### `sessions.go` — Sessions Panel

Takes `[]store.SessionRow`, renders a table with columns:
- `SYSTEM ID` — `session.SystemID`
- `MODE` — bind type (TX / RX / TRX)
- `SINCE` — format `connected_at` as HH:MM:SS, or show duration (elapsed)

Highlight active sessions in green, unbound ones greyed out (if shown at all — consider only showing currently connected sessions).

---

### `messages.go` — PDU Log Panel

Takes `[]store.MessageRow` (most recent N messages).

Each row shows:
- Timestamp (`submitted_at`)
- `system_id` (via JOIN with sessions, or store the system_id on the message row — simpler)
- Direction arrow and PDU type
- Key fields: message ID, dest addr (truncated), DR status

Keep a scrolling window of the last 50 rows. No interactive scroll needed for v1 — just show the most recent that fit.

> **Simplification:** Store `system_id` directly on `messages` to avoid a JOIN in the hot read path.

---

### `stats.go` — Stats Panel

Takes `store.Stats`, renders a simple key-value list with lipgloss styling. No logic — just formatting.

---

### The PDU Log Problem

The store records `submit_sm` and `deliver_sm` (message-level events), but the raw PDU log in the design above shows every PDU exchange including `enquire_link`. Two approaches:

**Option A (simpler, ship it first):** The PDU log only shows message-level events from the `messages` table. `enquire_link` is omitted. This is a clean v1.

**Option B (complete):** Add a `pdu_log` table with a row per PDU (direction, command_id, sequence_number, timestamp). The session writes to it on every send/receive. The TUI reads the last 50.

Start with **Option A**. Add `pdu_log` table if the log feels too sparse once the TUI is running.

---

## Part 3: Wiring Everything in `main.go`

Current `main.go` creates a `Manager` and a `Server`. After these changes:

```go
func main() {
    // 1. Open store
    store, err := store.New("smsc.db")
    if err != nil { log.Fatal(err) }
    defer store.Close()

    // 2. Create session manager (now takes a store)
    manager := session.NewManager(store)

    // 3. Create server (unchanged interface)
    s := server.NewServer(":2775", ...)

    // 4. Start server in a goroutine
    go s.Start()

    // 5. Start TUI (blocks until user quits)
    p := tea.NewProgram(tui.NewModel(store), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

The TUI runs on the main goroutine (Bubbletea requirement). The server runs on a background goroutine. Both share the store, which uses a single `*sql.DB` — SQLite handles concurrent writes via its internal mutex; the store methods can add their own `sync.Mutex` if write contention becomes a problem.

---

## Part 4: Session Changes Needed

`session.Session` needs two new fields:

```go
type Session struct {
    State          State
    Conn           net.Conn
    SystemID       string
    SequenceNumber uint32
    // new:
    store     *store.Store
    sessionID int64
}
```

`NewSession` signature change:
```go
func NewSession(conn net.Conn, s *store.Store) *Session
```

On construction, `NewSession` calls:
```go
id, _ := s.InsertSession("", conn.RemoteAddr().String(), "")
// SystemID and bind type not yet known — updated in handleBind
```

In `handleBind`, after setting `s.SystemID` and `s.State`:
```go
s.store.UpdateSessionState(s.sessionID, string(s.State))
// also update system_id and bind_type — add UpdateSessionBind method or fold into InsertSession
```

Cleanest approach: don't insert the session row until `handleBind` succeeds (we don't know `system_id` or `bind_type` until then). Keep a flag `bound bool` and only insert on first successful bind. If the client connects and never binds, no row is written.

Revised flow:
1. `NewSession` — no store write yet
2. `handleBind` success — `InsertSession(systemID, remoteAddr, bindType)` → store `sessionID`
3. `handleSubmitSM` — `InsertMessage(...)` and conditionally `MarkDelivered(...)`
4. `handleUnbind` — `CloseSession(sessionID)`
5. Manager's `RemoveSession` callback — `CloseSession(sessionID)` as a safety net for dropped connections

---

## Implementation Order

1. **Add dependencies** — `go get modernc.org/sqlite bubbletea lipgloss`
2. **Implement `store/store.go`** — schema, migrations, all methods, unit test with `:memory:` DB
3. **Update `session.Session`** — add store fields, wire store calls into handlers
4. **Update `session.Manager`** — pass store through, `NewManager` takes `*store.Store`
5. **Update `main.go`** — open store, pass to manager, start server in goroutine
6. **Implement `tui/app.go`** — root model, tick loop, window resize
7. **Implement `tui/stats.go`** — stats panel (simplest, build confidence)
8. **Implement `tui/sessions.go`** — sessions panel
9. **Implement `tui/messages.go`** — PDU log panel
10. **Wire TUI into `main.go`** — `tea.NewProgram`, move server to goroutine
11. **Manual end-to-end test** — connect Java ESME, send messages, verify TUI updates

---

## Open Questions

- **SQLite WAL mode?** With one writer (session goroutines) and one reader (TUI ticker), WAL mode (`PRAGMA journal_mode=WAL`) prevents the TUI read from blocking writes. Add this to `New()`.
- **Message log truncation?** The `messages` table will grow forever in a long-running dev session. For v1 this is fine. Later, add a `LIMIT` on inserts or a periodic delete of rows older than N hours.
- **Session ID in messages table?** Storing `system_id` directly on `messages` (denormalised) avoids a JOIN and simplifies the TUI read query. Do this.
