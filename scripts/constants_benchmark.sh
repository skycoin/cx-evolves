#!/bin/bash
  
# turn on bash's job control
# set -m
  
# if [ "$WORKERS" == "" ]
# then
# 	WORKERS=100
# fi

# # Start workers and put it in background
# basePort=9090
# i=0
# while [ "$i" -lt $WORKERS ]; do
#     ./server -port=$(($basePort+$i)) &
#     i=$(( i + 1 ))
# done 

  
  
# # Start the main process
# ./cx-evolves --constants=true --rounds=10 --name=Constants --population=100 --generations=1000 --expressions=50 --graphs=true --ast=false --use-log-fitness=false --workers=$WORKERS