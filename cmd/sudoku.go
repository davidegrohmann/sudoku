package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sudoku/pkg/sudoku"
	"sudoku/pkg/xorwow"
	"time"
)

const create = "create"
const solve = "solve"

var buildVersion = "<Version Not Set>"
var buildTime = "<Time Not Set>"

func main() {
	os.Exit(mainWithExit())
}

func mainWithExit() int {
	var file, cmd string
	var verbose, version bool
	flag.StringVar(&cmd, "cmd", "solve", "cmd to execute")
	flag.StringVar(&file, "file", "", "input file (if no file is given stdin is used)")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.BoolVar(&version, "version", false, "print version")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "%s [-cmd CMD] [-file FILE] [-verbose] [-version]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if version {
		_, _ = fmt.Fprintf(os.Stdout, "Build version=%s and time=%s\n", buildVersion, buildTime)
		return 0
	}

	switch cmd {
	case create:
		return createSudoku(verbose)
	case solve:
		r := os.Stdin
		if file != "" {
			var err error
			r, err = os.Open(file)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Could not open file %v: %v\n", file, err)
				return 1
			}
			defer func() {
				err := r.Close()
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Could not close file %v: %v\n", file, err)
				}
			}()
		}
		return solveSudoku(r, verbose)
	default:
		_, _ = fmt.Fprintf(os.Stderr, "Unknown cmd '%s'\n", cmd)
		return 3
	}
}

func randomShuffles(s *sudoku.Sudoku, rng xorwow.RNG, verbose bool) {
	shuffles := []func(){
		func() {
			s.RandomShuffleRows(rng.Rand)
			if verbose {
				_, _ = fmt.Fprint(os.Stderr, "Shuffle rows:\n")
				_ = s.Print(os.Stderr)
			}
		},
		func() {
			s.RandomShuffleColumns(rng.Rand)
			if verbose {
				_, _ = fmt.Fprint(os.Stderr, "Shuffle columns:\n")
				_ = s.Print(os.Stderr)
			}
		}, func() {
			s.RandomShuffleBlockRows(rng.Rand)
			if verbose {
				_, _ = fmt.Fprint(os.Stderr, "Shuffle block rows:\n")
				_ = s.Print(os.Stderr)
			}
		},
		func() {
			s.RandomShuffleBlockColumns(rng.Rand)
			if verbose {
				_, _ = fmt.Fprint(os.Stderr, "Shuffle block columns:\n")
				_ = s.Print(os.Stderr)
			}
		},
	}

	length := uint32(len(shuffles))
	for i := uint32(0); i < length; i++ {
		tmp := rng.Rand() % length
		for shuffles[tmp] == nil {
			tmp = (tmp + 1) % length
		}
		shuffles[tmp]()
		shuffles[tmp] = nil
	}
}

func createSudoku(verbose bool) int {
	rng := xorwow.NewRNG(uint32(time.Now().UnixNano()))
	s := &sudoku.Sudoku{}

	// optimization: fill first 9 spots in the sudoku with numbers from 1 to 9
	s.RandomFill(rng.Rand)

	// create a solved puzzle
	_ = s.Solve()
	if verbose {
		_, _ = fmt.Fprint(os.Stderr, "Solved to start with:\n")
		_ = s.Print(os.Stderr)
	}

	randomShuffles(s, rng, verbose)

	// try to remove cells until no more removals are possible
	var tmp uint8
	var row, column uint8
	var ptr *uint8
	attempt := 0
	for attempt < 100 {
		row = uint8(rng.Rand() % sudoku.Size)
		column = uint8(rng.Rand() % sudoku.Size)
		ptr = &(s.Table[row][column])
		if *ptr != 0 {
			tmp = *ptr
			*ptr = 0
			if s.Check() {
				attempt = 0
				if verbose {
					_, _ = fmt.Fprint(os.Stderr, "Intermediate:\n")
					_ = s.Print(os.Stderr)
				}
			} else {
				*ptr = tmp
				attempt++
			}
		}
	}

	_, _ = fmt.Fprint(os.Stdout, "New puzzle:\n")
	_ = s.Print(os.Stdout)

	return 0
}

func solveSudoku(reader io.Reader, verbose bool) (result int) {
	s := &sudoku.Sudoku{}
	err := parseSudoku(reader, s)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Bad input: %v\n", err)
		return 2
	}

	if verbose {
		_, _ = fmt.Fprint(os.Stderr, "Input puzzle:\n")
		_ = s.Print(os.Stderr)
	}

	success := s.Solve()

	if success {
		_, _ = fmt.Fprint(os.Stdout, "Solved puzzle:\n")
		_ = s.Print(os.Stdout)
		result = 0
	} else {
		_, _ = fmt.Fprint(os.Stderr, "Unsolvable puzzle:\n")
		_ = s.Print(os.Stderr)
		result = 127
	}
	return result
}

func parseSudoku(reader io.Reader, s *sudoku.Sudoku) error {
	row := 1
	column := 1
	buf := make([]byte, 1)
	var n int
	var num uint8
	var err error
	for true {
		n, err = reader.Read(buf)
		switch {
		case err == io.EOF:
			if row > sudoku.Size || (row == sudoku.Size && column == sudoku.Size) {
				return nil
			}
			return fmt.Errorf("malformed input reached EOF before parsing the whole board [%d, %d]", row, column)
		case err != nil:
			return err
		case n == 0 || buf[0] == '\r':
			// continue
		case buf[0] == '\n':
			if column < sudoku.Size {
				return fmt.Errorf("malformed input one row number %d has length %d", row, column)
			}
			row++
			column = 1
		default:
			num, err = parseChar(buf[0])
			if err != nil {
				return err
			}
			s.Table[row-1][column-1] = num
			column++
		}
	}
	return nil
}

func parseChar(b byte) (uint8, error) {
	// allow ' ' and '.'  as special marker for empty cell in input
	if b == '.' || b == ' ' {
		return 0, nil
	}

	v := b - '0'
	if v < 0 || v > 9 {
		return 0, fmt.Errorf("invalid char %v", b)
	}

	return v, nil
}
