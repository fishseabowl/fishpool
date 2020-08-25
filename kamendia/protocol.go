package kademlia

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// BucketSize is a constant value of the total number of peer ID entries a single routing table bucket hold.
const BucketSize int = 16

// PingTimeout is a constant value of ping timeout
const PingTimeout time.Duration = 3 * time.Second

// ErrBucketFull is an error of bucket is full
var ErrBucketFull = errors.New("bucket is full")

// Protocol represnets routing/discovery structure of the Kademlia protocol.
type Protocol struct {
	node   *Node
	logger *zap.Logger
	table  *Table

	events Events

	pingTimeout time.Duration
}

// NewProtocol returns a new instance of the Kademlia protocol.
func NewProtocol(opts ...ProtocolOption) *Protocol {
	p := &Protocol{
		pingTimeout: PingTimeout,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Find executes the Find RPC call to find the closest peers to some given target public key. It returns the IDs of
// the closest peers it finds.
func (p *Protocol) Find(target PublicKey, opts ...IteratorOption) []ID {
	return NewIterator(p.node, p.table, opts...).Find(target)
}

// Discover executes Find to discover new peers to your node through peers your node already knows
func (p *Protocol) Discover(opts ...IteratorOption) []ID {
	return p.Find(p.node.ID().ID, opts...)
}

// Ping sends a ping request to an address, and returns no error if a pong is received back, and returns error if
// disconnect or not get pong response
func (p *Protocol) Ping(ctx context.Context, addr string) error {
	msg, err := p.node.RequestMessage(ctx, addr, Ping{})
	if err != nil {
		return fmt.Errorf("failed to ping: %w", err)
	}

	if _, ok := msg.(Pong); !ok {
		return errors.New("did not get a pong back")
	}

	return nil
}

// Table returns the routing table.
func (p *Protocol) Table() *Table {
	return p.table
}

// Ack update the routing table
func (p *Protocol) Ack(id ID) {
	for {
		inserted, err := p.table.Update(id)
		if err == nil {
			if inserted {
				if p.events.OnPeerAdmitted != nil {
					p.events.OnPeerAdmitted(id)
				}
			} else {
				if p.events.OnPeerActivity != nil {
					p.events.OnPeerActivity(id)
				}
			}

			return
		}

		last := p.table.Last(id.ID)

		ctx, cancel := context.WithTimeout(context.Background(), p.pingTimeout)
		pong, err := p.node.RequestMessage(ctx, last.Address, Ping{})
		cancel()

		if err != nil {
			if id, deleted := p.table.Delete(last.ID); deleted {
				if p.events.OnPeerEvicted != nil {
					p.events.OnPeerEvicted(id)
				}
			}
		}

		if _, ok := pong.(Pong); !ok {
			if id, deleted := p.table.Delete(last.ID); deleted {
				if p.events.OnPeerEvicted != nil {
					p.events.OnPeerEvicted(id)
				}
			}
		}

		if p.events.OnPeerEvicted != nil {
			p.events.OnPeerEvicted(id)
		}

		return
	}
}
