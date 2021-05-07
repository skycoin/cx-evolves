#!/bin/bash
  
# turn on bash's job control
set -m

if [ "$WORKERS" == "" ]
then
	WORKERS=100
fi

# Start workers and put it in background
basePort=9090
i=0
while [ "$i" -lt $WORKERS ]; do
    ./server -port=$(($basePort+$i)) &
    i=$(( i + 1 ))
done 

  
# Start the main process
./cx-evolves --maze=true --name=MazeRunner -W=2 -H=2 --random-maze-size=false --population=100 --generations=1000 --expressions=50 --epoch-length=10 --graphs=true --ast=false --use-log-fitness=false --workers=$WORKERS