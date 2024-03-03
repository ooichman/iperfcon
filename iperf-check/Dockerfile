FROM golang:alpine as build
MAINTAINER Oren Oichman "Back to Root!!!"

WORKDIR /opt/app-root
COPY cmd cmd
RUN go build cmd/iperf-check/iperf-check.go

FROM scratch

WORKDIR /opt/app-root
COPY --from=build /opt/app-root/iperf-check /opt/app-root/iperf-check
COPY --from=build /etc/passwd /etc/passwd

ENTRYPOINT ["/opt/app-root/iperf-check"]
