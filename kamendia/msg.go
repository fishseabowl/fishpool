package kademlia

import (
	"fmt"
	"io"
)

// Ping is an empty ping message.
type Ping struct{}

// Marshal implements Serializable interface and returns a nil byte slice.
func (r *Ping) Marshal() []byte {
	return nil
}

// Unmarshal implements decode interface and returns a Ping message and never throws an error.
func (r *Ping) Unmarshal([]byte) error {
	r = &Ping{}
	return nil
}

// Pong is an empty pong message.
type Pong struct{}

// Marshal implements Serializable interface and returns a nil byte slice.
func (r *Pong) Marshal() []byte {
	return nil
}

// Unmarshal implements decode interface and returns a Pong instance and never throws an error.
func (r *Pong) Unmarshal([]byte) error {
	r = &Pong{}
	return nil
}

// FindNodeRequest represents a FIND_NODE RPC call in the Kademlia specification.
type FindNodeRequest struct {
	Target PublicKey
}

// Marshal implements Serializable interface and returns the public key of the target for this search
// request as a byte slice.
func (r *FindNodeRequest) Marshal() []byte {
	return r.Target[:]
}

// Unmarshal implements decode interface.
func (r *FindNodeRequest) Unmarshal(buf []byte) error {
	if len(buf) != SizePublicKey {
		r = &FindNodeRequest{}
		return fmt.Errorf("expected buf to be %d bytes, but got %d bytes: %w", SizePublicKey, len(buf), io.ErrUnexpectedEOF)
	}

	copy(r.Target[:], buf)

	return nil
}

// FindNodeResponse represents the results of a FIND_NODE RPC call
type FindNodeResponse struct {
	Results ID
}

// Marshal implements Serializable interface and encodes the list of closest peer ID results into list.
func (r *FindNodeResponse) Marshal() []byte {
	buf := []byte{byte(len(r.Results))}

	for _, result := range r.Results {
		buf = append(buf, result.Marshal()...)
	}

	return buf
}

// Unmarshal implements decode interface.
func (r *FindNodeResponse) Unmarshal(buf []byte) error {
	r = &FindNodeResponse{}

	if len(buf) < 1 {

		return io.ErrUnexpectedEOF
	}

	size := buf[0]
	buf = buf[1:]

	results := make([]ID, 0, size)

	for i := 0; i < cap(results); i++ {
		var id ID
		err := id.Unmarshal(buf).(ID)
		if err != nil {
			return io.ErrUnexpectedEOF
		}

		results = append(results, id)
		buf = buf[id.Size():]
	}

	r.Results = results

	return nil
}
