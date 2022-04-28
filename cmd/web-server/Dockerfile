FROM golang:1.17-buster as golang
COPY ./web /src/web
COPY ./pb /src/pb
COPY ./go.mod /src/
COPY ./go.sum /src/
COPY ./cmd/web-server /src/cmd/web-server
RUN cd /src/cmd/web-server && make linux

FROM debian:buster-20191014-slim as linux
WORKDIR /var/web-svc/
ADD ./application-web.yml .
COPY --from=golang /src/cmd/web-server/web-svc /bin/
EXPOSE 10000
CMD ["/bin/web-svc"]

