package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// 1. Get the filename for the target IP list from command-line arguments
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: Please provide a filename listing target IPs as an argument!")
		fmt.Fprintln(os.Stderr, "Usage: ./checkips_from_file <target_IP_list_filename>")
		os.Exit(1) // Usage error
	}
	targetIPsFilename := os.Args[1]

	// 2. First, read all CIDR ranges from stdin and store them in memory.
	//    This is necessary because stdin can only be read once, and we need to check each target IP against all CIDRs.
	var cidrNetworks []*net.IPNet
	stdinScanner := bufio.NewScanner(os.Stdin)
	fmt.Fprintln(os.Stderr, "INFO: Reading CIDR list from stdin...")

	for stdinScanner.Scan() {
		cidrStr := strings.TrimSpace(stdinScanner.Text())
		if cidrStr == "" {
			continue
		}
		_, network, err := net.ParseCIDR(cidrStr)
		if err != nil {
			// fmt.Fprintf(os.Stderr, "WARN: Failed to parse CIDR \"%s\" (skipping): %v\n", cidrStr, err)
			continue // Silently skip unparseable CIDRs
		}
		cidrNetworks = append(cidrNetworks, network)
	}
	if err := stdinScanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: Problem reading CIDR list from stdin:", err)
		os.Exit(2) // stdin read error
	}

	if len(cidrNetworks) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No valid CIDRs were read from stdin. Is the jq output correct?")
		os.Exit(3) // CIDR list empty
	}
	fmt.Fprintf(os.Stderr, "INFO: Loaded %d CIDR networks.\n", len(cidrNetworks))

	// 3. Next, open the specified target IP file
	targetIPsFile, err := os.Open(targetIPsFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not open target IP file \"%s\": %v\n", targetIPsFilename, err)
		os.Exit(4) // File open error
	}
	defer targetIPsFile.Close()

	// 4. Read each IP address from the file line by line and check each one
	fmt.Fprintf(os.Stderr, "INFO: Reading IPs from target file \"%s\" and starting checks...\n", targetIPsFilename)
	fileScanner := bufio.NewScanner(targetIPsFile)
	processedIPCount := 0

	for fileScanner.Scan() {
		targetIPStr := strings.TrimSpace(fileScanner.Text())
		if targetIPStr == "" {
			continue // Skip empty lines in the file
		}
		processedIPCount++

		ipToCheck := net.ParseIP(targetIPStr)
		if ipToCheck == nil {
			// Print errors for specific invalid IPs in the file to stderr
			fmt.Fprintf(os.Stderr, "WARN: Skipping invalid IP address format in target file: %s\n", targetIPStr)
			continue
		}

		foundMatchForThisIP := false
		var matchingCIDR string
		for _, network := range cidrNetworks {
			if network.Contains(ipToCheck) {
				foundMatchForThisIP = true
				matchingCIDR = network.String() // Convert net.IPNet back to string
				break
			}
		}

		if foundMatchForThisIP {
			fmt.Printf("SUCCESS: %s is within an ap-northeast-1 range (%s)!\n", targetIPStr, matchingCIDR)
		} else {
			fmt.Printf("NOT_FOUND: %s was not found in any ap-northeast-1 ranges.\n", targetIPStr)
		}
	}

	if err := fileScanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Problem reading from target IP file \"%s\": %v\n", targetIPsFilename, err)
		os.Exit(5) // Error during file read
	}
	fmt.Fprintf(os.Stderr, "INFO: Finished checking %d IP addresses.\n", processedIPCount)
}
