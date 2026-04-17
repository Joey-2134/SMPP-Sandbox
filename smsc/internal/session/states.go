package session

type State string

const (
	OPEN      State = "open"
	CLOSED    State = "closed"
	BOUND_TX  State = "bound_tx"
	BOUND_RX  State = "bound_rx"
	BOUND_TRX State = "bound_trx"
	UNBOUND   State = "unbound"
)
