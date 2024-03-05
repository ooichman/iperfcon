package main

import (
    "fmt"
    "os"
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "encoding/json"
	"io/ioutil"
	"time"
	"log"
)

 var (
 	iperfWarning = promauto.NewGauge(prometheus.GaugeOpts {
        Name: "iperf_warning_bandwidth",
        Help: "the custom set warning bandwidth",
    })

    iperfCritical = promauto.NewGauge(prometheus.GaugeOpts {
        Name: "iperf_critical_bandwidth",
        Help: "The custom set Critical bandwidth",
     })
     
     iperfStatusValue = promauto.NewGauge(prometheus.GaugeOpts {
        Name: "iperf_status",
        Help: "The Iperf Status regarding the bandwidth usage",
    })

    iperfRate = promauto.NewGauge(prometheus.GaugeOpts {
    	Name: "iperf_rate",
    	Help: "The Iperf Rate by the last reading",
    })

    iperfPacketSent = promauto.NewGauge(prometheus.GaugeOpts {
        Name: "iperf_packet_sent",
        Help: "the Iperf Number of Packet Sent",
    })

    iperfPacketReceived = promauto.NewGauge(prometheus.GaugeOpts {
        Name: "iperf_packet_received",
        Help: "The Iperf Number of Packet Received",
    })
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

type urlvars struct {
	url_path      string
	server_port   string
	export_port   string
	server_path   string
	warning_bw    string
	critical_bw   string
	outputFormat  string
	url_debug	  bool
	interval      string
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func  MetHandle(url urlvars) {

    var warning_bw int
	var critical_bw int
	var intreval_check int
	var cmdOputput IperfOutput
	var status int

	fullpath := "http://" + url.url_path + "/iperf/status?server="
	fullpath += url.server_path + ",port=" + url.server_port 
	fullpath += ",type=json" + ",format=" + url.outputFormat

    for {


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
           fmt.Sscan(url.interval, &intreval_check)

		   if cmdOputput.IperfResult.Status == "O.k" {
			   status = 0
		   } else if cmdOputput.IperfResult.Status == "Warning" {
			   status = 1
		   } else if cmdOputput.IperfResult.Status == "Critical" {
			   status = 2
		   }

		   iperfStatusValue.Set(float64(status))
   		   iperfWarning.Set(float64(warning_bw))
 	       iperfCritical.Set(float64(critical_bw))
    	   iperfPacketSent.Set(float64(cmdOputput.IperfResult.PacketSent))
	       iperfPacketReceived.Set(float64(cmdOputput.IperfResult.PacketReceived))

	       time.Sleep(time.Duration(intreval_check) * time.Second)
	}
}


func main() {

      url := urlvars{}
      
      url.url_path = getEnv("IPREF_CLIENT_URI", "nil")

	if url.url_path == "nil" {
		fmt.Println("Error - the IPREF_CLIENT_URI is not set")
		os.Exit(3)
	}

    url.interval = getEnv("INTERVAL_CHECK","10")

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

    http.Handle("/metrics", promhttp.Handler())
    go func() {
       log.Printf("Starting HTTP Service on port %v", url.export_port)
       log.Fatal(http.ListenAndServe(":"+url.export_port, nil))
    }()
    MetHandle(url)

}