#!/bin/bash
cd ..
GOOS=linux GOARCH=arm go build -v -o tools/agora
cd tools
