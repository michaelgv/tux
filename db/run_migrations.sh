#!/bin/bash

DIRECTION=$1

echo "Running migrations going $1"

for file in `ls ./migrations/ |grep $1| sort`; do
    cat ./migrations/$file | ./connectdb.sh
done