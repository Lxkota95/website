#!/bin/bash

set -xe

ls -latr .
ls -latr ../
tree ../
cat ../netlify.toml
stat ../server/main.go

go install github.com/air-verse/air@latest
# go build -o lxkota ../server/main.go
air