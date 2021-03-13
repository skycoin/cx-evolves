### Summary

This program generates, validates, and prints a maze. The algorithm used in generating the maze is recursive backtracker that uses a stack. The goal point is determined by using djikstra's algo to find the farthest point from the starting point.

### Usage

    go run main.go -W=[width of generated maze] -H=[height of generated maze] -player=[True if a random player will try to solve the maze] -runs=[# of runs the random player will try to solve the maze] -histogram=[true if you want to generate histogram for the result]
    go run main.go help 

### Example
   go run main.go -W=20 -H=20 -player=true -checkprev=true -runs=3 -histogram=true

### Notes
1. If no arguments are specified, the program will return an error.

2. If width and height values are specified, the program will run.

### Time to Generate Maze
16x16 Maze
| Trial Number  | Time          |
| ------------- |:-------------:| 
| 1             | 96.32µs       |
| 2             | 104.208µs     |  
| 3             | 96.414µs      |
| 4             | 140.279µs     |
| 5             | 89.953µs      |
| Average       | 105.435µs     |

32x32 Maze
| Trial Number  | Time          |
| ------------- |:-------------:| 
| 1             | 332.848µs     |
| 2             | 341.459µs     |  
| 3             | 361.315µs     |
| 4             | 353.017µs     |
| 5             | 355.363µs     |
| Average       | 348.800µs     |

64x64 Maze
| Trial Number  | Time          |
| ------------- |:-------------:| 
| 1             | 1.216691ms    |
| 2             | 1.38822ms     |  
| 3             | 1.416075ms    |
| 4             | 1.646954ms    |
| 5             | 1.436125ms    |
| Average       | 1.420813ms    |

124x124 Maze
| Trial Number  | Time          |
| ------------- |:-------------:| 
| 1             | 5.419467ms    |
| 2             | 4.550927ms    |  
| 3             | 5.019394ms    |
| 4             | 5.085402ms    |
| 5             | 5.305665ms    |
| Average       | 5.076171ms    |

### Time finding furthest point from starting point
16x16 Maze
| Trial Number  | Time          |
| ------------- |:-------------:| 
| 1             | 32.317µs      |
| 2             | 36.64µs       |  
| 3             | 29.159µs      |
| 4             | 36.156µs      |
| 5             | 30.898µs      |
| Average       | 33.034µs      |

32x32 Maze
| Trial Number  | Time          |
| ------------- |:-------------:| 
| 1             | 188.672µs     |
| 2             | 218.576µs     |  
| 3             | 230.086µs     |
| 4             | 224.502µs     |
| 5             | 223.16µs      |
| Average       | 216.999µs     |

64x64 Maze
| Trial Number  | Time          |
| ------------- |:-------------:| 
| 1             | 0.91057ms     |
| 2             | 1.024449ms    |  
| 3             | 1.040892ms    |
| 4             | 1.179318ms    |
| 5             | 1.057039ms    |
| Average       | 1.0424536ms   |

124x124 Maze
| Trial Number  | Time          |
| ------------- |:-------------:| 
| 1             | 2.781328ms    |
| 2             | 2.767195ms    |  
| 3             | 2.521474ms    |
| 4             | 2.490726ms    |
| 5             | 2.758457ms    |
| Average       | 2.65119696ms  |

