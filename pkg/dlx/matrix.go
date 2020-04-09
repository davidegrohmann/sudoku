package dlx

import "fmt"

type node struct {
	Row int64
	Column int64
	Count uint32
	Left *node
	Right *node
	Up *node
	Down *node
}

func (n node) isHeader() bool {
	return n.Row == -1
}

func (n node) isSentinel() bool {
	return n.Row == -1 && n.Column == -1
}

type candidate struct {
	stack []*node
}

func (s *candidate) addNode(n *node) {
	s.stack = append(s.stack, n)
}

func (s *candidate) removeNode(n *node) {
	last := len(s.stack) -1
	if n != s.stack[last] {
		panic(fmt.Sprintf("Nodes %v and %v do not match", n, s.stack[last]))
	}
	s.stack = s.stack[:last]
}

func (s *candidate) notify(callback SolutionCallback) {
	result := make([]uint32, len(s.stack))
	j := 0
	for i := len(s.stack)-1; i >= 0; i-- {
		result[j] = uint32(s.stack[i].Row)
		j++
	}
	callback.Found(result)
}
