#!/usr/bin/env bash

main="main.go"
tcpServer="tcp-server"
webServer="web-server"

echo "golang building tcp-server ..."
GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o tcp/${tcpServer} ../tcp/${main}
echo "build tcp-server successfully!"
echo "golang building web-server ..."
GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o web/${webServer} ../web/${main}
echo "build web-server successfully!"

echo "docker building tcp-server image ..."
cd tcp && docker build -f Dockerfile -t entry-tcp-server:1.0 .
echo "build tcp-server image successfully!"
echo "docker building tcp-server image ..."
cd .. && cd web && docker build -f Dockerfile -t entry-web-server:1.0 .
echo "build web-server image successfully!"
