# ESME Spring Boot Implementation Plan

## Phase 1 — Spring Boot setup

- Add `spring-boot-starter-web` to `esme/pom.xml`
- Add `spring-boot-maven-plugin` so the app can be run with `mvn spring-boot:run`
- Replace `Main.java` with a `@SpringBootApplication` entry point
- Verify the app starts and serves on port 8080 with no routes yet
- All existing PDU classes stay untouched

---

## Phase 2 — SmppClient refactor

Replace the current synchronous SmppClient with one that supports concurrent sessions and server-initiated PDUs.

### Background read loop
- `connect()` starts a daemon thread that loops on `readPdu()` indefinitely
- The background thread is the only code that ever reads from the socket
- On read error the thread marks the session as disconnected and completes all pending futures exceptionally

### Blocking bind / unbind
- `ConcurrentHashMap<Integer, CompletableFuture<byte[]>> pendingResponses`
- Before writing a bind or unbind PDU, register a `CompletableFuture` keyed by sequence number
- Call `future.get(5, TimeUnit.SECONDS)` after writing — blocks until background thread delivers the response
- Background thread on receiving a response with high bit set: remove from `pendingResponses`, call `future.complete(raw)`

### Fire-and-forget submit
- `ConcurrentHashMap<Integer, Consumer<SubmitSmResp>> submitCallbacks`
- `submitSm()` registers a callback keyed by sequence number then returns immediately
- Background thread on receiving `submit_sm_resp`: remove callback, call it with the parsed response

### Server-initiated PDU handling in background thread
- `deliver_sm` — parse, call registered `DeliverSmHandler`, write `deliver_sm_resp`
- `enquire_link` — write `enquire_link_resp`
- Unknown command ID — write `generic_nack`

### Background thread dispatch logic
```
read PDU header
if commandId high bit set:
    check pendingResponses (bind_resp, unbind_resp)
    check submitCallbacks  (submit_sm_resp)
else:
    switch commandId:
        DELIVER_SM    → handle + respond
        ENQUIRE_LINK  → respond
        default       → generic_nack
```

### SmppClient fields added
- `String id` — UUID assigned at construction
- `List<SessionEvent> eventLog` — ordered log of all events on this session (bind, submits, receipts, errors)
- `Consumer<SessionEvent> eventListener` — called by background thread on each new event; wired to SSE push by Spring

---

## Phase 3 — SessionRegistry

`@Service` — application-scoped singleton.

- `ConcurrentHashMap<String, SmppClient> sessions`
- `createSession(host, port, systemId, password, bindType)` — constructs SmppClient, calls connect() + bind(), registers it, returns the session ID
- `getSession(id)` — returns SmppClient or throws if not found
- `removeSession(id)` — calls unbind(), removes from map
- `listSessions()` — returns all active sessions

---

## Phase 4 — REST API

Single `@RestController` at `/api/sessions`.

| Method | Path | Body | Response | Notes |
|--------|------|------|----------|-------|
| POST | `/api/sessions` | `{ host, port, systemId, password, bindType }` | `{ id, systemId, state }` | Blocking — returns after bind completes |
| GET | `/api/sessions` | — | array of session summaries | |
| DELETE | `/api/sessions/{id}` | — | 204 | Blocking — returns after unbind completes |
| POST | `/api/sessions/{id}/submit` | `{ from, to, message }` | 202 | Fire-and-forget — returns immediately |
| GET | `/api/sessions/{id}/events` | — | array of `SessionEvent` | Full event log for this session |

`SessionEvent` shape:
```json
{ "type": "SUBMIT_SENT | SUBMIT_ACKED | DELIVER_SM | ERROR", "timestamp": "...", "detail": "..." }
```

---

## Phase 5 — SSE (Server-Sent Events)

### How it works
- Browser opens a long-lived GET request to `/api/sessions/{id}/events/stream`
- Spring holds the connection open with `SseEmitter`
- Background thread pushes events down it as they arrive
- No config needed — built into `spring-boot-starter-web`

### Event push
- `SmppClient` holds a `Consumer<SessionEvent> eventListener` set at construction by `SessionRegistry`
- Background thread calls `eventListener.accept(event)` on:
  - `submit_sm_resp` received
  - `deliver_sm` received
  - Session error / disconnect
- `SessionRegistry` wires the listener to call `sseEmitter.send(event)` on the emitter for that session

### Emitter lifecycle
- `SseEmitter` is created when the browser opens the stream endpoint
- Stored in `SessionRegistry` alongside the `SmppClient`
- Completed (closed) when the session is unbound or the connection drops

### Browser subscription
- After creating a session, browser opens `EventSource("/api/sessions/{id}/events/stream")`
- Each event is appended to the session's message log in real time
- No JS library needed — `EventSource` is built into all browsers

---

## Phase 6 — Frontend

Single HTML page served from `esme/src/main/resources/static/index.html`. Vanilla JS, no libraries needed.

### Create session panel
- Form: host, port, systemId, password, bind type (TX / RX / TRX)
- On submit: POST `/api/sessions`, on success add a session card

### Sessions panel
- One card per active session
- Shows: sessionId, systemId, bind type, connection duration, state
- Disconnect button — DELETE `/api/sessions/{id}`

### Per-session card
- Send message form: from, to, message → POST `/api/sessions/{id}/submit`
- Live event log — `EventSource` subscription to `/api/sessions/{id}/events/stream`, each event appended as a row
- Event row shows: timestamp, type, detail (message ID, receipt status, error message)
