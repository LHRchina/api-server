#!/usr/bin/env bash

exec_file="bin/api-server"
if [ -e "$exec_file" ]
then
    rm $exec_file
else
    echo "ok"
fi

go build -o ./bin/api-server ./src/main.go
chmod u+x bin/api-server
./bin/api-server