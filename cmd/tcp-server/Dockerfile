FROM golang:1.17-buster as golang
COPY ./tcp /src/tcp
COPY ./pb /src/pb
COPY ./go.mod /src/
COPY ./go.sum /src/
COPY ./cmd/tcp-server /src/cmd/tcp-server
RUN cd /src/cmd/tcp-server && make linux

FROM centos:7.9.2009 as linux
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
WORKDIR /var/tcp-svc/
ADD ./application-tcp.yml .
COPY --from=golang /src/cmd/tcp-server/tcp-svc /bin/
EXPOSE 20000
CMD ["/bin/tcp-svc"]

