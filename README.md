# iperfcon
iperfcon is for Openshift/Kubernetes bandwidth testing the SDN network.
the iperf-client is working with the iperf-server container and outputs 
the iperf3 results summery after we run a GET command  to the iperf-client 
route with the right values.

we can use it for several used cases:

- network bandwidth within the worker
- network bandwidth between 2 workers
- network bandwidth between a Worker and and external Server

## how to use it 

