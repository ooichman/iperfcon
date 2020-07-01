#!/bin/bash

#curl -s -X GET http://${SERVER_IP}:8080/iperf/status?server=${SERVER_IP},port=5001,type=log,critical=3000,warnging=7000,format=M
curl -s -X GET http://${CLIENT_IP}:8080/iperf/status?server=${SERVER_IP},port=5001,type=log,critical=3000,warnging=7000,format=M
