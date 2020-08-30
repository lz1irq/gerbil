package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s i.p.to.check dns.rbl.domain\n", os.Args[0])
		os.Exit(5)
	}
	rawIP := os.Args[1]
	fmt.Printf("Original IP: %s\n", rawIP)
	dnsRbl := os.Args[2]
	fmt.Printf("DNS RBL domain: %s\n", dnsRbl)

	rbl := NewRBL(dnsRbl, ".")

	result, err := rbl.CheckIP(rawIP)
	if err != nil {
		fmt.Printf("Failed to check IP %s: %s\n", rawIP, err.Error())
		os.Exit(5)
	}
	fmt.Println(result)

}
