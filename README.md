Sudoku Solver
=============

Sudoku solver based on Dacing Links algorithm (see [paper](https://arxiv.org/abs/cs/0011047))

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

For example, try to execute the following for solving the example sudoku puzzle:
```shell script
cmd/sudoku.bin --verbose --file examples/sudoku.txt
```

For creating puzzle, run the program with the `create` command (very slow!):
```shell script
cmd/sudoku.bin --cmd create
```

Instruction on how to use the program can be printed by
```shell script
cmd/sudoku.bin --help
```