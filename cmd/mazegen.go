package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/RyanCarrier/dijkstra"
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

// Distance between cells.
// Used for Dijkstra's algo for finding farthest point from goal point.
const (
	distanceBetweenPoints = 1
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

// NewMaze creates a new maze.
// The starting point is set randomly.
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

//  Generate generates a maze.
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
	maze.Cells[maze.getIndex(point.X, point.Y)] |= direction
	next := point.Advance(direction)
	maze.Cells[maze.getIndex(next.X, next.Y)] |= Opposite[direction]

	return next
}

// Contains judges whether the argument point is inside the maze or not
func (maze *Maze) Contains(point *Point) bool {
	return 0 <= point.X && point.X < maze.Width && 0 <= point.Y && point.Y < maze.Height
}

// Neighbors gathers the nearest undecided points
func (maze *Maze) Neighbors(point *Point) (neighbors []int) {
	for _, direction := range Directions {
		next := point.Advance(direction)
		if maze.Contains(next) && maze.Cells[maze.getIndex(next.X, next.Y)] == 0 {
			neighbors = append(neighbors, direction)
		}
	}
	return neighbors
}

// PrintMaze prints the entire maze onto cli with +,-, *space*, and | characters.
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
						if !IsBitSet(byte(maze.Cells[maze.getIndex(x, y)]), bitUp) {
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

					if !IsBitSet(byte(maze.Cells[maze.getIndex(x, y)]), bitRight) {
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

// IsBitSet checks if the bit in b located in pos is true.
func IsBitSet(b byte, pos int) bool {
	return (b & (1 << pos)) != 0
}

// Validates the maze.
// If the cell is open to the NORTH, the cell on its NORTH must be open to the SOUTH,
// if the cell is open to the SOUTH, the cell on its SOUTH must be open to the NORTH.
// if the cell is open to the EAST, the cell on its EAST must be open to the WEST, and
// if the cell is open to the WEST, the cell on its WEST must be open to the EAST,
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
			if IsBitSet(byte(maze.Cells[maze.getIndex(x, y)]), bitUp) {
				next := point.Advance(Up)
				if maze.Contains(next) {
					// Up cell should be open down
					if !IsBitSet(byte(maze.Cells[maze.getIndex(next.X, next.Y)]), bitDown) {
						panic("cells did not match")
					}
				}
			}

			// If cell is open DOWN
			if IsBitSet(byte(maze.Cells[maze.getIndex(x, y)]), bitDown) {
				next := point.Advance(Down)
				if maze.Contains(next) {
					// Down cell should be open Up
					if !IsBitSet(byte(maze.Cells[maze.getIndex(next.X, next.Y)]), bitUp) {
						panic("cells did not match")
					}
				}
			}

			// If cell is open LEFT
			if IsBitSet(byte(maze.Cells[maze.getIndex(x, y)]), bitLeft) {
				next := point.Advance(Left)
				if maze.Contains(next) {
					// left cell should be open Right
					if !IsBitSet(byte(maze.Cells[maze.getIndex(next.X, next.Y)]), bitRight) {
						panic("cells did not match")
					}
				}
			}

			// If cell is open RIGHT
			if IsBitSet(byte(maze.Cells[maze.getIndex(x, y)]), bitRight) {
				next := point.Advance(Right)
				if maze.Contains(next) {
					// right cell should be open Left
					if !IsBitSet(byte(maze.Cells[maze.getIndex(next.X, next.Y)]), bitLeft) {
						panic("cells did not match")
					}
				}
			}
		}
	}
	fmt.Printf("Finished...\n")
	fmt.Printf("Maze is valid.\n")
}

// SetGoalPoint sets maze goal point by using djikstra's algorithm to find farthest
// point from the starting point.
func (maze *Maze) SetGoalPoint() {
	graph := dijkstra.NewGraph()

	// Make vertex for all indexes
	for i := 0; i < maze.Height*maze.Width; i++ {
		graph.AddVertex(i)
	}

	// Add arcs to all points
	for h := 0; h < maze.Height; h++ {
		for w := 0; w < maze.Width; w++ {
			point := Point{
				X: w,
				Y: h,
			}

			// check up
			up := point.Advance(Up)
			if maze.Contains(up) {
				err := graph.AddArc(maze.getIndex(point.X, point.Y), maze.getIndex(up.X, up.Y), distanceBetweenPoints)
				if err != nil {
					log.Fatal(err)
				}
			}

			// check down
			down := point.Advance(Down)
			if maze.Contains(down) {
				err := graph.AddArc(maze.getIndex(point.X, point.Y), maze.getIndex(down.X, down.Y), distanceBetweenPoints)
				if err != nil {
					log.Fatal(err)
				}
			}

			// check right
			right := point.Advance(Right)
			if maze.Contains(right) {
				err := graph.AddArc(maze.getIndex(point.X, point.Y), maze.getIndex(right.X, right.Y), distanceBetweenPoints)
				if err != nil {
					log.Fatal(err)
				}
			}

			// check left
			left := point.Advance(Left)
			if maze.Contains(left) {
				err := graph.AddArc(maze.getIndex(point.X, point.Y), maze.getIndex(left.X, left.Y), distanceBetweenPoints)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	startPoint := maze.getIndex(maze.Start.X, maze.Start.Y)
	var longestDistance int64
	var farthestPoint Point
	for h := 0; h < maze.Height; h++ {
		for w := 0; w < maze.Width; w++ {
			if startPoint != maze.getIndex(w, h) {
				best, err := graph.Shortest(startPoint, maze.getIndex(w, h))
				if err != nil {
					log.Fatal(err)
				}

				if best.Distance > longestDistance {
					longestDistance = best.Distance
					farthestPoint = Point{
						X: w,
						Y: h,
					}
				}
			}

		}
	}
	maze.Goal = &farthestPoint
}

func (maze *Maze) getIndex(x, y int) int {
	return x + (y * maze.Width)
}

func startMaze(width, height int) {
	maze := NewMaze(height, width)
	maze.Generate()
	maze.ValidateMaze()
	maze.SetGoalPoint()
	fmt.Printf("Start Point (x,y)=(%v,%v)\nGoal Point (x,y)=(%v,%v)\n", maze.Start.X, maze.Start.Y, maze.Goal.X, maze.Goal.Y)
	maze.PrintMaze()
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
