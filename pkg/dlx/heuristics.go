package dlx

import (
	"fmt"
	"math/rand"
)

type picker interface {
	pick(sentinel *node) *node
}

// heuristic: pick the leftmost column with the lowest count
// in this way the branching factory of the search space is minimized
// indeed the top part of the tree will have a small branching factor
type leftmost struct {}

func (leftmost) pick(sentinel *node) (result *node) {
	if !sentinel.isSentinel() {
		panic(fmt.Sprintf("Node %v is not a sentinel", sentinel))
	}

	for current := sentinel.Right; current != sentinel; current = current.Right {
		if result == nil || result.Count > current.Count {
			result = current
		}
	}

	return
}

type randomized struct {
	columns uint32
}

func (r randomized) pick(sentinel *node) (result *node) {
	if !sentinel.isSentinel() {
		panic(fmt.Sprintf("Node %v is not a sentinel", sentinel))
	}

	picked := rand.Uint32() % r.columns
	result = sentinel.Right // move at least once on the right
	for i := uint32(0); i < picked; i++ {
		result = result.Right
	}

	return
}
