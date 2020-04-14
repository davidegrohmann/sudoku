package dlx

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmptyDLX(t *testing.T) {
	dlx := NewDLX(make([][]uint32, 0))
	c := testCallback{}
	dlx.Solve(&c)

	require.Equal(t, [][]uint32{{}}, c.found)
}

func Test2x2LatinSquare(t *testing.T) {
	rows := 8
	columns := 12
	matrix := make([][]uint32, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]uint32, columns)
	}

	matrix[0][0] = 1
	matrix[0][4] = 1
	matrix[0][8] = 1
	matrix[1][0] = 1
	matrix[1][6] = 1
	matrix[1][10] = 1
	matrix[2][1] = 1
	matrix[2][5] = 1
	matrix[2][8] = 1
	matrix[3][1] = 1
	matrix[3][7] = 1
	matrix[3][10] = 1
	matrix[4][2] = 1
	matrix[4][4] = 1
	matrix[4][9] = 1
	matrix[5][2] = 1
	matrix[5][6] = 1
	matrix[5][11] = 1
	matrix[6][3] = 1
	matrix[6][5] = 1
	matrix[6][9] = 1
	matrix[7][3] = 1
	matrix[7][7] = 1
	matrix[7][11] = 1

	partialSolution := []uint32{4}

	dlx := NewDLXWithPartialSolution(matrix, partialSolution)
	c := testCallback{}
	dlx.Solve(&c)

	require.Equal(t, [][]uint32{{7, 2, 1, 4}}, c.found)
}

type testCallback struct {
	found [][]uint32
}

func (t *testCallback) More() bool {
	return true
}

func (t *testCallback) Found(solution []uint32) {
	t.found = append(t.found, solution)
}
