package main

import (
	"fmt"
	"log"
	"os"

	"github.com/und3f/lanscan/scanner"
)

func main() {
	scanRange, err := scanner.ParseCIDR(os.Args[1])
	if err != nil {
		log.Fatal("Failed to parse scanning range %s", err)
	}
	networkScanner := scanner.NewPingScanner()

	networkScanner.SetHostFoundHandler(func(host scanner.Host) {
		fmt.Println(host)
	})

	if err := networkScanner.Scan(scanRange); err != nil {
		log.Fatalf("Failed to start scanning: %s", err)
	}
}
