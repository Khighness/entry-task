#!/usr/bin/env bash

main="main.go"
tcpServer="tcp-server"
webServer="web-server"

echo "golang building tcp-server ..."
GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o tcp/${tcpServer} ../tcp/${main}
echo "build tcp-server successfully!"
echo "golang building web-server ..."
GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o web/${webServer} ../web/${main}
echo "build web-server successfully!"
