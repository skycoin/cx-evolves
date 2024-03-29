# Installation

- Install and setup go 1.14 and above.
- Install CX (https://github.com/skycoin/cx).
- `git pull` your CX repository so you have the latest CX version.
- Make sure `go module` is enabled. Install go dependencies by running:
```
$go mod verify 
$go mod tidy
$go mod download
```

# Benchmarks
| Task Name              | Task Code Name      | Task Version |
| ---------------------- |:-------------------:|:------------:| 
|  Maze                  | "maze"              |      1       |
|  Constants             | "constants"         |      1       | 
|  Evens                 | "evens"             |      1       |
|  Evens greater than 0  | "evens"             |      2       |
|  Odds                  | "odds"              |      1       |
|  Primes                | "primes"            |      1       |
|  Composites            | "composites"        |      1       |
|  Network Simulator     | "network_simulator" |      1       |
|  Range                 | "range"             |      1       |


### Summary

This program runs evolve with Maze. Generated programs are given multiple i32 as inputs and outputs one i32. The output is their move for the maze. 

### Usage
For more information, run
```
go run main.go help 
```

```
go run main.go 
--task=[Name of task to benchmark (in lower case--refer to task code name above)]
--task-version=[Task vesion value, default is 1]
--population=[Population size or number of programs generated per generation]
--generations=[number of generations]
--expressions=[number of expressions a generated program can have]
--graphs=[Set true if average fitness and fittest per generation graphs will be saved] 
--ast=[Set true if best ASTs per generation will be saved]
--use-log-fitness=[Set true if fitness will be fitness(log base 2)]
--workers=[number of workers available to use, if not set default is 1]
--random-search=[set true to have no mutation on individuals]
--select-rank-cufoff=[set true if selection is select, rank, and cutoff]

<!-- Set if maze benchmark -->
-W=[width of the maze to solve] 
-H=[height of maze to solve]  
--random-maze-size=[set true if generated maze for every N epochs will have random sizes (NxN 2,3,4,5,6,7, or 8)]
--epoch-length=[Maze changes for every N generation(Example. Epochs=5, maze changes every 5 generations)]

<!-- Set if range benchmark -->
--upper-range=[upper range (int) for range benchmark]
--lower-range=[lower range (int) for range benchmark]

<!-- Set if benchmark is either constants, evens, odds, primes, composites, range, or network simulator -->
--rounds=[Number of rounds per program]

<!-- Set if benchmark is constants -->
--constants-target=[target number for constants benchmark]

```

### Example Benchmarking

For Maze
```
go run main.go --task=maze -W=2 -H=2 --random-maze-size=false --population=50 --generations=100 --expressions=30 --epoch-length=100 --graphs=true --ast=false --use-log-fitness=false --workers=1
```

For Range
```
go run main.go --task=range --upper-range=9 --lower-range=2 --rounds=10 --population=300 --generations=1000 --expressions=30 --graphs=true --ast=false --use-log-fitness=false --workers=1
```

For Constants
```
go run main.go --task=constants --rounds=10 --population=300 --generations=1000 --expressions=30 --graphs=true --ast=false --use-log-fitness=false --workers=1
```

For Constants with explicit target
```
go run main.go --task=constants --rounds=10 --constants-target=5 --population=300 --generations=1000 --expressions=30 --graphs=true --ast=false --use-log-fitness=false --workers=1

Note: For constants benchmarking 0-256 range, use --constants-target=256.
```

For Odds
```
go run main.go --task=odds --rounds=10 --population=300 --generations=1000 --expressions=30 --graphs=true --ast=false --use-log-fitness=false --workers=1
```

For Evens
```
go run main.go --task=evens --rounds=10 --population=300 --generations=1000 --expressions=30 --graphs=true --ast=false --use-log-fitness=false --workers=1
```

For Primes
```
go run main.go --tasks=primes --rounds=10 --population=300 --generations=1000 --expressions=30 --graphs=true --ast=false --use-log-fitness=false --workers=1
```

For Composites
```
go run main.go --task=composites --rounds=10 --population=300 --generations=1000 --expressions=30 --graphs=true --ast=false --use-log-fitness=false --workers=1
```

For Network Simulator
```
 go run main.go --task=network_simulator --rounds=20 --population=300 --generations=2000 --expressions=30 --graphs=true --ast=false --use-log-fitness=false --workers=1  
```

### Notes
1. If no arguments are specified, the program will run default values.

# Benchmarking in Docker Env

First edit the ./scripts/maze_benchmark.sh to the benchmark options needed and the number of workers to deploy.
Then use this command:

To build and run
```
NAME=[Name for the docker image]
MOUNT=[Local directory to mount the "Result" directory of the container, this is where the graphs and asts can be found] 
make deploy 
```

Example
```
NAME=testbenchmark MOUNT=/Benchmarking/Results make deploy 
```

To build and push to dockerhub
```
NAME=[Name for the docker image]
make push-docker
```

Example
```
NAME=testbenchmark make push-docker
```

# High priority tasks
- [x] Maze
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
- [x] Check if a number is prime