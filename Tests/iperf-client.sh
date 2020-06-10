#!/bin/bash

curl -X GET http://localhost:8080/iperf/status?server=localhost,port=5001,type=json,critical=30000m,warnging=50000m
