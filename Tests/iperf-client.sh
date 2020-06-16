#!/bin/bash

curl -X GET http://localhost:8080/iperf/status?server=localhost,port=5001,type=json,critical=3,warnging=5,format=g
#curl -X POST http://localhost:8080/iperf/status?server=localhost,port=5001,type=json,critical=30000,warnging=50000,format=M
