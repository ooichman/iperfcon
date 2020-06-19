#!/bin/bash

curl -s -X GET http://10.15.2.19:8080/iperf/status?server=10.15.2.19,port=5001,type=json,critical=3000,warnging=7000,format=M
#curl -X GET http://localhost:8080/iperf/status?server=localhost,port=5001,type=html,critical=3,warnging=5,format=g
