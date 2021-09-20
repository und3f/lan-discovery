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
	hs := scanner.NewHostsStorage()

	observer := &NewHostsObserver{}
	hs.HostUpdatePublisher.Subscribe(observer)

	fmt.Println("Discovering hosts...")

	var arp scanner.ARPDiscovery
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
		var hosts []*scanner.Host
		hosts, scanRange, err = discovery.Interfaces()
		if err != nil {
			log.Fatalf("Failed to discover scanning range: %s", err)
		}
		for _, host := range hosts {
			hs.Update(host)
		}
	}

	networkScanner := scanner.NewPingScanner()

	hostFoundHandler := func(host *scanner.Host) { hs.Update(host) }

	networkScanner.SetHostFoundHandler(hostFoundHandler)
	arp.SetHostFoundHandler(hostFoundHandler)

	for i := 0; i < SCAN_TIMES; i++ {
		if err := networkScanner.Scan(scanRange); err != nil {
			log.Fatalf("Failed to start scanning: %s", err)
		}

		if err := arp.Discover(); err != nil {
			log.Fatalf("Failed to start scanning: %s", err)
		}
	}

	PrintExternalHosts(&hs)
}

func PrintExternalHosts(hs *scanner.HostsStorage) {
	var hosts []*scanner.Host
	for _, host := range hs.GetHosts() {
		if len(host.HardwareAddr) == 0 {
			hosts = append(hosts, host)
		}
	}

	if len(hosts) > 0 {
		fmt.Println("Next hosts are out of the LAN", hosts)
	}
}

type NewHostsObserver struct{}

func (observer *NewHostsObserver) Update(host *scanner.Host) *scanner.Host {
	if len(host.HardwareAddr) == 0 {
		return nil
	}

	fmt.Println(host.String())
	return nil
}
