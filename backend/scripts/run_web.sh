#!/bin/bash

# wait for MYSQL server to start
./wait-for-it.sh db:3306 -t 30
./db_migrate
./main
