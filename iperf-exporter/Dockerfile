FROM ubi8/go-toolset as build
MAINTAINER Oren Oichman "Back to Root!!!"

WORKDIR /opt/app-root
COPY cmd cmd
RUN go build cmd/iperf-exporter/iperf-exporter.go

FROM ubi8/ubi-minimal

USER 1001
WORKDIR /opt/app-root
COPY --from=build /opt/app-root/iperf-exporter /opt/app-root/iperf-exporter

EXPOSE 9100
ENTRYPOINT ["./iperf-exporter"]
