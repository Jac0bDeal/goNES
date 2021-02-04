# goNES
A NES emulator in pure Go.

Currently only the 6502 cpu and disassembler are implemented. The main 
program currently loads a hardcoded test program into memory, disassembles
it, and prints out the disassembled program to the terminal. 

![build-status](https://github.com/Jac0bDeal/goNES/workflows/Build/badge.svg?branch=main)
![tests-status](https://github.com/Jac0bDeal/goNES/workflows/Tests/badge.svg?branch=main)
[![go-report-card](https://goreportcard.com/badge/github.com/Jac0bDeal/goNES)](https://goreportcard.com/report/github.com/Jac0bDeal/goNES)

## Build
In order to build this, you need Go 1.14+ and Make installed.

From the project root, run
```shell script
make goNES
```

This will build the binary at `bin/goNES`.

## Running
Once built, the binary is run with
```shell script
./bin/goNES
```

## Tests
If you want to run the tests (for some reason) use
```shell script
make test
```
