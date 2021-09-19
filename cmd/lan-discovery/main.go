package main

import (
	"fmt"
	"log"
	"os"

	"github.com/und3f/lan-discovery/discovery"
	"github.com/und3f/lan-discovery/scanner"
)

const SCAN_TIMES = 3

func main() {
	var scanRange scanner.Range
	var err error

	arp, err := scanner.NewARPDiscovery()
	if err != nil {
		log.Fatalf("ARP Discovery initialization failed: %s", err)
	}

	if len(os.Args) >= 2 {
		scanRange, err = scanner.ParseCIDR(os.Args[1])
		if err != nil {
			log.Fatalf("Failed to parse scanning range: %s", err)
		}
	} else {
		scanRange, err = discovery.Interfaces()
		if err != nil {
			log.Fatalf("Failed to discover scanning range: %s", err)
		}
	}

	hs := scanner.NewHostsStorage()

	networkScanner := scanner.NewPingScanner()

	hostFoundHandler := func(host *scanner.Host) {
		if updated := hs.Update(host); updated != nil {
			fmt.Println("Discovered", host.String())
		}
	}

	networkScanner.SetHostFoundHandler(hostFoundHandler)
	arp.SetHostFoundHandler(hostFoundHandler)

	fmt.Println("Discovering hosts...")
	for i := 0; i < SCAN_TIMES; i++ {
		if err := networkScanner.Scan(scanRange); err != nil {
			log.Fatalf("Failed to start scanning: %s", err)
		}

		if err := arp.Discover(); err != nil {
			log.Fatalf("Failed to start scanning: %s", err)
		}
	}
}
