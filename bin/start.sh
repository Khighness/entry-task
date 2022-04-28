#!/bin/sh

docker run --name et-tcp-svc -d -p 20000:20000 entry/tcp-svc
docker run --name et-web-svc -d -p 10000:10000 entry/web-svc
docker run --name et-fe-svc -d -p 80:80 entry/fe-svc
