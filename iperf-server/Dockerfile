FROM centos
MAINTAINER Oren Oichman

COPY run-server.sh /usr/sbin/
RUN  useradd iperf && yum update -y && yum install -y iperf3 && yum clean all
USER iperf
CMD ["/usr/sbin/run-server.sh"]
EXPOSE 5001/tcp
