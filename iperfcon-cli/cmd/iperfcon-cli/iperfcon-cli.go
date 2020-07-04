package main

import (
		"fmt"
		"flag"
		"net/http"
		"log"
		"io"
		"os"
)

func main() {
	
	warngingPtr := flag.String("w", "5000" , "The Warnging Value")
	criticalPtr := flag.String("c", "3000" , "The Critical Value")
	urlPtr := flag.String("u" , "nil" , "The URL client (with PORT) ")
	serverPtr := flag.String("s" ,"nil" ,"The Server FQDN/IP")
	portPtr := flag.String("p", "5001" ,"the Iperf Server Port")
	formatPtr := flag.String("f" , "M", "The Format output")
	outputPtr := flag.String("o" , "log", "the Output format {html|log|json}")
	HelpPtr := flag.Bool("h" , false , "for the Help Menu")
	debugPtr := flag.Bool("d" , false , "for DEBUG")
	

	flag.Parse()
	if *HelpPtr == true {
		fmt.Printf("iperfcon-cli:\n")
		fmt.Printf("\t -w - for setting the Warnging, default Value: 5000\n")
		fmt.Printf("\t -c - for setting the Critical, default Value: 3000\n")
		fmt.Printf("\t -u - The URL of the iperf Client, default Value: nil\n")
		fmt.Printf("\t -s - The iperf Server FQDN/IP, default Value: nil\n")
		fmt.Printf("\t -p - The Port for the iperf Server, default Value: 5001\n")
		fmt.Printf("\t -f - The Output Format Option: {KMGTkmgt }, default Value: M\n")
		fmt.Printf("\t -o - The Output type {log|html|json} , default Value: log\n")
		fmt.Printf("\n")
		return
	}

	if *debugPtr == true {
		fmt.Printf("\t -w - Corrent Value: %s\n", *warngingPtr)
		fmt.Printf("\t -c - Corrent Value: %s\n", *criticalPtr)
		fmt.Printf("\t -u - Corrent Value: %s\n", *urlPtr)
		fmt.Printf("\t -s - Corrent Value: %s\n", *serverPtr)
		fmt.Printf("\t -p - Corrent Value: %s\n", *portPtr)
		fmt.Printf("\t -f - Corrent Value: %s\n", *formatPtr)
		fmt.Printf("\t -o - Corrent Value: %s\n", *outputPtr)
		fmt.Printf("\n")
	}


	fullpath := "http://" + *urlPtr + "/iperf/status?server=" + *serverPtr +",format=" + *formatPtr
	fullpath += ",port=" + *portPtr + ",warnging=" + *warngingPtr + ",critical=" + *criticalPtr
	fullpath += ",type=" + *outputPtr
//	fullpath += ",critical=" + *criticalPtr

	if *debugPtr == true {
		fmt.Println("the fullpath string is: ", fullpath)
	}

	resp, err := http.Get(fullpath)
	if err != nil {
    	log.Fatal(err)
    }
    defer resp.Body.Close()

    _, oserr := io.Copy(os.Stdout, resp.Body)
            	if oserr != nil {
                	log.Fatal(oserr)
            	}
}