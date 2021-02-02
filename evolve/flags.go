package main

import (
	"flag"
	"fmt"
	"os"
)

const VERSION = "1.0"

type flags struct {
	iterations        int
	printVersion      bool
}

func defaultCmdFlags() flags {
	return flags{
		iterations:        0,
		printVersion:      false,
	}
}

var commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

func parseFlags(options *flags, args []string) {
	commandLine.IntVar(&options.iterations, "iterations", options.iterations, "Iterations")
	commandLine.BoolVar(&options.printVersion, "version", options.printVersion, "Version")

	commandLine.Parse(args)
}

func printVersion() {
	fmt.Println("Sequence Predictors v", VERSION)
}
