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

To run the hello-world example:

- `go run main.go`

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
