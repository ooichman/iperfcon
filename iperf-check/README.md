<img alt="Rook" src="iperf-logo.png" width="20%" height="20%">

# iperf-check

the iperf-check is an simple web client written in GO which expects several environment variables
and then use those to query the ipref-client URL in the given interval.

### Environment Variables

  - URL_INTERVAL - set the interval time in second (most be an integer value) 
  - IPREF_CLIENT_URI - the route or service of the ipref client 
  - IPREF_SERVER - the service of the iperf server 
  - CRITICAL_LIMIT - ( Defualt : "30000m")
  - WARNING_LIMIT - ( Defualt : "50000m")
  - USE_DEBUG - Use this flag if you want to run the iperf-check in debug mode


### Deployment 

#### Without Operator

if you wish to deploy the tool without operator all you need to do is edit the deployment.yaml file 
with the right environment variables and the run the command :

    env:
    - name: URL_INTERVAL
      value: 30
    - name: IPREF_CLIENT_URI
      value: ''
    - name: IPREF_SERVER 
      value: ''
    - name: WARNING_LIMIT
      value: 50000m
    - name: CRITICAL_LIMIT
      value: 30000m
    - name: USE_DEBUG
      value: False

one you changed all the values you can go ahead and create the deployment

    # oc create -f pod-deployment.yaml

once it is deployed you can view the output on the logs view (or send it to ELK)

by running 

    # oc logs iperf-check

#### with Operator

Well , you don't need to do much , just create the CR in your desired namespace and the operator will\
deploy all the 