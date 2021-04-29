#!/bin/bash
  
# turn on bash's job control
set -m
  
# Start worker and put it in background
./server -workers=3 &
  
# Start the main process
./cx-evolves --maze=true --name=MazeRunner -W=2 -H=2 --random-maze-size=false --population=100 --generations=100 --expressions=300 --epoch-length=100 --graphs=true --ast=true --use-log-fitness=false --workers=3