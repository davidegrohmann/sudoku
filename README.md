Sudoku Solver
=============

Sudoku solver based on Dacing Links algorithm (see [paper](https://arxiv.org/abs/cs/0011047https://arxiv.org/abs/cs/0011047))

Build
-----

To build the project for linux/amd64, simply run
```shell script
make
```

To clean all produced objects, run
```shell script
make clean
```

Run
---

The produced program can be found in the `cmd` folder.

For example, try to execute the following:
```shell script
cmd/sudoku.bin --verbose --file examples/sudoku.txt
```

Instruction on how to use the program can be printed by
```shell script
cmd/sudoku.bin --help
```