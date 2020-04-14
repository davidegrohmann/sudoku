package sudoku

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSudoku(t *testing.T) {
	s := &Sudoku{}
	/*
		6 0 4 0 2 0 5 0 7
		3 0 0 0 1 0 0 0 9
		0 0 0 0 0 0 0 0 0
		0 6 1 4 0 7 8 3 0
		0 0 0 2 0 8 0 0 0
		9 0 0 1 0 5 0 0 6
		2 0 0 0 0 0 0 0 5
		4 0 0 9 0 3 0 0 8
		0 0 9 0 0 0 3 0 0
	*/

	s.Table[0][0] = 6
	s.Table[0][2] = 4
	s.Table[0][4] = 2
	s.Table[0][6] = 5
	s.Table[0][8] = 7

	s.Table[1][0] = 3
	s.Table[1][4] = 1
	s.Table[1][8] = 9

	s.Table[3][1] = 6
	s.Table[3][2] = 1
	s.Table[3][3] = 4
	s.Table[3][5] = 7
	s.Table[3][6] = 8
	s.Table[3][7] = 3

	s.Table[4][3] = 2
	s.Table[4][5] = 8

	s.Table[5][0] = 9
	s.Table[5][3] = 1
	s.Table[5][5] = 5
	s.Table[5][8] = 6

	s.Table[6][0] = 2
	s.Table[6][8] = 5

	s.Table[7][0] = 4
	s.Table[7][3] = 9
	s.Table[7][5] = 3
	s.Table[7][8] = 8

	s.Table[8][2] = 9
	s.Table[8][6] = 3

	success := s.Solve()
	require.True(t, success, "sudoku solver failed")

	/*
		684329517
		372516489
		195784263
		561497832
		743268951
		928135746
		237841695
		456973128
		819652374
	*/
	expected := [9][9]uint8{
		{6, 8, 4, 3, 2, 9, 5, 1, 7},
		{3, 7, 2, 5, 1, 6, 4, 8, 9},
		{1, 9, 5, 7, 8, 4, 2, 6, 3},
		{5, 6, 1, 4, 9, 7, 8, 3, 2},
		{7, 4, 3, 2, 6, 8, 9, 5, 1},
		{9, 2, 8, 1, 3, 5, 7, 4, 6},
		{2, 3, 7, 8, 4, 1, 6, 9, 5},
		{4, 5, 6, 9, 7, 3, 1, 2, 8},
		{8, 1, 9, 6, 5, 2, 3, 7, 4},
	}

	require.Equal(t, expected, s.Table)
}
