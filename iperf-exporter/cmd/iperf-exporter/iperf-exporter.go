//iperf-exporter.go
package main

import (
		"fmt"
		"net/http"
		"os"
		"log"
		"text/template"
		"encoding/json"
		"io/ioutil"


)

type IperfOutput struct {
	IperfResult struct {
		Date           string  `json:"date"`
		Status         string  `json:"Status"`
		Rate           string  `json:"Rate"`
		PacketSent     float64 `json:"PacketSent"`
		PacketReceived float64 `json:"PacketReceived"`
	} `json:"iperfResult"`
}

type TempResult struct {
	Status         int
	Warnging       int
	Critical       int
	PacketSent	   float64
	PacketReceived float64
}

type urlvars struct {
	url_path      string
	server_port   string
	export_port   string
	server_path   string
	warning_bw    string
	critical_bw   string
	outputFormat  string
	url_debug	  bool
}

const doc = `
<!DOCTYPE html>
<html>
<head lang="en">
	<meta charset="UTF-8">
	<title> Prometheus Iperf Exporter </title>
</head>
<body>
<p style="margin-bottom:0; padding-top:0;">iperf_exporter_status {{.Status}}</p>
<p style="margin-bottom:0; padding-top:0;">iperf_exporter_warnging {{.Warnging}}</p>
<p style="margin-bottom:0; padding-top:0;">iperf_exporter_critical {{.Critical}}</p>
<p style="margin-bottom:0; padding-top:0;">iperf_exporter_packet_sent {{.PacketSent}}</p>
<p style="margin-bottom:0; padding-top:0;">iperf_exporter_packet_received {{.PacketReceived}}</p>
</body>
</html>
`


func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}


func (url *urlvars) MetHandle(w http.ResponseWriter, r *http.Request) {

	var warning_bw int
	var critical_bw int
	var cmdOputput IperfOutput


	fullpath := "http://" + url.url_path + "/iperf/status?server="
	fullpath += url.server_path + ",port=" + url.server_port 
	fullpath += ",type=json" + ",format=" + url.outputFormat

	w.Header().Add("Content Type", "text/html")
	tmpl ,temperr := template.New("PrometheusTemplate").Parse(doc)

	resp, err := http.Get(fullpath)
	if err != nil {
    	log.Fatal(err)
    }
    defer resp.Body.Close()
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    jsonerr := json.Unmarshal(bodyBytes , &cmdOputput)

    if jsonerr != nil {
    	log.Fatal(jsonerr)
    }

	fmt.Sscan(url.warning_bw , &warning_bw)
	fmt.Sscan(url.critical_bw , &critical_bw)

	if temperr == nil {
		temp_res := TempResult{0, 
		warning_bw , 
		critical_bw , 
		cmdOputput.IperfResult.PacketSent , 
		cmdOputput.IperfResult.PacketReceived }
		tmpl.Execute(w, temp_res)
	}

}

func main() {
	
	var url urlvars


	url.url_path = getEnv("IPREF_CLIENT_URI", "nil")

	if url.url_path == "nil" {
		fmt.Println("Error - the IPREF_CLIENT_URI is not set")
		os.Exit(3)
	}

	url.server_path = getEnv("IPREF_SERVER","nil")
	
	if url.server_path == "nil" {
		fmt.Println("Error - the IPREF_SERVER is not set")
	}

	url.server_port = getEnv("SERVER_PORT", "5001")
	url.critical_bw = getEnv("CRITICAL_LIMIT", "3000")
	url.warning_bw = getEnv("WARNING_LIMIT", "5000")

	url.outputFormat = getEnv("IPERF_FORMAT", "m")
	
	_ , exists := os.LookupEnv("USE_DEBUG")
	//urlbool, exists := os.LookupEnv("USE_DEBUG")
    if !exists {
    	url.url_debug = false
//	urlbool = "false"
	} else {
		url.url_debug = true
	//	urlbool = "true"
	}

	url.export_port = getEnv("IPERF_EXPORTER_PORT", "9100")
    url.export_port = ":"+url.export_port

    http.HandleFunc("/metrics", url.MetHandle)

    if err := http.ListenAndServe(url.export_port, nil) ; err != nil {
    	log.Fatalf("Could not listen on port %s %v", url.export_port , err)
    }

}