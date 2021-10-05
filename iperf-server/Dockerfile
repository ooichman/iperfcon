FROM registry.redhat.io/ubi8/ubi-minimal
MAINTAINER Oren Oichman "Back to ROOT !!!"

WORKDIR /opt/app-root
COPY run-server.sh /usr/sbin/
COPY iperf3-3.5-3.el8.x86_64.rpm /opt/app-root
RUN  rpm -Uvi iperf3-3.5-3.el8.x86_64.rpm && rm -f iperf3-3.5-3.el8.x86_64.rpm
USER 1001

CMD ["/usr/sbin/run-server.sh"]
EXPOSE 5001/tcp
