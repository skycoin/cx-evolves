package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/skycoin/cx-evolves/cmd/maze"
	"github.com/urfave/cli/v2"
)

func main() {
	var width int
	var height int
	var numberOfRuns int
	var randomPlayer bool
	var plotHisto bool

	mazeApp := &cli.App{
		Name:    "Maze Generator",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "random player",
				Aliases:     []string{"RP", "player"},
				Usage:       "set true if a random player will try to solve the generated maze",
				Destination: &randomPlayer,
			},
			&cli.BoolFlag{
				Name:        "plot histogram",
				Aliases:     []string{"histogram"},
				Usage:       "set true if a histogram should be plotted",
				Destination: &plotHisto,
			},
			&cli.IntFlag{
				Name:        "width",
				Aliases:     []string{"W"},
				Usage:       "width of the generated maze",
				Destination: &width,
			},
			&cli.IntFlag{
				Name:        "height",
				Aliases:     []string{"H"},
				Usage:       "height of the generated maze",
				Destination: &height,
			},
			&cli.IntFlag{
				Name:        "runs",
				Aliases:     []string{"N"},
				Usage:       "number of how many runs the player will try to finish the maze",
				Destination: &numberOfRuns,
			},
		},
		Action: func(c *cli.Context) error {
			if width == 0 && height == 0 {
				log.Error("No width and/or height specified")
				os.Exit(1)
			}

			if randomPlayer && numberOfRuns == 0 {
				log.Error("need to specify the number of runs the random player will do to solve the maze")
				os.Exit(1)
			}

			switch randomPlayer {
			case true:
				game := maze.Game{}
				game.PlotHistogram = plotHisto
				game.Init(width, height)
				game.MazeGame(numberOfRuns, nil)
			case false:
				maze.StartMaze(width, height)
			}
			return nil
		},
	}

	err := mazeApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
