package dlx

type problem struct {
	// input
	rows            uint32
	columns         uint32
	matrix          [][]uint32
	partialSolution []uint32

	// internal
	nodes     [][]*node
	sentinel  *node
	candidate *candidate
}

// DLX : dancing link solver interface to solve the defined problem
type DLX interface {
	Solve(callback SolutionCallback)
	SolveRandomized(callback SolutionCallback, randomFunction func() uint32)
}

func newDLX(matrix [][]uint32) *problem {
	rows := uint32(len(matrix))
	columns := uint32(0)
	if rows > 0 {
		columns = uint32(len(matrix[0]))
	}
	p := problem{
		rows:    rows,
		columns: columns,
		matrix:  matrix,
	}
	return &p
}

// NewDLX : create a new dancing link solver for the constraint matrix provided
func NewDLX(matrix [][]uint32) DLX {
	return newDLX(matrix)
}

// NewDLXWithPartialSolution : create a new dancing link solver for the constraint matrix provided
//							   and it starts from the given partial solution
func NewDLXWithPartialSolution(matrix [][]uint32, partialSolution []uint32) DLX {
	p := newDLX(matrix)
	p.partialSolution = partialSolution
	return p
}

func (d *problem) Solve(callback SolutionCallback) {
	p := leftmost{}
	d.solveInternal(&p, callback)
}

func (d *problem) SolveRandomized(callback SolutionCallback, randomFunction func() uint32) {
	p := randomized{
		columns: d.columns,
		rf:      randomFunction,
	}
	d.solveInternal(&p, callback)
}

func (d *problem) solveInternal(heuristic picker, callback SolutionCallback) {
	d.createNodeMatrix()
	d.linkNodes()

	d.createHeaders()
	d.linkHeaders()

	d.candidate = &candidate{}

	d.fillInPartialSolution()
	d.solveRec(heuristic, callback)
}

func (d *problem) createNodeMatrix() {
	d.nodes = make([][]*node, d.rows)
	for i := uint32(0); i < d.rows; i++ {
		d.nodes[i] = make([]*node, d.columns)
		for j := uint32(0); j < d.columns; j++ {
			if d.matrix[i][j] != 0 {
				d.nodes[i][j] = &node{
					Row:    int64(i),
					Column: int64(j),
				}
			}
		}
	}
}

func (d *problem) linkNodes() {
	var first, current, previous *node
	for i := uint32(0); i < d.rows; i++ {
		current = nil
		first = nil
		previous = nil
		for j := uint32(0); j < d.columns; j++ {
			if d.nodes[i][j] != nil {
				current = d.nodes[i][j]
				if first == nil {
					first = current
				}
				if previous != nil {
					previous.Right = current
					current.Left = previous
				}
				previous = current
			}
		}
		if current == nil {
			panic("LinkNodes: LR Current is nil!")
		}
		current.Right = first
		first.Left = current
	}

	for j := uint32(0); j < d.columns; j++ {
		current = nil
		first = nil
		previous = nil
		for i := uint32(0); i < d.rows; i++ {
			if d.nodes[i][j] != nil {
				current = d.nodes[i][j]
				if first == nil {
					first = current
				}
				if previous != nil {
					previous.Down = current
					current.Up = previous
				}
				previous = current
			}
		}

		if current == nil {
			panic("LinkNodes: UD Current is nil!")
		}
		current.Down = first
		first.Up = current
	}
}

func (d *problem) createHeaders() {
	var first, last, current *node
	for i := int64(d.columns) - 1; i >= 0; i-- {
		current = &node{
			Row:    -1,
			Column: i,
		}
		if last == nil {
			last = current
		}
		if first == nil {
			first = current
		}
		current.Right = first
		current.Left = last
		first.Left = current
		last.Right = current

		first = current
	}

	d.sentinel = &node{
		Row:    -1,
		Column: -1,
	}
	d.sentinel.Right = first
	d.sentinel.Left = last

	first.Left = d.sentinel
	last.Right = d.sentinel
}

func (d *problem) linkHeaders() {
	var first *node
	header := d.sentinel.Right
	var count uint32

	for i := uint32(0); i < d.columns; i++ {
		count = 0
		for j := int64(d.rows) - 1; j >= 0; j-- {
			if d.nodes[j][i] != nil {
				count++
				first = d.nodes[j][i]
			}
		}

		header.Up = first.Up
		first.Up.Down = header
		first.Up = header
		header.Down = first
		header.Count = count

		header = header.Right
	}
}

func coverUpDown(node *node) {
	node.Up.Down = node.Down
	node.Down.Up = node.Up
}

func uncoverUpDown(node *node) {
	node.Up.Down = node
	node.Down.Up = node
}

func coverLeftRight(node *node) {
	node.Right.Left = node.Left
	node.Left.Right = node.Right
}

func uncoverLeftRight(node *node) {
	node.Right.Left = node
	node.Left.Right = node
}

func cover(header *node) {
	coverLeftRight(header)

	for node := header.Down; node != header; node = node.Down {
		if node.isHeader() {
			node.Count = node.Count - 1
		}

		for r := node.Right; r != node; r = r.Right {
			coverUpDown(r)
		}
	}
}

func uncover(header *node) {
	for node := header.Up; node != header; node = node.Up {
		if node.isHeader() {
			node.Count = node.Count + 1
		}
		for l := node.Left; l != node; l = l.Left {
			uncoverUpDown(l)
		}
	}
	uncoverLeftRight(header)

}

func headerFor(node *node) *node {
	n := node
	for !n.isHeader() {
		n = n.Up
	}
	return n
}

func (d *problem) solveRec(heuristic picker, callback SolutionCallback) bool {
	if d.sentinel.Right == d.sentinel {
		d.candidate.notify(callback)
		return true
	}

	column := heuristic.pick(d.sentinel)
	cover(column)

	for node := column.Down; node != column; node = node.Down {
		d.candidate.addNode(node)
		for r := node.Right; r != node; r = r.Right {
			cover(headerFor(r))
		}

		if d.solveRec(heuristic, callback) && !callback.More() {
			return true
		}

		d.candidate.removeNode(node)
		for l := node.Left; l != node; l = l.Left {
			uncover(headerFor(l))
		}

	}

	uncover(column)
	return false
}

func (d *problem) fillInPartialSolution() {
	var line uint32
	var node *node
	for i := 0; i < len(d.partialSolution); i++ {
		line = d.partialSolution[i]
		for j := uint32(0); j < d.columns; j++ {
			if d.nodes[line][j] != nil {
				node = d.nodes[line][j]
				d.candidate.addNode(node)
				cover(headerFor(node))
				for r := node.Right; r != node; r = r.Right {
					cover(headerFor(r))
				}
				break
			}
		}
	}
}
