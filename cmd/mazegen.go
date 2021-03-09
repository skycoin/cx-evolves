package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Maze cell configurations
// The paths of the maze is represented in the binary representation.
const (
	Up = 1 << iota
	Down
	Left
	Right
)

// Bit positions
const (
	bitUp = iota
	bitDown
	bitLeft
	bitRight
)

// Directions is the set of all the directions
var Directions = []int{Up, Down, Left, Right}

// The differences in the x-y coordinate
var dy = map[int]int{Up: -1, Down: 1, Left: 0, Right: 0}
var dx = map[int]int{Up: 0, Down: 0, Left: -1, Right: 1}

// Opposite directions
var Opposite = map[int]int{Up: Down, Down: Up, Left: Right, Right: Left}

type Point struct {
	X, Y int
}

// Advance the point forward by the argument direction
func (point *Point) Advance(direction int) *Point {
	return &Point{point.X + dx[direction], point.Y + dy[direction]}
}

type Maze struct {
	Width  int
	Height int
	// Each cell of maze has 4 bits (for whether there is an opening N, opening S, opening W, opening E) on the current cell
	// index=x+(y*width) each cell of maze has 4 bits
	Cells       []int
	CurrentMove int    // starts at zero, increments every move
	Goal        *Point // Goal position random
	Start       *Point // Start Position random
}

// NewMaze creates a new maze
func NewMaze(height int, width int) *Maze {
	rand.Seed(time.Now().UnixNano())

	cells := make([]int, width*height)
	start := &Point{
		X: (rand.Int() % width),
		Y: (rand.Int() % height),
	}
	return &Maze{
		Width:       width,
		Height:      height,
		Cells:       cells,
		CurrentMove: 0,
		Start:       start,
	}
}

func (maze *Maze) Generate() {
	point := maze.Start
	stack := []*Point{maze.Start}
	for len(stack) > 0 {
		for {
			point = maze.Next(point)
			if point == nil {
				break
			}
			stack = append(stack, point)
		}

		if len(stack) > 0 {
			stack = stack[:len(stack)-1] // Pop
			if len(stack) > 0 {
				point = stack[len(stack)-1]
			}
		}
	}
}

// Next advances the maze path randomly and returns the new point
func (maze *Maze) Next(point *Point) *Point {
	neighbors := maze.Neighbors(point)
	if len(neighbors) == 0 {
		return nil
	}
	direction := neighbors[rand.Int()%len(neighbors)]
	maze.Cells[point.X+(point.Y*maze.Width)] |= direction
	next := point.Advance(direction)
	maze.Cells[next.X+(next.Y*maze.Width)] |= Opposite[direction]

	maze.CurrentMove += 1
	return next
}

// Contains judges whether the argument point is inside the maze or not
func (maze *Maze) Contains(point *Point) bool {
	return 0 <= point.X && point.X < maze.Height && 0 <= point.Y && point.Y < maze.Width
}

// Neighbors gathers the nearest undecided points
func (maze *Maze) Neighbors(point *Point) (neighbors []int) {
	for _, direction := range Directions {
		next := point.Advance(direction)
		if maze.Contains(next) && maze.Cells[next.X+(next.Y*maze.Width)] == 0 {
			neighbors = append(neighbors, direction)
		}
	}
	return neighbors
}

