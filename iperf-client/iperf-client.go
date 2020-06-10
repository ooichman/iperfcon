package main

import (
	"fmt"
//	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"net/http"
	"strings"
	"regexp"
)

type iperfValues struct {
		serverAddr string
		portNumber string
		outputType string
		iperfCritical string
		iperfWarning string
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func RunningIperf(ipc iperfValues,w http.ResponseWriter) {


	iperfCommand := "/usr/bin/iperf3"
	log.Printf("the comman is %s\n", iperfCommand)

	cmd := exec.Command(iperfCommand,"-c",ipc.serverAddr,"-p",ipc.portNumber,"-f","M","-t","2","-i","0.5")
	out, err := cmd.Output()

	if err != nil {
        fmt.Fprintf(w,"return Error code from command %s", err)
    }

//    fmt.Fprintf(w, "%s\n", out)
//	log.Printf("%s\n", out)


	fmt.Printf("the data is : %s", out)
}

func GetHandle(w http.ResponseWriter, r *http.Request) {

	var iperfClient iperfValues

	if r.Method != "GET" {
		fmt.Fprintf(w, "Only GET Method is allowed\n")
		log.Printf("Received request Method %s not allowed, GET Only", r.Method)
		os.Exit(2)
	}
	fullquery := r.URL.RawQuery
	fullquery = strings.ReplaceAll(fullquery,",","&")
//	fmt.Fprintf(w,"Your Query String is %s\n", fullquery)
  
    splitquery := strings.Split(fullquery, "&")
    
    if len(splitquery) < 5 {
		fmt.Fprintf(w, "Make sure your Query has all the values %s\n", fullquery)
		log.Printf("Make sure your Query has all the values %s\n", fullquery)
    }

    for i := 0 ; i < len(splitquery); i++ {
//    	fmt.Fprintf(w, "%s\n" , splitquery[i])
//    	fmt.Fprintf(os.Stdout, "%s\n" , splitquery[i])

		serverReg := regexp.MustCompile(`(server){1}`)
		if serverReg.MatchString(splitquery[i]) {
				serverSplit := strings.Split(splitquery[i], "=")
				iperfClient.serverAddr = serverSplit[1]
				log.Printf("the Iperf Server IP Address is : %s", iperfClient.serverAddr)
		}

		portReg := regexp.MustCompile(`(port){1}`)
		if portReg.MatchString(splitquery[i]) {
				portSplit := strings.Split(splitquery[i], "=")
				iperfClient.portNumber = portSplit[1]
				log.Printf("the Iperf Port Number is %s", iperfClient.portNumber)
		}

		typeReg := regexp.MustCompile(`(type){1}`)
		if typeReg.MatchString(splitquery[i]) {
				typeSplit := strings.Split(splitquery[i], "=")
				iperfClient.outputType = typeSplit[1]
				log.Printf("the Iperf output type is %s", iperfClient.outputType)
		}
		warngingReg := regexp.MustCompile(`(warnging){1}`)
		if warngingReg.MatchString(splitquery[i]) {
				warngSplit := strings.Split(splitquery[i], "=")
				iperfClient.iperfWarning = warngSplit[1]
				log.Printf("the warnging limit is %s", iperfClient.iperfWarning)
		}
		criticalReg := regexp.MustCompile(`(critical){1}`)
		if criticalReg.MatchString(splitquery[i]) {
				crticalSplit := strings.Split(splitquery[i], "=")
				iperfClient.iperfCritical = crticalSplit[1]
				log.Printf("the critical limit is %s", iperfClient.iperfCritical)
		}
    }

    RunningIperf(iperfClient,w)
    
}

func main() {
	
    servicePort := getEnv("IPERF_CLIENT_PORT", "8080")
    servicePort = ":"+servicePort
	http.HandleFunc("/iperf/api.cgi", GetHandle)
	http.HandleFunc("/iperf/status", GetHandle)

	if err := http.ListenAndServe(servicePort, nil) ; err != nil {
    	log.Fatalf("Could not listen on port %s %v", servicePort , err)
    }
}