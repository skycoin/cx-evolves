#!/bin/bash
  
# turn on bash's job control
set -m
  
if [ "$WORKERS" == "" ]
then
	WORKERS=500
fi

# Start workers and put it in background
basePort=9090
i=0
while [ "$i" -lt $WORKERS ]; do
    ./server -port=$(($basePort+$i)) &
    i=$(( i + 1 ))
done 

  
  
# Start the main process
./cx-evolves --task=constants --rounds=10 --constants-target=5 --population=500 --generations=20000 --expressions=100 --graphs=true --ast=false --use-log-fitness=false --workers=$WORKERS