FROM registry.redhat.io/ubi8/go-toolset as build
MAINTAINER Oren Oichman "Back to Root!!!"

WORKDIR /opt/app-root

COPY iperf-client.go iperf-client.go
RUN go build iperf-client.go

FROM registry.redhat.io/ubi8/ubi-minimal
MAINTAINER Oren Oichman "Back to Root!!!"

WORKDIR /opt/app-root
COPY --from=build /opt/app-root/iperf-client /opt/app-root/iperf-client

COPY iperf3-3.5-3.el8.x86_64.rpm /opt/app-root
RUN  rpm -Uvi iperf3-3.5-3.el8.x86_64.rpm && rm -f iperf3-3.5-3.el8.x86_64.rpm

USER 1001
EXPOSE 8080

ENTRYPOINT ["./iperf-client"]

