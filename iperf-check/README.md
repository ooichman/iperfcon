<img alt="Rook" src="media/iperf-logo.png" width="20%" height="20%">

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

    # oc create -f deployment.yaml

one it is deployed you can view the output on the logs view (or send it to ELK)

