package sudoku

import (
	"fmt"
)

func assertValid(n... uint8) {
	for i := range n {
		if i < 0 || i > 9 {
			panic(fmt.Sprintf("invalid value %d", i))
		}
	}
}

type storeFirst struct {
	table *[Size][Size]uint8
	hasSolution bool
}

func (s *storeFirst) More() bool {
	return false
}

func (s *storeFirst) Found(solution []uint32) {
	if len(solution) != Size*Size {
		panic(fmt.Sprintf("Solution %v has wrong size", solution))
	}

	var row, column, number uint8
	var x uint32
	for i := 0; i < len(solution); i++ {
		x = solution[i]
		row = toSudokuRow(x)
		column = toSudokuColumn(x)
		number = toSudokuNumber(x)
		assertValid(row, column, number)
		s.table[row][column] = number
	}

	s.hasSolution = true
}

type checkMultiple struct {
	found int32
}

func (s *checkMultiple) More() bool {
	return s.found < 2
}

func (s *checkMultiple) Found(solution []uint32) {
	s.found++
}