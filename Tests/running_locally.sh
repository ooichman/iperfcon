podman ps | grep iperf-server 2>&1 > /dev/null && RETVAL=0 || RETVAL=1

if [[ ${RETVAL} -eq 1 ]]; then
podman run -d -p 5001:5001 --rm --name iperf-server quay.io/ooichman/iperf-server
sleep 2
fi

podman ps | grep iperf-client 2>&1 > /dev/null && RETVAL=0 || RETVAL=1

if [[ ${RETVAL} -eq 1 ]]; then
podman run -d -p 8080:8080 --rm --name iperf-client quay.io/ooichman/iperf-client
sleep 2
fi

IF_NAME=`netstat -rnv | egrep -v "Gateway|Kernel" | head -1 | awk '{print $8}'`
IF_ADDR=`ip addr show dev ${IF_NAME} | grep "inet " | awk '{print $2}' | cut -d \/ -f 1`

podman ps | grep iperf-exporter 2>&1 > /dev/null && RETVAL=0 || RETVAL=1

if [[ ${RETVAL} -eq 1 ]]; then
podman run -d -p 9100:9100 \
	-e IPREF_CLIENT_URI=${IF_ADDR}:8080 \
	-e IPREF_SERVER=${IF_ADDR} \
	--name iperf-exporter \
	quay.io/ooichman/iperf-exporter
sleep 2
fi

podman ps | grep iperf-check 2>&1 > /dev/null && RETVAL=0 || RETVAL=1

if [[ ${RETVAL} -eq 1 ]]; then
podman run -d --rm \
	--name iperf-check \
	-e IPREF_CLIENT_URI=${IF_ADDR}:8080 \
	-e IPREF_SERVER=${IF_ADDR} \
	quay.io/ooichman/iperf-check
fi