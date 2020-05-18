<img alt="Rook" src="media/iperf-logo.png" width="20%" height="20%">

# iperfcon
iperfcon is for Openshift/Kubernetes bandwidth testing the SDN network.
the iperf-client is working with the iperf-server container and outputs 
the iperf3 results summery after we run a GET command  to the iperf-client 
route with the right values.

we can use it for several used cases:

- network bandwidth within the worker
- network bandwidth between 2 workers
- network bandwidth between a Worker and and external Server

## Deploying iperfcon 
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

The iperf-server container has 2 environment variables you can run in the deployment:

- IPERF_PROTOCOL - choose between tcp and udp (default: tcp)
- IPERF_PORT - choose the port on which the iperf server will listen upon (default: 5001)

## How to Use it
now run the curl command to the route to get the results:

    # curl -X GET \
    http://iperf-client-router/iperf/api.cgi?server=iperf-server-service,port=5001,type=json

The RESTAPI only expect a GET request with the following values :

- server - the iperf-server service IP address or name
- port - the port you are using on the iperf-server (the default is 5001)
- type - the type of output you want to see , that can be either HTML or JSON (lowercap latter ONLY!!)

if you want to look at the results in a nicer output you can pipe it to jq

    # curl -X GET  \
      http://iperf-client-router/iperf/api.cgi?server=iperf-server-service,port=5001,type=json | \
      jq


