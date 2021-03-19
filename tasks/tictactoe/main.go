package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var path = "~/tictatoe-output.txt"

type Game struct {
	boards     [9]string
	board      []string
	player     string
	turnNumber int
}

func PrintBoard(board []string, size int) {
	ClearScreen()
	for i, v := range board {
		if v == "" {
			fmt.Printf(" ")
		} else {
			fmt.Printf(v)
		}
		if i > 0 && (i+1)%size == 0 {
			fmt.Printf("\n")
		} else {
			fmt.Printf("|")
		}
	}
}

func ClearScreen() {
	c := exec.Command("cmd", "/c", "cls")
	c.Stdout = os.Stdout
	c.Run()
}

func CheckForWinner(b []string, n int, size int) (bool, string) {

	test := false
	i := 0

	//horizantal test
	for i < size {
		test = b[i] == b[i+1] && b[i+1] == b[i+2] && b[i] != ""
		if !test {
			i += 3
		} else {
			return true, b[i]
		}
	}
	i = 0
	//vertical test
	for i < 3 {
		test = b[i] == b[i+3] && b[i+3] == b[i+6] && b[i] != ""
		if !test {
			i += 1
		} else {
			return true, b[i]
		}
	}
	//diagonal 1 test
	if b[0] == b[4] && b[4] == b[8] && b[0] != "" {
		return true, b[i]
	}
	//diagonal 2 test
	if b[2] == b[4] && b[4] == b[6] && b[2] != "" {
		return true, b[i]
	}
	if n == 9 {
		return true, ""
	}
	return false, ""
}

func (game *Game) SwitchPlayer() {
	if game.player == "X" {
		game.player = "O"
		return
	}
	game.player = "X"
}

func (game *Game) play(position int) error {
	if game.board[position-1] == "" {
		game.board[position-1] = game.player
		game.SwitchPlayer()
		game.turnNumber += 1
		return nil
	}
	return errors.New("Try another move")
}

func askforplay() int {
	rand.Seed(time.Now().Unix())
	randomMove := rand.Intn(10)
	if randomMove == 0 {
		randomMove = 1
	}
	fmt.Println("Position to play")
	fmt.Println(randomMove)
	return randomMove
}

func askforplaymanual() int {
	var moveInt int
	fmt.Println("Enter Pos to play: ")
	fmt.Scan(&moveInt)
	return moveInt
}

/* print errors*/
func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

/*writeFile write the data into file*/
func writeFile(p []int) {

	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if isError(err) {
		return
	}
	defer file.Close()

	// write into file
	_, err = file.WriteString(fmt.Sprintln(p))
	if isError(err) {
		return
	}

	// save changes
	err = file.Sync()
	if isError(err) {
		return
	}

	//fmt.Println("==> done writing to file")
}

func main() {
	var game Game
	game.player = "X"

	gameOver := false
	var winner string

	// get the parameter from the execution flag go run --
	getSize := os.Args[1:]
	size, _ := strconv.Atoi(getSize[0])
	game.board = make([]string, (size * size))

	for gameOver != true {
		PrintBoard(game.board, size)
		move := askforplaymanual()
		err := game.play(move)
		if err != nil {
			fmt.Println(err)
			continue
		}

		//gameOver, winner = CheckForWinner(game.board, game.turnNumber)
		gameOver = false

	}

	//PrintBoard(game.board)

	if winner == "" {
		fmt.Println("it's a draw")
	} else {
		fmt.Printf("Wohoo %s is winner", winner)
	}
}
