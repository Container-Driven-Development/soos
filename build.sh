#!/bin/bash

rm -rf out

GOOS=linux GOARCH=amd64 go build -o out/soos-Linux-x86_64 -ldflags='-s -w' soos.go && upx --brute out/soos-Linux-x86_64
GOOS=darwin GOARCH=amd64 go build -o out/soos-Darwin-x86_64 -ldflags='-s -w' soos.go && upx --brute out/soos-Darwin-x86_64
GOOS=windows GOARCH=amd64 go build -o out/soos-Windows-x86_64.exe -ldflags='-s -w' soos.go && upx --brute out/soos-Windows-x86_64.exe
