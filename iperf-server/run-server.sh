#!/bin/bash
#script for running the iperf Server

if [ -z $IPERF_PROTOCOL ]; then
	IPERF_PROTOCOL="tcp";
fi

if [ -z $IPERF_PORT ]; then
	IPERF_PORT=5001
fi

case "${IPERF_PROTOCOL}" in

	'tcp')
		IPERFCMD="/usr/bin/iperf3 -s -p ${IPERF_PORT}"
	;;
	'udp')
		IPERFCMD="/usr/bin/iperf3 -s -u -p ${IPERF_PORT}"
esac
${IPERFCMD}
