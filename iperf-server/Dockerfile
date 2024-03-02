FROM docker.io/library/alpine:latest as iperf-builder
MAINTAINER Oren Oichman "Back to Root!!!"
WORKDIR /app/src

RUN apk add --no-cache build-base && \
    wget https://github.com/esnet/iperf/releases/download/3.16/iperf-3.16.tar.gz -P /app/src && \
    tar -zxvf iperf-3.16.tar.gz
WORKDIR  /app/src/iperf-3.16/

RUN ./configure --prefix=/app/output && make && make install


#FROM scratch
FROM docker.io/library/alpine:latest

ENV PATH="/usr/local/bin"
WORKDIR /opt/app-root

COPY --from=iperf-builder /app/output /usr/local


USER 1001
EXPOSE 5001

CMD ["/usr/sbin/run-server.sh"]
