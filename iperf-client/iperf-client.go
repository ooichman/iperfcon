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
	"encoding/json"
)

type IperfOutput struct {
	End struct {
		SumSent struct {
//			Start         int     `json:"start"`
//			End           float64 `json:"end"`
//			Seconds       float64 `json:"seconds"`
//			Bytes         int64   `json:"bytes"`
			BitsPerSecond float64 `json:"bits_per_second"`
//			Retransmits   int     `json:"retransmits"`
		} `json:"sum_sent"`
		SumReceived struct {
//			Start         int     `json:"start"`
//			End           float64 `json:"end"`
//			Seconds       float64 `json:"seconds"`
//			Bytes         int64   `json:"bytes"`
			BitsPerSecond float64 `json:"bits_per_second"`
		} `json:"sum_received"`
	} `json:"end"`
}

type iperfValues struct {
		serverAddr string
		portNumber string
		outputType string
		iperfCritical string
		iperfWarning string
		iperfFormat string
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func checkFormat(ipc iperfValues) int64 {

	var duplicateNum int64
	duplicateNum = 1
	if ipc.iperfFormat == "k" || ipc.iperfFormat == "K" {
		duplicateNum = duplicateNum * 1024
	} else if ipc.iperfFormat == "m" || ipc.iperfFormat == "M" {
		duplicateNum = duplicateNum * 1024 * 1024
	} else if ipc.iperfFormat == "g" || ipc.iperfFormat == "G" {
		duplicateNum = duplicateNum * 1024 * 1024 * 1024
	} else if ipc.iperfFormat == "t" || ipc.iperfFormat == "T" {
		duplicateNum = duplicateNum * 1024 * 1024 * 1024 * 1024
	} else {
		duplicateNum = 1
	}

	return duplicateNum
}

func RunningIperf(ipc iperfValues,w http.ResponseWriter) {

	iperfCommand := "/usr/bin/iperf3"
//	log.Printf("the comman is %s\n", iperfCommand)

	cmd := exec.Command(iperfCommand,"-c",ipc.serverAddr,"-p",ipc.portNumber,"-J","-t","2","-i","0.5")
	out, cmderr := cmd.Output()

	if cmderr != nil {
		fmt.Fprintf(w,"there was an error running iperf with the given arguments: %+v\n",cmderr)
		http.Error(w,"", 502)
		return
	}

	var statusFlag string
	var iperfSumSent float64
	var iperfSumRec float64
	var critalNum float64
	var warningNum float64
	var ipout IperfOutput
	err := json.Unmarshal([]byte(out), &ipout)

	if err != nil {
		fmt.Fprintf(w,"there was an error parsing the output: %+v\n",err)
		http.Error(w,"", 502)
		return
	}
	duplicateSum := checkFormat(ipc)

	if duplicateSum != 1 {

	iperfSumSent = ipout.End.SumSent.BitsPerSecond / float64(duplicateSum)
	iperfSumRec = ipout.End.SumReceived.BitsPerSecond / float64(duplicateSum)
	
	} else {
		fmt.Fprintf(w,"Error: wrong Format type provided\n")
		http.Error(w,"", 504)
		return
	}
//	if checkFormat(ipc.iperfFormat)
	sumAvrg := (iperfSumSent + iperfSumRec) / float64(2)

	fmt.Sscan(ipc.iperfCritical, &critalNum)
	fmt.Sscan(ipc.iperfWarning, &warningNum)

//	critalNum = critalNum * float64(duplicateSum)
//	warningNum = warningNum * float64(duplicateSum)
	
	if sumAvrg <= critalNum {
		statusFlag = "Critical"
	} else if sumAvrg <= warningNum {
		statusFlag = "Warning"
	} else {
		statusFlag = "O.k"
	}
//	fmt.Fprintf(w,"critalNum is : %f , warningNum is %f", critalNum, warningNum)
	fmt.Fprintf(w,"{ \"Status\": \"%s\" , \"PacketSent\": \"%f\" , \"PacketReceived: \"%f\"\n", statusFlag, iperfSumSent , iperfSumRec )

}

func GetHandle(w http.ResponseWriter, r *http.Request) {

	var iperfClient iperfValues

	if r.Method != "GET" {
//		fmt.Fprintf(w, "Only GET Method is allowed\n")
		log.Printf("Received request Method %s not allowed, GET Only", r.Method)
		http.Error(w,"Only GET Method is allowed\n", 500)
		return
	}
	fullquery := r.URL.RawQuery
	fullquery = strings.ReplaceAll(fullquery,",","&")
//	fmt.Fprintf(w,"Your Query String is %s\n", fullquery)
  
    splitquery := strings.Split(fullquery, "&")
    
//    if len(splitquery) < 5 {
//		fmt.Fprintf(w, "Make sure your Query has all the values %s\n", fullquery)
//		log.Printf("Make sure your Query has all the values %s\n", fullquery)
//    }

    for i := 0 ; i < len(splitquery); i++ {
//    	fmt.Fprintf(w, "%s\n" , splitquery[i])
//    	fmt.Fprintf(os.Stdout, "%s\n" , splitquery[i])

		serverReg := regexp.MustCompile(`(server){1}`)
		if serverReg.MatchString(splitquery[i]) {
				serverSplit := strings.Split(splitquery[i], "=")
				iperfClient.serverAddr = serverSplit[1]
//				log.Printf("the Iperf Server IP Address is : %s", iperfClient.serverAddr)
		}

		portReg := regexp.MustCompile(`(port){1}`)
		if portReg.MatchString(splitquery[i]) {
				portSplit := strings.Split(splitquery[i], "=")
				iperfClient.portNumber = portSplit[1]
//				log.Printf("the Iperf Port Number is %s", iperfClient.portNumber)
		}

		typeReg := regexp.MustCompile(`(type){1}`)
		if typeReg.MatchString(splitquery[i]) {
				typeSplit := strings.Split(splitquery[i], "=")
				iperfClient.outputType = typeSplit[1]
//				log.Printf("the Iperf output type is %s", iperfClient.outputType)
		}
		warngingReg := regexp.MustCompile(`(warnging){1}`)
		if warngingReg.MatchString(splitquery[i]) {
				warngSplit := strings.Split(splitquery[i], "=")
				iperfClient.iperfWarning = warngSplit[1]
//				log.Printf("the warnging limit is %s", iperfClient.iperfWarning)
		}
		criticalReg := regexp.MustCompile(`(critical){1}`)
		if criticalReg.MatchString(splitquery[i]) {
				crticalSplit := strings.Split(splitquery[i], "=")
				iperfClient.iperfCritical = crticalSplit[1]
//				log.Printf("the critical limit is %s", iperfClient.iperfCritical)
		}
		formatReg := regexp.MustCompile(`(format){1}`)
		if formatReg.MatchString(splitquery[i]) {
				formatSplit := strings.Split(splitquery[i], "=")
				iperfClient.iperfFormat = formatSplit[1]
//				log.Printf("the Format is %s", iperfClient.iperfFormat)
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