package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

type Game struct {
	board      [9]string
	player     string
	turnNumber int
}

func PrintBoard(board [9]string) {
	ClearScreen()
	for i, v := range board {
		if v == "" {
			fmt.Printf(" ")
		} else {
			fmt.Printf(v)
		}
		if i > 0 && (i+1)%3 == 0 {
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

func CheckForWinner(b [9]string, n int) (bool, string) {

	test := false
	i := 0

	//horizantel test
	for i < 9 {
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

func main() {
	var game Game
	game.player = "X"

	gameOver := false
	var winner string

	for gameOver != true {
		PrintBoard(game.board)
		move := askforplay()
		err := game.play(move)
		if err != nil {
			fmt.Println(err)
			continue
		}

		gameOver, winner = CheckForWinner(game.board, game.turnNumber)

	}

	PrintBoard(game.board)

	if winner == "" {
		fmt.Println("it's a draw")
	} else {
		fmt.Printf("Wohoo %s is winner", winner)
	}
}
