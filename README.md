# SMPP Sandbox

A self-contained SMPP sandbox that simulates full SMS infrastructure on a single machine with no external dependencies or paid services. Two components communicate over real SMPP — actual binary PDUs over TCP, not a simulation of the protocol. The only thing mocked is final delivery to a real phone number.

## Architecture

```
Java Client (ESME)  →  Go Mock SMSC  →  (in production this would be a carrier SMSC)
```

The Go SMSC accepts SMPP connections from the Java client, processes PDUs, and simulates message delivery and receipts.

## Components

### Go SMSC (`smsc/`)
Mock SMSC server. Handles TCP connections, parses SMPP PDUs, manages session state, and simulates delivery receipts. Built with a Bubbletea TUI dashboard showing live session and message state.

### Java ESME Client (`client/`)
SMPP client that connects to the Go SMSC. Sends `submit_sm` PDUs and handles incoming `deliver_sm` delivery receipts.

Both sides implement PDU parsing from scratch — no SMPP libraries used.

## PDU Support

- `bind_transmitter`, `bind_receiver`, `bind_transceiver` + responses
- `submit_sm`, `submit_sm_resp`
- `deliver_sm`, `deliver_sm_resp`
- `enquire_link`, `enquire_link_resp`
- `unbind`, `unbind_resp`
- `generic_nack`