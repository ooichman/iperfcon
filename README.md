<img alt="Rook" src="media/iperf-logo.png" width="25%" height="25%">

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
the usage of the containers is very simple.

first let's build the namespace for them :

    # oc new-project iperf

clone the github to your current working directory

    # git clone https://github.com/ooichman/iperfcon.git

now let's run the deployment for both Deployments

    # oc create -f iperfcon/iperf-server/pod-deployment.yaml
    # oc create -f iperfcon/iperf-client/pod-deployment.yaml

now make sure the pods are deployed as you expected them to be :

    # oc get pods -n iperf -o wide

if you want to check the communication between 2 workers make sure the pods are spraed out

- create a service for the iperf-server with port 5001 and a service for the iperf-client with port 8080
- create a route for the iperf-client

now run the curl command to the route to get the results:

   #curl -X GET  http://iperf-client-router/iperf/api.cgi?nan=wwww,server=iperf-server-service,port=5001,type=json

if you want to look at the results in a nicer output you can pipe it to jq

   #curl -X GET  http://iperf-client-router/iperf/api.cgi?nan=wwww,server=iperf-server-service,port=5001,type=json | jq


