#!/bin/sh

cd ..

# tcp
echo "$(TZ=UTC-8 date +%Y-%m-%d" "%H:%M:%S) building image tcp-server ..."
docker build -t khighness/entry-tcp-server:v1 -f ./tcp/Dockerfile .
echo "$(TZ=UTC-8 date +%Y-%m-%d" "%H:%M:%S) build tcp-server successfully!"

# web
echo "$(TZ=UTC-8 date +%Y-%m-%d" "%H:%M:%S) building image web-server ..."
docker build -t khighness/entry-web-server:v1 -f ./web/Dockerfile .
echo "$(TZ=UTC-8 date +%Y-%m-%d" "%H:%M:%S) build web-server successfully!"

