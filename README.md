# SMPP Sandbox

A self-contained SMPP sandbox that simulates full SMS infrastructure on a single machine with no external dependencies or paid services. Two components communicate over real SMPP — actual binary PDUs over TCP, not a simulation of the protocol. The only thing mocked is final delivery to a real phone number.

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

## PDU Support

- `bind_transmitter`, `bind_receiver`, `bind_transceiver` + responses
- `submit_sm`, `submit_sm_resp`
- `deliver_sm`, `deliver_sm_resp`
- `enquire_link`, `enquire_link_resp`
- `unbind`, `unbind_resp`
- `generic_nack`