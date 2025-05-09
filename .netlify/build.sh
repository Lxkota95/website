#!/bin/bash

set -xe

# Debug information
echo "Current directory:"
pwd
ls -latr .
echo "Parent directory:"
ls -latr ../

# Create necessary directories
mkdir -p ../public/static
mkdir -p ../functions/server

# Copy static files to public directory
cp -r ../static/* ../public/static/

# Install AWS Lambda dependency if not already in go.mod
cd ..
go get github.com/aws/aws-lambda-go

# Replace the import path in the serverless adapter with your actual module path
MODULE_NAME=$(go list -m)
echo "Module name: $MODULE_NAME"

# Build the serverless function
go build -o functions/server/server functions/server/main.go

# Copy required files for the function
mkdir -p functions/server/content
cp -r content/* functions/server/content/
cp template.html functions/server/

# For Netlify debugging
echo "Build completed"
ls -la public/
ls -la functions/