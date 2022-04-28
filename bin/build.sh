#!/bin/sh

# tcp
echo "$(TZ=UTC-8 date +%Y-%m-%d" "%H:%M:%S) building image tcp-svc ..."
docker build -t entry/tcp-svc -f cmd/tcp-server/Dockerfile .

# web
echo "$(TZ=UTC-8 date +%Y-%m-%d" "%H:%M:%S) building image web-svc ..."
docker build -t entry/web-svc -f cmd/web-server/Dockerfile .

# vue
echo "$(TZ=UTC-8 date +%Y-%m-%d" "%H:%M:%S) building image front-svc  ..."
docker build -t entry/front-svc -f cmd/front-end/Dockerfile .
