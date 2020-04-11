package dlx

type picker interface {
	pick(sentinel *node) *node
}

// heuristic: pick the leftmost column with the lowest count
// in this way the branching factory of the search space is minimized
// indeed the top part of the tree will have a small branching factor
type leftmost struct{}

func (leftmost) pick(sentinel *node) (result *node) {
	for current := sentinel.Right; current != sentinel; current = current.Right {
		if result == nil || result.Count > current.Count {
			result = current
		}
	}

	return
}
