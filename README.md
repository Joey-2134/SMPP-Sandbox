# SMPP Sandbox

A self-contained SMPP sandbox that simulates full SMS infrastructure on a single machine with no external dependencies or paid services. Two components communicate over real SMPP — actual binary PDUs over TCP, not a simulation of the protocol. The only thing mocked is final delivery to a real phone number.

SMSC
<p>
    <img src="/assets/SMSC_example.png" width="100%"alt="SMSC Example">
</p>
ESME
<p>
    <img src="/assets/ESME_example.png" width="100%" alt="ESME Example">
</p>


## Architecture

```
Browser  →  Java Spring Boot ESME  →  Go Mock SMSC  →  (in production this would be a carrier SMSC)
```

The Java ESME is a Spring Boot app exposing a REST API and WebSocket endpoint. A browser frontend lets you create multiple SMPP sessions, send messages interactively, and watch delivery receipts arrive in real time. The Go SMSC accepts those connections, processes PDUs, and simulates delivery.

## Components

### Go SMSC (`smsc/`)

Mock SMSC server. Handles TCP connections, parses SMPP PDUs, manages session state, and simulates delivery receipts. Built with a Bubbletea TUI dashboard showing live session and message state.

### Java ESME (`esme/`)

Spring Boot app that manages multiple concurrent SMPP sessions. Exposes a REST API for creating sessions and submitting messages, and pushes live events (delivery receipts, submit confirmations) to the browser over SSE. Each session runs a background read loop to handle server-initiated PDUs. Bind/unbind are blocking operations; submits are fire-and-forget with WebSocket callbacks.

Both sides implement PDU parsing from scratch — no SMPP libraries used.

## Local Setup

### Prerequisites


| Tool  | Version |
| ----- | ------- |
| Go    | 1.22+   |
| Java  | 21      |
| Maven | 3.x     |


### 1. Start the SMSC

```bash
cd smsc
go run ./cmd/smsc
```

The TUI dashboard will appear. The server listens on `:2775`. A SQLite database (`smsc.db`) is created in the same directory on first run.

Press `q` to quit.

### 2. Start the ESME

```bash
cd esme
mvn spring-boot:run
```

The Spring Boot app starts on `http://localhost:8080`. Open that URL in a browser to access the frontend.

### 3. Connect and send messages

From the browser frontend:

1. Create a session — this binds to the SMSC over SMPP
2. Submit a message — fill in source, destination, and body
3. Watch delivery receipts arrive in real time

Both the browser (via SSE) and the SMSC TUI will show the message flow as it happens.

---

## PDU Support

- `bind_transmitter`, `bind_receiver`, `bind_transceiver` + responses
- `submit_sm`, `submit_sm_resp`
- `deliver_sm`, `deliver_sm_resp`
- `enquire_link`, `enquire_link_resp`
- `unbind`, `unbind_resp`
- `generic_nack`