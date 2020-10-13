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
var CheckedHost string

func init() {
	flag.StringVar(
		&RBLAddress, "rbl", "zen.spamhaus.org",
		"Address of a single DNS RBL service to query",
	)
	flag.StringVar(
		&RBLFile, "rbl.file", "",
		"File with list of RBLs to check, one per line",
	)

	flag.StringVar(&CheckedHost, "host", "", "IP address or domain to check in RBL")
	flag.Parse()
}

func main() {
	if CheckedHost == "" || (RBLAddress == "" && RBLFile == "") {
		flag.Usage()
		os.Exit(5)
	}

	rbls, err := loadRBLs(RBLAddress, RBLFile)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(5)
	}

	var status string = "OK: " + CheckedHost + " at "
	for _, rbl := range rbls {
		result, err := rbl.CheckIP(CheckedHost)
		if err != nil {
			fmt.Printf("CRITICAL: failed to check %s at %s: %s\n", CheckedHost, rbl.Domain, err.Error())
			os.Exit(2)
		}
		if result.IsBlocked {
			fmt.Printf("CRITICAL: %s blocked at %s: %s\n", CheckedHost, rbl.Domain, result.BlockReason)
			os.Exit(2)
		} else {
			status += rbl.Domain + ", "
		}
	}
	fmt.Println(status)
	os.Exit(0)
}

func loadRBLs(rblAddress, rblFile string) ([]*RBL, error) {
	var rblDomains []string
	var err error
	if RBLFile != "" {
		rblDomains, err = loadRBLFromFile(RBLFile)
		if err != nil {
			return nil, fmt.Errorf("Error loading RBL domains from file '%s': %s", RBLFile, err.Error())
		}
	} else {
		rblDomains = []string{RBLAddress}
	}

	var rbls []*RBL
	for _, domain := range rblDomains {
		rbls = append(rbls, NewRBL(domain, ""))
	}

	return rbls, nil
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
