#!/bin/bash

python ./display.py &
dpid=$!
sleep 1
python ./server.py &
spid=$!
sleep 1
python ./gui.py
wait $spid
wait $dpid
