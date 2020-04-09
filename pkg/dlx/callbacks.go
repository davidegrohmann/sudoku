package dlx

import (
	"fmt"
	"io"
)

// SolutionCallback : to be notified when a solution has found
type SolutionCallback interface {
	More() bool
	Found(solution []uint32)
}

type printFirst struct {
	writer io.Writer
}

// PrintFirstOnly : a solution callback that prints only the first solution found
func PrintFirstOnly(writer io.Writer) SolutionCallback {
	s := printFirst{
		writer: writer,
	}
	return &s
}

func (printFirst) More() bool {
	return false
}

func (pf printFirst) Found(solution []uint32) {
	for i := range solution {
		_, _ = fmt.Fprintf(pf.writer, "%d\n", i)
	}
}

type printAll struct {
	writer io.Writer
}

// PrintAll : a solution callback that prints all solutions found
func PrintAll(writer io.Writer) SolutionCallback {
	s := printAll{
		writer: writer,
	}
	return &s
}

func (printAll) More() bool {
	return true
}

func (pf printAll) Found(solution []uint32) {
	for i := range solution {
		_, _ = fmt.Fprintf(pf.writer, "%d\n", i)
	}
	_, _ = fmt.Fprint(pf.writer, "\n")
}
