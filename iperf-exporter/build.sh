#!/bin/bash

if [[ -z $1 ]]; then
	CONTAINER_TAG=iperf-exporter
	buildah bud -f Containerfile -t ${CONTAINER_TAG} .
else
	CONTAINER_TAG=$1
	buildah bud -f Containerfile -t ${CONTAINER_TAG} . 
	buildah push ${CONTAINER_TAG}