package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var RBLAddress string
var RBLFile string
var CheckedIP string

func init() {
	flag.StringVar(
		&RBLAddress, "rbl", "zen.spamhaus.org",
		"Address of a single DNS RBL service to query",
	)
	flag.StringVar(
		&RBLFile, "rbl.file", "",
		"File with list of RBLs to check, one per line",
	)

	flag.StringVar(&CheckedIP, "ip", "", "IP address to check in RBL")
	flag.Parse()
}

func main() {
	if CheckedIP == "" || (RBLAddress == "" && RBLFile == "") {
		flag.Usage()
		os.Exit(5)
	}

	var rblDomains []string
	var err error
	if RBLFile != "" {
		rblDomains, err = loadRBLFromFile(RBLFile)
		if err != nil {
			fmt.Printf("Error loading RBL domains from file '%s': %s\n", RBLFile, err.Error())
		}
	} else {
		rblDomains = []string{RBLAddress}
	}

	var rbls []*RBL
	for _, domain := range rblDomains {
		rbls = append(rbls, NewRBL(domain, ""))
	}

	for _, rbl := range rbls {
		result, err := rbl.CheckIP(CheckedIP)
		if err != nil {
			fmt.Printf("Failed to check IP %s: %s\n", CheckedIP, err.Error())
			os.Exit(5)
		}
		fmt.Printf("rbl=%s, ip=%s, blocked=%t, reason=%s\n", rbl.Domain, CheckedIP, result.IsBlocked, result.BlockReason)
	}

}

func loadRBLFromFile(fname string) ([]string, error) {
	var domains []string

	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return domains, nil
}
