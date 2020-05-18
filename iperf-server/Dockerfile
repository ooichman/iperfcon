FROM ubi8/ubi-minimal
MAINTAINER Oren Oichman "Back to ROOT !!!"

COPY run-server.sh /usr/sbin/
COPY iperf3 /usr/bin/
COPY libiperf.so.0 libiperf.so.0.0.0 /usr/lib64/
USER 1001

CMD ["/usr/sbin/run-server.sh"]
EXPOSE 5001/tcp
