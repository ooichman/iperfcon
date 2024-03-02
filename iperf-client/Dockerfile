FROM golang:alpine as build
MAINTAINER Oren Oichman "Back to Root!!!"

WORKDIR /opt/app-root

COPY iperf-client.go iperf-client.go
RUN go build iperf-client.go

FROM docker.io/library/alpine:latest as iperf-builder
WORKDIR /app/src

RUN apk add --no-cache build-base && \
    wget https://github.com/esnet/iperf/releases/download/3.16/iperf-3.16.tar.gz -P /app/src && \
    tar -zxvf iperf-3.16.tar.gz
WORKDIR  /app/src/iperf-3.16/

RUN ./configure --prefix=/app/output && make && make install


#FROM scratch
FROM docker.io/library/alpine:latest
MAINTAINER Oren Oichman "Back to Root!!!"

ENV PATH="/usr/local/bin"
WORKDIR /opt/app-root
COPY --from=build /opt/app-root/iperf-client /opt/app-root/iperf-client
COPY --from=build /etc/passwd /etc/passwd

COPY --from=iperf-builder /app/output /usr/local


USER 1001
EXPOSE 8080

ENTRYPOINT ["./iperf-client"]
