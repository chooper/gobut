#!/bin/bash

test -z "$1" && { echo "Usage: $0 <version tag>"; exit 1; }
$tag = $1

docker build -t chooper/gobut:$tag .
