#!/bin/bash
CGO_ENABLED=0
go build -o ./build/v2ray-exporter.linux-amd64
GOOS=linux GOARCH=amd64 go build -o ./build/v2ray-exporter.linux-amd64
GOOS=darwin GOARCH=amd64 go build -o ./build/v2ray-exporter.darwin-amd64