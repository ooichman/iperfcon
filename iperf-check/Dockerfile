FROM ubi8/go-toolset as build
MAINTAINER Oren Oichman "Back to Root!!!"

WORKDIR /opt/app-root
COPY cmd cmd
RUN go build cmd/iperf-check/iperf-check.go

FROM ubi8/ubi-minimal

WORKDIR /opt/app-root
COPY --from=build /opt/app-root/iperf-check /opt/app-root/iperf-check

ENTRYPOINT ["./iperf-check"]