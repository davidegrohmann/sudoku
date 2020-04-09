package sudoku

import (
	"fmt"
	"io"
	"sudoku/pkg/dlx"
)

// Size : the size of the squared sudoku board
const Size = 9

// Sudoku : a sudoku game
type Sudoku struct {
	Table [Size][Size]uint8
}

// Print : prints the current status of the sudoku board into the given writer
func (s *Sudoku) Print(writer io.Writer) (err error) {
	for i := 0; i < Size; i++ {
		for j := 0; j < Size; j++ {
			_, err = fmt.Fprintf(writer, "%d", s.Table[i][j])
			if err != nil {
				return
			}
		}
		_, err = fmt.Fprintf(writer, "\n")
		if err != nil {
			return
		}
	}
	return
}

// Check : checks if a sudoku puzzle has only one solution (with an optimal leftmost strategy)
func (s *Sudoku) Check() bool {
	callback := &checkMultiple{
		found: 0,
	}
	s.solve(callback, false)
	return callback.found == 1
}

// Solve : solve the sudoku puzzle
func (s *Sudoku) Solve() bool {
	callback := &storeFirst{
		table:       &s.Table,
		hasSolution: false,
	}
	s.solve(callback, false)
	return callback.hasSolution
}

// SolveRandomized : solve the sudoku puzzle with a randomized heuristic (useful when creating new puzzles)
func (s *Sudoku) SolveRandomized() bool {
	callback := &storeFirst{
		table:       &s.Table,
		hasSolution: false,
	}
	s.solve(callback, true)
	return callback.hasSolution
}

func (s *Sudoku) initialize() [][]uint32 {
	sudokuCoverSet := make([][]uint32, rowColumnNumberEntries)
	for i := 0; i < rowColumnNumberEntries; i++ {
		sudokuCoverSet[i] = make([]uint32, constraintEntries)
		for j := 0; j < constraintEntries; j++ {
			sudokuCoverSet[i][j] = 0
		}
	}

	var row, column uint32
	for r := uint8(0); r < Size; r++ {
		for c := uint8(0); c < Size; c++ {
			for n := uint8(1); n <= Size; n++ {
				row = toConstraintLineEntry(r, c, n)
				// row column constraint
				column = toConstraintRowColumn(r, c)
				sudokuCoverSet[row][column] = 1
				// row number constraint
				column = toConstraintRowNumber(r, n)
				sudokuCoverSet[row][column] = 1
				// column number constraint
				column = toConstraintColumnNumber(c, n)
				sudokuCoverSet[row][column] = 1
				// box number constraint
				column = toConstraintBoxNumber(r, c, n)
				sudokuCoverSet[row][column] = 1
			}
		}
	}

	return sudokuCoverSet
}

func (s *Sudoku) createPartialSolution() (result []uint32) {
	var n uint8
	for r := uint8(0); r < Size; r++ {
		for c := uint8(0); c < Size; c++ {
			n = s.Table[r][c]
			if n > 0 {
				result = append(result, toConstraintLineEntry(r, c, n))
			}
		}
	}
	return
}

func (s *Sudoku) solve(callback dlx.SolutionCallback, randomized bool) {
	sudokuCoverSet := s.initialize()
	partialSolution := s.createPartialSolution()

	DLX := dlx.NewDLXWithPartialSolution(sudokuCoverSet, partialSolution)

	if randomized {
		DLX.SolveRandomized(callback)
	} else {
		DLX.Solve(callback)
	}
}
