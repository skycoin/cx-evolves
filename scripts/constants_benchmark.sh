#!/bin/bash
  
# turn on bash's job control
set -m
  
# Start worker and put it in background
go run ./cxexecutes/worker/cmd/server.go -workers=3 &
  
# Start the main process
go run main.go --constants=true --rounds=10 --name=Constants --population=300 --generations=1000 --expressions=30 --graphs=true --ast=true --use-log-fitness=false --workers=3