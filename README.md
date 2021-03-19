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
go run main.go 
--name=[Name of generated program]
--maze=[set true if benchmark evolve with maze]
-W=[width of the maze to solve] 
-H=[height of maze to solve]  
--random-maze-size=[set true if generated maze for every N epochs will have random sizes (NxN 2,3,4,5,6,7, or 8)]
--epoch-length=[Maze changes for every N generation(Example. Epochs=5, maze changes every 5 generations)]
--constants=[set true if benchmark evolve with constants]
--evens=[set true if benchmark evolve with evens]
--odds=[set true if benchmark evolve with odds]
--primes=[set true if benchmark evolve with primes]
--range=[set true if benchmark evolve with range]
--upper-range=[upper range (int) for range benchmark]
--lower-range=[lower range (int) for range benchmark]
--rounds=[Number of rounds per program for benchmarks contants, evens, odds, primes, and range]
--population=[Population size or number of programs generated per generation]
--generations=[number of generations]
--expressions=[number of expressions a generated program can have]
--graphs=[Set true if average fitness and fittest per generation graphs will be saved] 
--ast=[Set true if best ASTs per generation will be saved]
--use-log-fitness=true 
```

### Example Benchmarking

For Maze
```
go run main.go --maze=true --name=MazeRunner -W=2 -H=2 --random-maze-size=false --population=50 --generations=100 --expressions=30 --epoch-length=100 --graphs=true --ast=false --use-log-fitness=false
```

For Range
```
go run main.go --range=true --upper-range=9 --lower-range=2 --rounds=10 --name=Range --population=300 --generations=1000 --expressions=30 --graphs=true --ast=false --use-log-fitness=false
```

For Constants
```
go run main.go --constants=true --rounds=10 --name=Constants --population=300 --generations=1000 --expressions=30 --graphs=true --ast=false --use-log-fitness=false
```

For Odds
```
go run main.go --odds=true --rounds=10 --name=Odds --population=300 --generations=1000 --expressions=30 --graphs=true --ast=false --use-log-fitness=false
```

For Evens
```
go run main.go --evens=true --rounds=10 --name=Evens --population=300 --generations=1000 --expressions=30 --graphs=true --ast=false --use-log-fitness=false
```

For Primes
```
go run main.go --primes=true --rounds=10 --name=Primes --population=300 --generations=1000 --expressions=30 --graphs=true --ast=false --use-log-fitness=false
```

### Notes
1. If no arguments are specified, the program will run default values.
