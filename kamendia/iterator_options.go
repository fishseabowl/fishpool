package kademlia

import (
	"time"

	"go.uber.org/zap"
)

// IteratorOption is a functional option to configure Iterator.
type IteratorOption func(it *Iterator)

// WithIteratorLogger configures the logger instance for an iterator.
func WithIteratorLogger(logger *zap.Logger) IteratorOption {
	return func(it *Iterator) {
		it.logger = logger
	}
}

// WithIteratorMaxNumResults sets the max number of result peer IDs for an iterator.
func WithIteratorMaxNumResults(maxNumResults int) IteratorOption {
	return func(it *Iterator) {
		it.maxNumResults = maxNumResults
	}
}

// WithIteratorNumParallelLookups sets the max number of parallel lookup
func WithIteratorNumParallelLookups(numParallelLookups int) IteratorOption {
	return func(it *Iterator) {
		it.numParallelLookups = numParallelLookups
	}
}

// WithIteratorNumParallelRequestsPerLookup sets the max number of parallel requests peer single lookup
func WithIteratorNumParallelRequestsPerLookup(numParallelRequestsPerLookup int) IteratorOption {
	return func(it *Iterator) {
		it.numParallelRequestsPerLookup = numParallelRequestsPerLookup
	}
}

// WithIteratorLookupTimeout sets the max timeout to wait
func WithIteratorLookupTimeout(lookupTimeout time.Duration) IteratorOption {
	return func(it *Iterator) {
		it.lookupTimeout = lookupTimeout
	}
}
