package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"log"
	"io"
)

type urlvars struct {
	url_int     string
	url_path    string
	server_port string
	server_path string
	warning_bw  string
	critical_bw string
	outputFormat string
	outputType string
	url_debug	bool
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func RunCheck(myvars urlvars) {

	var x int
	fullurl := myvars.url_path
	fullurl += ",critical=" + myvars.critical_bw
	fullurl += ",warnging=" + myvars.warning_bw
	fullurl += ",format=" + myvars.outputFormat

	if myvars.url_debug == true {
		fmt.Printf("the url is %s",fullurl)
	}
	
		if _, err := fmt.Sscan(myvars.url_int, &x); err == nil {
			if myvars.url_debug == true {
				fmt.Printf("the inteval is %d seconds\n", x)
			}
		} else {
     		os.Exit(2)
		}

    	for {
			resp, err := http.Get(fullurl)
			if err == nil {
				if myvars.url_debug == true {
					fmt.Printf("GET for %s was successful\n", fullurl)
				}
				defer resp.Body.Close()
		    	_, err = io.Copy(os.Stdout, resp.Body)
            	if err != nil {
                	log.Fatal(err)
            	}
			} else {
				fmt.Printf("Unable to reach %s\n", fullurl)
				os.Exit(3)
			}
		
			time.Sleep(time.Second * time.Duration(x)) 
		
    	}
	
}

func main() {

	// checking and testing the environment variables
	var url urlvars

	url.url_int = getEnv("URL_INTERVAL", "300")
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

	url.outputType = getEnv("OUTPUT_TYPE", "log" )
	//if urlbool == "true" {
	//	fmt.Println("the debug level is set to 1")
	//}
	
	url.url_path = "http://" + url.url_path + "/iperf/status?server="
	url.url_path += url.server_path + ",port=" + url.server_port 
	url.url_path += ",type=" + url.outputType



	RunCheck(url)
}
