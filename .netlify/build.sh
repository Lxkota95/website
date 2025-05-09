#!/bin/bash

set -x

ls -latr .
ls -latr ../
tree ../
cat netlify.toml
stat ../server/main.go

go build -o lxkota ../server/main.go