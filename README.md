# Installation

- Install and setup go 1.13 and above.
- Install CX (https://github.com/skycoin/cx).
- `git pull` your CX repository so you have the latest CX version.
- Make sure `go module` is enabled. Install go dependencies by running:
```
$go get -u ./...
$go mod verify 
$go mod vendor
```

# High priority tasks
- [ ] Maze
- [ ] Snake
- [ ] Tic-Tac-Toe
- [ ] Multi-armed bandit
- [ ] Rock, paper & scissors

# Low priority tasks
- [ ] Time-series mean estimation
- [ ] Time-series variance estimation
- [ ] Time-series covariance estimation
- [ ] Stock market simulator
- [ ] Committee of experts
- [ ] Linear feedback shift register prediction
- [ ] Non-linear feedback shift register prediction
- [ ] Check if a number is prime

### Summary

This program runs evolve with Maze. Generated programs are given multiple i32 as inputs and outputs one i32. The output is their move for the maze. 

### Usage
For more information, run
```
go run main.go help 
```

```
go run main.go -W=[width of the maze to solve] 
-H=[height of maze to solve]  
--random=[set true if generated maze for every N epochs will have random sizes (NxN 2,3,4,5,6,7, or 8)]
--population=[Population size or number of programs generated per generation]
--generations=[number of generations]
--expressions=[number fo expressions a generated program can have]
--epoch-length=[Maze changes for every N generation(Example. Epochs=5, maze changes every 5 generations)]
--name=[Name of generated program]
--graphs=[Set true if average fitness and fittest per generation graphs will be saved] 
--ast=[Set true if best ASTs per generation will be saved]
```

### Example
```
go run main.go -W=2 -H=2 --random=false --population=50 --generations=100 --expressions=30 --epoch-length=100 --name=MazeRunner --graphs=true --ast=false
```

### Notes
1. If no arguments are specified, the program will run default values.
