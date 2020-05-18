package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type urlvars struct {
	url_int     string
	url_path    string
	warning_bw  string
	critical_bw string
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
	fullurl += ",warning=" + myvars.warning_bw
	fullurl += ",critical=" + myvars.critical_bw

	var x int
	if _, err := fmt.Sscan(myvars.url_int, &x); err == nil {
		fmt.Printf("the inteval is %d seconds\n", x)
	} else {
     	os.Exit(2)
	}

    for {
		resp, err := http.Get(fullurl)
		if err == nil {
			fmt.Printf("GET for %s was successful", fullurl)
		} else {
			fmt.Printf("Unable to reach %s", fullurl)
		}
		time.Sleep(time.Second * time.Duration(x)) 
		defer resp.Body.Close()
    }
}

func main() {

	// checking and testing the environment variables

	urlinterval := getEnv("URL_INTERVAL", "300s")
	urlpath := getEnv("URL_PATH", "nil")

	if urlpath == "nil" {
		fmt.Println("Error - the URL_PATH is not set")
		os.Exit(3)
	}

	urlcritical := getEnv("CRITICAL_LIMIT", "30000m")
	urlwarning := getEnv("WARNING_LIMIT", "50000m")

	url := urlvars{url_int: urlinterval, url_path: urlpath, warning_bw: urlwarning, critical_bw: urlcritical}

	RunCheck(url)
}
