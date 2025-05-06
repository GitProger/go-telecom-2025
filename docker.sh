#!/bin/bash

docker run -v "$(pwd)"/"$1":/root/"$1" \
           -v "$(pwd)"/"$2":/root/"$2" \
           go-telecom-2025-docker "$1" "$2"