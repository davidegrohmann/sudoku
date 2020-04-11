package sudoku

// RandomFill : fille randomly 9 cells with numbers from 1 to 9
func (s *Sudoku) RandomFill(rng func() uint32) {
	var row, column uint8
	for i := uint8(1); i <= Size; i++ {
		row = uint8(rng() % Size)
		column = uint8(rng() % Size)
		// move until the cell is not taken
		for s.Table[row][column] != 0 {
			column = (column + 1) % Size
			if column == 0 {
				row = (row + 1) % Size
			}
		}
		s.Table[row][column] = i
	}
}

// RandomShuffleRows : shuffle rows in each block of 3 in order to preserve sudoku properties
func (s *Sudoku) RandomShuffleRows(rng func() uint32) {
	var row, tmp uint8
	for i := uint8(0); i < Size; i+=3 {
		for j := i; j < i+blockSize; j++ {
			row = i + uint8(rng()%blockSize)
			if j != row {
				for k := uint8(0); k < Size; k++ {
					tmp = s.Table[j][k]
					s.Table[j][k] = s.Table[row][k]
					s.Table[row][k] = tmp
				}
			}
		}
	}
}

// RandomShuffleColumns : shuffle columns in each block of 3 in order to preserve sudoku properties
func (s *Sudoku) RandomShuffleColumns(rng func() uint32) {
	var column, tmp uint8
	for i := uint8(0); i < Size; i+=3 {
		for j := i; j < i+blockSize; j++ {
			column = i + uint8(rng()%blockSize)
			if j != column {
				for k := uint8(0); k < Size; k++ {
					tmp = s.Table[k][j]
					s.Table[k][j] = s.Table[k][column]
					s.Table[k][column] = tmp
				}
			}
		}
	}
}

// RandomShuffleBlockRows : shuffle blocks of 3 rows in order to preserve sudoku properties
func (s *Sudoku) RandomShuffleBlockRows(rng func() uint32) {
	var rBlock, tmp uint8
	for iBlock := uint8(0); iBlock < blockSize; iBlock++ {
		rBlock = uint8(rng() % blockSize)
		if iBlock != rBlock {
			for k := uint8(0); k < Size; k++ {
				for j := uint8(0); j < blockSize; j++ {
					tmp = s.Table[iBlock*blockSize+j][k]
					s.Table[iBlock*blockSize+j][k] = s.Table[rBlock*blockSize+j][k]
					s.Table[rBlock*blockSize+j][k] = tmp
				}
			}
		}
	}
}

// RandomShuffleBlockColumns : shuffle blocks of 3 columns in order to preserve sudoku properties
func (s *Sudoku) RandomShuffleBlockColumns(rng func() uint32) {
	var cBlock, tmp uint8
	for iBlock := uint8(0); iBlock < blockSize; iBlock++ {
		cBlock = uint8(rng() % blockSize)
		if iBlock != cBlock {
			for k := uint8(0); k < Size; k++ {
				for j := uint8(0); j < blockSize; j++ {
					tmp = s.Table[k][iBlock*blockSize+j]
					s.Table[k][iBlock*blockSize+j] = s.Table[k][cBlock*blockSize+j]
					s.Table[k][cBlock*blockSize+j] = tmp
				}
			}
		}
	}
}
