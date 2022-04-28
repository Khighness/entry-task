ARG BIN_DIR=build

FROM golang:1.17-alpine AS go-build
MAINTAINER KHighness "parakovo@gmail.com"
WORKDIR /${BIN_DIR}
COPY . .
RUN GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o web-server web/main.go

FROM alpine:latest AS run-prod
WORKDIR /var/entry
COPY --from=go-build /${BIN_DIR}/web-server .
COPY --from=go-build /${BIN_DIR}/web/conf ./web/conf
COPY --from=go-build /${BIN_DIR}/web/public ./web/public
EXPOSE 10000/tcp
ENTRYPOINT ["./web-server"]
