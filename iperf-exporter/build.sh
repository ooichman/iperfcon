#!/bin/bash

if [[ -z $1 ]]; then
	CONTAINER_TAG=iperf-exporter
else
	CONTAINER_TAG=$1
fi
	
   buildah bud -f Containerfile -t ${CONTAINER_TAG} . 

if [[ ! -z $1 ]]; then
   buildah push ${CONTAINER_TAG}
fi
