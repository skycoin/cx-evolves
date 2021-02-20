package main

import (
	"fmt"
	"math/rand"
	"time"
	"unicode"

	"github.com/nsf/termbox-go"

	"github.com/itchyny/maze"
)

type keyDire struct {
	key  termbox.Key
	char rune
	dir  int
}

var keyDirss = []*keyDire{
	{termbox.KeyArrowUp, 'k', maze.Up},
	{termbox.KeyArrowDown, 'j', maze.Down},
	{termbox.KeyArrowLeft, 'h', maze.Left},
	{termbox.KeyArrowRight, 'l', maze.Right},
}

func auto(maze *maze.Maze, format *maze.Format) {
	rand.Seed(time.Now().Unix())
	eventsa := make(chan termbox.Event)
	go func() {
		for {
			eventsa <- termbox.PollEvent()
		}
	}()
	move := make(chan string)
	go func() {
		move <- "go"
	}()

	strwriter := make(chan string)
	ticker := time.NewTicker(10 * time.Millisecond)
	go printTermboxa(maze, strwriter, time.Now())
	maze.Started = true
	maze.Write(strwriter, format)
loop:
	for {
		select {
		case <-move:
			if !maze.Finished {
				randomMove := rand.Intn(4) //rand int n de keyDirs
				keydir := keyDirss[randomMove]
				if keydir.key != 00 {
					maze.Move(keydir.dir)

					if maze.Finished {
						maze.Solve()
					}

					maze.Write(strwriter, format)
					go func() {
						move <- "go"
					}()
					continue
				}
			}
		case event := <-eventsa:
			if event.Ch == 'q' || event.Ch == 'Q' || event.Key == termbox.KeyCtrlC || event.Key == termbox.KeyCtrlD {
				break loop
			}
		case <-ticker.C:
			if !maze.Finished {
				strwriter <- "\u0000"
			}
		}
	}
	ticker.Stop()
}

func printTermboxa(maze *maze.Maze, strwriter chan string, start time.Time) {
	x, y := 1, 0
	for {
		str := <-strwriter
		switch str {
		case "\u0000":
			printFinisheda(maze, time.Now().Sub(start))
			termbox.Flush()
			x, y = 1, 0
		default:
			printStringa(str, &x, &y)
		}
	}
}

func printStringa(str string, x *int, y *int) {
	attr, skip, d0, d1, d := false, false, '0', '0', false
	fg, bg := termbox.ColorDefault, termbox.ColorDefault
	for _, c := range str {
		if c == '\n' {
			*x, *y = (*x)+1, 0
		} else if c == '\x1b' || attr && c == '[' {
			attr = true
		} else if attr && unicode.IsDigit(c) {
			if !skip {
				if d {
					d1 = c
				} else {
					d0, d = c, true
				}
			}
		} else if attr && c == ';' {
			skip = true
		} else if attr && c == 'm' {
			if d0 == '7' && d1 == '0' {
				fg, bg = termbox.AttrReverse, termbox.AttrReverse
			} else if d0 == '3' {
				fg, bg = termbox.Attribute(uint64(d1-'0'+1)), termbox.ColorDefault
			} else if d0 == '4' {
				fg, bg = termbox.ColorDefault, termbox.Attribute(uint64(d1-'0'+1))
			} else {
				fg, bg = termbox.ColorDefault, termbox.ColorDefault
			}
			attr, skip, d0, d1, d = false, false, '0', '0', false
		} else {
			termbox.SetCell(*y, *x, c, fg, bg)
			*y = *y + 1
		}
	}
}

func printFinisheda(maze *maze.Maze, duration time.Duration) {
	str := fmt.Sprintf("%8d.%02ds      ", int64(duration/time.Second), int64((duration%time.Second)/1e7))
	fg, bg := termbox.ColorDefault, termbox.ColorDefault
	if maze.Finished {
		x, y := maze.Height, 2*maze.Width-6
		if y < 0 {
			y = 0
		}
		for j, s := range []string{
			"                  ",
			"    Finished!     ",
			str,
			"                  ",
			"  Press q to quit ",
			"                  ",
		} {
			for i, c := range s {
				termbox.SetCell(y+i, x+j, c, fg, bg)
			}
		}
	} else {
		for i, c := range str {
			termbox.SetCell(4*maze.Width+i-8, 1, c, fg, bg)
		}
	}
}
