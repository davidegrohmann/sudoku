package sudoku

const rowColumnNumberEntries = Size * Size * Size
const constraints = 4
const constraintMatrixSize = Size * Size
const constraintEntries = constraints * constraintMatrixSize

// map from line entry in the set cover map to the row in the sudoku
func toSudokuRow(x uint32) uint8 {
	result := x / 81
	return uint8(result)
}

// map from line entry in the set cover map to the column in the sudoku
func toSudokuColumn(x uint32) uint8 {
	result := (x % 81) / 9
	return uint8(result)
}

// map from line entry in the set cover map to the value of the cell
func toSudokuNumber(x uint32) uint8 {
	result := (x % 9) + 1
	return uint8(result)
}

const rowColumnConstraintOffset = 0 * constraintMatrixSize
const rowNumberConstraintOffset = 1 * constraintMatrixSize
const columnNumberConstraintOffset = 2 * constraintMatrixSize
const boxNumberConstraintOffset = 3 * constraintMatrixSize

func toConstraintLineEntry(row, column, number uint8) uint32 {
	return (uint32(row) * 81) + (uint32(column) * 9) + (uint32(number) - 1)
}

func toConstraintRowColumn(row, column uint8) uint32 {
	return rowColumnConstraintOffset + (uint32(row) * 9) + uint32(column)
}

func toConstraintRowNumber(row, number uint8) uint32 {
	return rowNumberConstraintOffset + (uint32(row) * 9) + (uint32(number) - 1)
}

func toConstraintColumnNumber(column, number uint8) uint32 {
	return columnNumberConstraintOffset + (uint32(column) * 9) + (uint32(number) - 1)
}

func toConstraintBoxNumber(row, column, number uint8) uint32 {
	box := ((uint32(row) / 3) * 3) + (uint32(column) / 3)
	return boxNumberConstraintOffset + (box * 9) + (uint32(number) - 1)
}
