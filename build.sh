#!/bin/bash
CGO_ENABLED=0 go build -o ./build/v2ray_exporter
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/v2ray_exporter.linux-amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./build/v2ray_exporter.darwin-amd64