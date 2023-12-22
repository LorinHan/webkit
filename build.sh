#!/bin/bash

statik -src=./template
echo "Successfully generated statik file ..."
go build -o cmd/webkit ./gen/gen.go
echo "OK"