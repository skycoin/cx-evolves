### Summary

This program generates, validates, and prints maze. The algorithm use in generating the maze is recursive backtracker that uses a stack. The goal point is determined by using djikstra's algo to find the farthest point from the starting point.

### Usage

    go run mazegen.go -W=[width of generated maze] -H=[height of generated maze]
    go run mazegen.go help 

### Example
   go run mazgen.go -W=20 -H=20

### Notes
1. If no arguments are specified, the program will return an error.

2. If width and height values are specified, the program will run.
