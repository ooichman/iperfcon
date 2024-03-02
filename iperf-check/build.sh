#!/bin/bash

	if [[ -z ${1} ]]; then
		buildah bud -f Dockerfile -t iperf-client .
    else
    	CONTAINER_TAG=${1}
    	buildah bud -f Dockerfile -t ${CONTAINER_TAG}
 		buildah push ${CONTAINER_TAG}
 	fi
