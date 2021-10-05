#!/bin/bash
	buildah bud -f Dockerfile -t quay.io/ooichman/iperf-server .
	buildah push quay.io/ooichman/iperf-server
