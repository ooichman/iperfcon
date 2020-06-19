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
	warning_bw  string
	critical_bw string
	outputFormat string
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

	fullurl := myvars.url_path
	fullurl += ",critical=" + myvars.critical_bw
	fullurl += ",warnging=" + myvars.warning_bw
	fullurl += ",format=" + myvars.outputFormat

	var x int
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

	urlinterval := getEnv("URL_INTERVAL", "300")
	url_client := getEnv("IPREF_CLIENT_URI", "nil")

	if url_client == "nil" {
		fmt.Println("Error - the IPREF_CLIENT_URI is not set")
		os.Exit(3)
	}

	iprefserver := getEnv("IPREF_SERVER","nil")
	
	if iprefserver == "nil" {
		fmt.Println("Error - the IPREF_SERVER is not set")
	}

	urlcritical := getEnv("CRITICAL_LIMIT", "3000")
	urlwarning := getEnv("WARNING_LIMIT", "5000")

	iperformat := getEnv("IPERF_FORMAT", "m")
	var usedebug bool
	_ , exists := os.LookupEnv("USE_DEBUG")
	//urlbool, exists := os.LookupEnv("USE_DEBUG")
    if !exists {
		usedebug = false
	//	urlbool = "false"
	} else {
		usedebug = true
	//	urlbool = "true"
	}

	//if urlbool == "true" {
	//	fmt.Println("the debug level is set to 1")
	//}
	
	urlpath := "http://" + url_client + "/iperf/status?server=" + iprefserver + ",port=5001,type=json"
	url := urlvars{url_int: urlinterval, url_path: urlpath, warning_bw: urlwarning, critical_bw: urlcritical, outputFormat: iperformat , url_debug: usedebug}

	RunCheck(url)
}
