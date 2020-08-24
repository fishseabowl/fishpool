package kademlia

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Iterator is used for peer discovery, and for finding peers by their public key.
type Iterator struct {
	sync.Mutex

	node   *Node
	table  *Table
	logger *zap.Logger

	visited map[PublicKey]struct{}
	results chan ID
	buckets [][]ID

	maxNumResults                int
	numParallelLookups           int
	numParallelRequestsPerLookup int

	lookupTimeout time.Duration
}

// NewIterator instantiates a new iterator bounded to a node and routing table
func NewIterator(node *Node, table *Table, opts ...IteratorOption) *Iterator {
	it := &Iterator{
		node:   node,
		table:  table,
		logger: node.Logger(),

		maxNumResults:                BucketSize,
		numParallelLookups:           3,
		numParallelRequestsPerLookup: 8,

		lookupTimeout: 3 * time.Second,
	}

	for _, opt := range opts {
		opt(it)
	}

	return it
}

// Find attempts to salvage through the Kademlia overlay network through the peers of the node this iterator
// is bound to for a target public key. It blocks the current goroutine until the search is complete.
func (it *Iterator) Find(target PublicKey) []ID {
	var closest []ID

	it.init(target)
	go it.iterate()

	for id := range it.results {
		closest = append(closest, id)
	}

	closest = SortByDistance(target, closest)

	if len(closest) > it.maxNumResults {
		closest = closest[:it.maxNumResults]
	}

	return closest
}

func (it *Iterator) init(target PublicKey) {
	it.results = make(chan ID, 1)

	it.visited = map[PublicKey]struct{}{
		it.node.ID().ID: {},
		target:          {},
	}

	it.buckets = make([][]ID, it.numParallelLookups)

	for i, id := range it.table.Peers() {
		it.visited[id.ID] = struct{}{}
		it.buckets[i%it.numParallelLookups] = append(it.buckets[i%it.numParallelLookups], id)
	}
}

func (it *Iterator) iterate() {
	var wg sync.WaitGroup
	wg.Add(len(it.buckets))

	for i := 0; i < len(it.buckets); i++ {
		i := i

		go func() {
			defer wg.Done()
			it.processLookupBucket(i)
		}()
	}

	wg.Wait()

	close(it.results)
}

func (it *Iterator) lookupRequest(id ID, out chan<- []ID) {
	ctx, cancel := context.WithTimeout(context.Background(), it.lookupTimeout)
	defer cancel()

	obj, err := it.node.RequestMessage(ctx, id.Address, FindNodeRequest{Target: id.ID})
	if err != nil {
		out <- nil
		return
	}

	res, ok := obj.(FindNodeResponse)
	if !ok {
		out <- nil
		return
	}

	it.results <- id

	out <- res.Results
}

func (it *Iterator) processLookupBucket(i int) {
	queue := make(chan ID, it.numParallelRequestsPerLookup)
	results := make(chan []ID, it.numParallelRequestsPerLookup)
	pending := 0

	var wg sync.WaitGroup
	wg.Add(it.numParallelRequestsPerLookup)

	for i := 0; i < it.numParallelRequestsPerLookup; i++ {
		go func() {
			defer wg.Done()
			it.processLookupRequests(queue, results)
		}()
	}

	for {
		for len(it.buckets[i]) > 0 && len(queue) < cap(queue) {
			popped := it.buckets[i][0]
			it.buckets[i] = it.buckets[i][1:]

			queue <- popped
			pending++
		}

		if pending == 0 {
			break
		}

		ids := <-results

		it.Lock()
		for _, id := range ids {
			if _, visited := it.visited[id.ID]; !visited {
				it.visited[id.ID] = struct{}{}
				it.buckets[i] = append(it.buckets[i], id)
			}
		}
		it.Unlock()

		pending--
	}

	close(queue)
	wg.Wait()
	close(results)
}

func (it *Iterator) processLookupRequests(in <-chan ID, out chan<- []ID) {
	for id := range in {
		it.lookupRequest(id, out)
	}
}
