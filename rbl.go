package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

// RBLResult represents the result of an DNS RBL check, including the
// checked IP, whether it is blocked or not, and optionally a reason for the
// block
type RBLResult struct {
	IP          string
	IsBlocked   bool
	BlockReason string
}

func (r *RBLResult) String() string {
	if r.IsBlocked {
		return fmt.Sprintf("IP %s blocked: %s", r.IP, r.BlockReason)
	}
	return fmt.Sprintf("IP %s not blocked", r.IP)
}

const defaultSeparator = "."
const defaultTimeout = 1 * time.Second

// RBL abstracts a DNS RBL service.
type RBL struct {
	Domain    string
	Separator string
	Timeout   time.Duration
	resolver  *net.Resolver
}

// NewRBL initializes a new RBL struct.
func NewRBL(domain, separator string) *RBL {
	if separator == "" {
		separator = defaultSeparator
	}
	return &RBL{
		Domain:    domain,
		Separator: separator,
		Timeout:   defaultTimeout,
		resolver:  &net.Resolver{PreferGo: false},
	}
}

// CheckIP checks whether is considered blocked in this particular RBL.
func (r *RBL) CheckIP(ip string) (*RBLResult, error) {
	queryDomain := r.formatQuery(ip)

	isBlocked, err := r.isBlocked(queryDomain)
	if err != nil {
		return nil, err
	}

	result := &RBLResult{
		IP:        ip,
		IsBlocked: isBlocked,
	}

	if result.IsBlocked {
		reason, err := r.getBlockReason(queryDomain)
		if err != nil {
			return result, nil
		}
		if reason == "" {
			reason = "unspecified reason"
		}
		result.BlockReason = reason
	}

	return result, nil
}

func (r *RBL) isBlocked(query string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	_, err := r.resolver.LookupIPAddr(ctx, query)
	if err != nil {
		dnsErr, ok := err.(*net.DNSError)
		if !ok {
			fmt.Printf("err : %s\n", err.Error())
			return false, err
		}

		if dnsErr.IsNotFound {
			return false, nil
		}
	}
	return true, nil
}

func (r *RBL) getBlockReason(query string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	defer cancel()

	reasons, err := r.resolver.LookupTXT(ctx, query)
	if err != nil {
		err, ok := err.(*net.DNSError)
		if !ok {
			return "", err
		}

		if err.IsNotFound {
			return "", nil
		}
	}

	var reason string
	for _, r := range reasons {
		reason = r + " " + reason
	}
	return reason, nil
}

func (r *RBL) formatQuery(ip string) string {
	reversedIP := r.reverseIP(ip)
	return reversedIP + "." + r.Domain
}

func (r *RBL) reverseIP(ip string) string {
	octets := strings.Split(ip, ".")
	reversedIP := octets[len(octets)-1]
	for i := len(octets) - 2; i >= 0; i-- {
		reversedIP += r.Separator + octets[i]
	}
	return reversedIP
}