func (maze *Maze) PrintMaze() {
	hWall := []byte("+---")
	hOpen := []byte("+   ")
	vWall := []byte("|   ")
	vOpen := []byte("    ")
	rightCorner := []byte("+\n")
	rightWall := []byte("|\n")
	var b []byte

	for y := 0; y < maze.Height; y++ {
		for z := 0; z < 3; z++ {
			for x := 0; x < maze.Width; x++ {
				switch z {
				case 0:
					// Top
					if y == 0 {
						// Top wall
						b = append(b, hWall...)
						// End of top
						if x == (maze.Width)-1 {
							b = append(b, rightCorner...)
						}
					}

					if y > 0 {
						if !IsBitSet(byte(maze.Cells[x+(y*maze.Width)]), bitUp) {
							b = append(b, hWall...)
						} else {
							b = append(b, hOpen...)
						}
						// End of top
						if x == (maze.Width)-1 {
							b = append(b, rightWall...)
						}
					}

				case 1:
					// Middle
					if x == 0 {
						b = append(b, vWall...)
					}

					if !IsBitSet(byte(maze.Cells[x+(y*maze.Width)]), bitRight) {
						// End of middle
						if x == (maze.Width)-1 {
							b = append(b, rightWall...)
						} else {
							b = append(b, vWall...)
						}
					} else {
						b = append(b, vOpen...)
					}

				case 2:
					// Bottom
					if y == (maze.Height)-1 {
						b = append(b, hWall...)
						if x == (maze.Width)-1 {
							b = append(b, rightCorner...)
						}
					}
				}
			}
		}
	}
	fmt.Print(string(b))
}

func IsBitSet(b byte, pos int) bool {
	return (b & (1 << pos)) != 0
}

func (maze *Maze) ValidateMaze() {
	fmt.Printf("Validating Maze...\n")
	var point Point
	for y := 0; y < maze.Height; y++ {
		for x := 0; x < maze.Width; x++ {
			point = Point{
				X: x,
				Y: y,
			}
			// If cell is open UP
			if IsBitSet(byte(maze.Cells[x+(y*maze.Width)]), bitUp) {
				next := point.Advance(Up)
				if maze.Contains(next) {
					// Up cell should be open down
					if !IsBitSet(byte(maze.Cells[next.X+(next.Y*maze.Width)]), bitDown) {
						panic("cells did not match")
					}
				}
			}

			// If cell is open DOWN
			if IsBitSet(byte(maze.Cells[x+(y*maze.Width)]), bitDown) {
				next := point.Advance(Down)
				if maze.Contains(next) {
					// Down cell should be open Up
					if !IsBitSet(byte(maze.Cells[next.X+(next.Y*maze.Width)]), bitUp) {
						panic("cells did not match")
					}
				}
			}

			// If cell is open LEFT
			if IsBitSet(byte(maze.Cells[x+(y*maze.Width)]), bitLeft) {
				next := point.Advance(Left)
				if maze.Contains(next) {
					// left cell should be open Right
					if !IsBitSet(byte(maze.Cells[next.X+(next.Y*maze.Width)]), bitRight) {
						panic("cells did not match")
					}
				}
			}

			// If cell is open RIGHT
			if IsBitSet(byte(maze.Cells[x+(y*maze.Width)]), bitRight) {
				next := point.Advance(Right)
				if maze.Contains(next) {
					// right cell should be open Left
					if !IsBitSet(byte(maze.Cells[next.X+(next.Y*maze.Width)]), bitLeft) {
						panic("cells did not match")
					}
				}
			}
		}
	}
	fmt.Printf("Finished...\n")
	fmt.Printf("Maze is valid.\n")
}

func main() {
	var width string
	var height string
	var intWidth int
	var intHeight int

	mazeApp := &cli.App{
		Name:    "Maze Generator",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "width",
				Value:       "",
				Aliases:     []string{"W"},
				Usage:       "width of the generated maze",
				Destination: &width,
			},
			&cli.StringFlag{
				Name:        "height",
				Value:       "",
				Aliases:     []string{"H"},
				Usage:       "height of the generated maze",
				Destination: &height,
			},
		},
		Action: func(c *cli.Context) error {
			if width != "" && height != "" {
				intWidth, _ = strconv.Atoi(width)
				intHeight, _ = strconv.Atoi(height)
				startMaze(intWidth, intHeight)
			} else {
				log.Error("No width and/or height specified")
				os.Exit(1)
			}
			return nil
		},
	}

	err := mazeApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func startMaze(width, height int) {
	maze := NewMaze(height, width)
	maze.Generate()
	maze.ValidateMaze()
	fmt.Printf("Starting Point (x,y)=(%v,%v)\n", maze.Start.X, maze.Start.Y)
	fmt.Printf("Number of moves=%v\n", maze.CurrentMove)
	maze.PrintMaze()
}
