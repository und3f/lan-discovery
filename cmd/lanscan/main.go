package main

import (
	"fmt"
	"log"
	"os"

	"github.com/und3f/lanscan/scanner"
)

const SCAN_TIMES = 3

func main() {
	hosts := make(map[string]scanner.Host)

	scanRange, err := scanner.ParseCIDR(os.Args[1])
	if err != nil {
		log.Fatal("Failed to parse scanning range %s", err)
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
