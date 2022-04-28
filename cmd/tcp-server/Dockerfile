FROM golang:1.17-buster as golang
COPY ./tcp /src/tcp
COPY ./pb /src/pb
COPY ./go.mod /src/
COPY ./go.sum /src/
COPY ./cmd/tcp-server /src/cmd/tcp-server
RUN cd /src/cmd/tcp-server && make linux

FROM debian:buster-20191014-slim as linux
WORKDIR /var/tcp-svc/
ADD ./application-tcp.yml .
COPY --from=golang /src/cmd/tcp-server/tcp-svc /bin/
EXPOSE 20000
CMD ["/bin/tcp-svc"]

