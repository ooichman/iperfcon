#!/bin/bash
buildah bud -f Dockerfile -t quay.io/ooichman/iperf-client . && \
buildah push quay.io/ooichman/iperf-client
