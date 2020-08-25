package kademlia

import (
	"time"

	"go.uber.org/zap"
)

// ProtocolOption represents a functional option to configure a Protocol.
type ProtocolOption func(p *Protocol)

// WithProtocolEvents configures an event for a Protocol.
func WithProtocolEvents(events Events) ProtocolOption {
	return func(p *Protocol) {
		p.events = events
	}
}

// WithProtocolLogger configures the logger instance for a Protocol.
func WithProtocolLogger(logger *zap.Logger) ProtocolOption {
	return func(p *Protocol) {
		p.logger = logger
	}
}

// WithProtocolPingTimeout configures timeout for a Protocol.
func WithProtocolPingTimeout(pingTimeout time.Duration) ProtocolOption {
	return func(p *Protocol) {
		p.pingTimeout = pingTimeout
	}
}
