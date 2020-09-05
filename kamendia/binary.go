package kademlia

import (
	"bytes"
	"math/bits"
	"sort"
)

// XOR allocates a new byte slice with the result of XOR(a, b).
func XOR(a, b []byte) []byte {
	if len(a) != len(b) {
		return a
	}

	c := make([]byte, len(a))

	for i := 0; i < len(a); i++ {
		c[i] = a[i] ^ b[i]
	}

	return c
}

// PrefixLen returns the number of prefixed zero bits in a.
func PrefixLen(a []byte) int {
	for i, b := range a {
		if b != 0 {
			return i*8 + bits.LeadingZeros8(b)
		}
	}

	return len(a) * 8
}

// SortByDistance sorts ids by descending XOR distance with respect to id.
func SortByDistance(id ID, ids []ID) []ID {
	sort.Slice(ids, func(i, j int) bool {
		return bytes.Compare(XOR(ids[i].ID[:], id[:]), XOR(ids[j].ID[:], id[:])) == -1
	})

	return ids
}
