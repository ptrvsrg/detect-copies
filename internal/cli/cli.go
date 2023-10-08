package cli

import (
	"flag"
	"fmt"
	"net"
	"os"
)

// Parse function is used to parse command-line arguments and return the address and port.
func Parse() (addr string, port int) {
	// Create options
	flag.StringVar(&addr, "address", "", "Multicast group address")
	flag.IntVar(&port, "port", 0, "Multicast group port")

	// Parse command-line arguments
	flag.Parse()

	// Check if required options are provided
	seen := make(map[string]bool)
	flag.Visit(func(flag *flag.Flag) {
		seen[flag.Name] = true
	})
	if !seen["address"] || !seen["port"] {
		fmt.Println("Missing required flags: -address, -port")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Validate the address
	validateAddress(addr)

	// Validate the port
	validatePort(port)

	return
}

// validateAddress validates whether the provided address is a valid IP address.
func validateAddress(addr string) {
	_, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		fmt.Println(addr, "is not a valid IP address")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

// validatePort validates whether the provided port is within a valid range.
func validatePort(port int) {
	if port < 0 || port > 65535 {
		fmt.Println(port, "is not a valid port")
		flag.PrintDefaults()
		os.Exit(1)
	}
}
