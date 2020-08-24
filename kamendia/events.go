package kademlia

// Events follows the Kademlia protocol.
type Events struct {
	// OnPeerAdmitted is called when a peer is admitted to being inserted into your nodes' routing table.
	OnPeerAdmitted func(id ID)

	// OnPeerActivity is called when your node interacts with a peer, causing the peer's entry in your nodes' routing
	// table to be jumped to the head of its respective bucket.
	OnPeerActivity func(id ID)

	// OnPeerEvicted is called when your node fails to ping/dial a peer that was previously admitted into your nodes'
	// routing table, which leads to evict the peers ID from your nodes' routing table.
	OnPeerEvicted func(id ID)
}
