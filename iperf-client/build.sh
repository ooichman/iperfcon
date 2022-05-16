#!/bin/bash

	if [[ -z ${1} ]]; then
		buildah bud -f Containerfile -t iperf-client .
    else
    	CONTAINER_TAG=${1}
    	buildah bud -f Containerfile -t ${CONTAINER_TAG}
 		buildah push ${CONTAINER_TAG}
 	fi
