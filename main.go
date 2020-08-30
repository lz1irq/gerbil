package main

import (
	"flag"
	"fmt"
	"os"
)

var RBLAddress string
var CheckedIP string

func init() {
	flag.StringVar(
		&RBLAddress, "rbl", "zen.spamhaus.org",
		"Address of a single DNS RBL service to query",
	)

	flag.StringVar(&CheckedIP, "ip", "", "IP address to check in RBL")
	flag.Parse()
}

func main() {
	if CheckedIP == "" {
		flag.Usage()
		os.Exit(5)
	}

	rbl := NewRBL(RBLAddress, ".")

	result, err := rbl.CheckIP(CheckedIP)
	if err != nil {
		fmt.Printf("Failed to check IP %s: %s\n", CheckedIP, err.Error())
		os.Exit(5)
	}
	fmt.Println(result)

}
