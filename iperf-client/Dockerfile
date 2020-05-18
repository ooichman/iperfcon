FROM registry.redhat.io/ubi8
MAINTAINER Oren Oichman

COPY iperf3-3.5-3.el8.x86_64.rpm /
RUN  yum update -y && yum install -y iperf3-3.5-3.el8.x86_64.rpm python3 && yum clean all && rm -f iperf3-3.5-3.el8.x86_64.rpm
USER 1001
COPY iperfcon.py /home/iperf/

EXPOSE 8080

CMD ["python3", "/home/iperf/iperfcon.py"]


