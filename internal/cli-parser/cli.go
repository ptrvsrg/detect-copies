package cliparser

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func Parse() (addr string, port int) {
	flag.StringVar(&addr, "address", "", "Multicast group address")
	flag.IntVar(&port, "port", 0, "Multicast group port")

	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(flag *flag.Flag) {
		seen[flag.Name] = true
	})
	if !seen["address"] || !seen["port"] {
		fmt.Println("Missing required flags: -address, -port")
		flag.PrintDefaults()
		os.Exit(1)
	}

	validateAddress(addr)
	validatePort(port)

	return
}

func validateAddress(addr string) {
	ip := net.ParseIP(addr)
	if ip == nil {
		fmt.Println(addr, "is not valid IP address")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func validatePort(port int) {
	if port < 0 || port > 65535 {
		fmt.Println(port, "is not valid port")
		flag.PrintDefaults()
		os.Exit(1)
	}
}