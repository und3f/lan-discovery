package main

import (
	"fmt"
	"log"
	"os"

	"github.com/und3f/lanscan/discoverer"
	"github.com/und3f/lanscan/scanner"
)

const SCAN_TIMES = 3

func main() {
	hosts := make(map[string]scanner.Host)
	var scanRange scanner.Range
	var err error

	if len(os.Args) >= 2 {
		scanRange, err = scanner.ParseCIDR(os.Args[1])
		if err != nil {
			log.Fatalf("Failed to parse scanning range: %s", err)
		}
	} else {
		scanRange, err = discoverer.Interfaces()
		if err != nil {
			log.Fatalf("Failed to discover scanning range: %s", err)
		}
	}

	networkScanner := scanner.NewPingScanner()

	networkScanner.SetHostFoundHandler(func(host scanner.Host) {
		if _, exists := hosts[host.IP.String()]; exists {
			return
		}

		hosts[host.IP.String()] = host
		fmt.Println("Discovered", host.String())
	})

	fmt.Println("Discovering hosts...")
	for i := 0; i < SCAN_TIMES; i++ {
		if err := networkScanner.Scan(scanRange); err != nil {
			log.Fatalf("Failed to start scanning: %s", err)
		}
	}
}
