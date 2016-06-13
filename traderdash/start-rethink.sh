#!/bin/bash

echo Starting a 2 node rethink cluster

rethinkdb &

sleep 5

rethinkdb --port-offset 1 --directory rethinkdb_data2 --join localhost:29015 &
