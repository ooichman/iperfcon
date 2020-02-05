<img alt="Rook" src="media/iperf-logo.png" width="50%" height="50%">

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
