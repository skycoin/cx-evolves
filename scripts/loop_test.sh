#!/bin/bash

if [ "$WORKERS" == "" ]
then
	WORKERS=24
fi

basePort=9090
i=0
while [ "$i" -lt $WORKERS ]; do
    echo "Base Port $(($basePort+$i)) "
    i=$(( i + 1 ))
done 