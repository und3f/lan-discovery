package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sort"

	"github.com/und3f/lan-discovery/discovery"
	"github.com/und3f/lan-discovery/scanner"
)

const SCAN_TIMES = 3

func main() {
	var scanRange scanner.Range
	hs := scanner.NewHostsStorage()

	observer := &NewHostsObserver{}

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

	fmt.Print("Discovering hosts...")
	hs.HostUpdatePublisher.Subscribe(observer)

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

	hs.HostUpdatePublisher.Unsubscribe(observer)
	fmt.Println()

	PrintExternalHosts(&hs)
	PrintNetwork(&hs)
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

func IsIPLess(a, b net.IP) bool {
	if lenDiff := len(a) - len(b); lenDiff != 0 {
		return lenDiff < 0
	}
	for i, part := range a {
		if diff := part - b[i]; diff != 0 {
			return diff < 0
		}
	}
	return false
}

func PrintNetwork(hs *scanner.HostsStorage) {
	var hosts []*scanner.Host
	for _, host := range hs.GetHosts() {
		if host.IsOnline() {
			hosts = append(hosts, host)
		}
	}

	sort.Slice(hosts, func(i, j int) bool { return IsIPLess(hosts[i].IP, hosts[j].IP) })
	for _, host := range hosts {
		fmt.Println(host.String())
	}
}

type NewHostsObserver struct{}

func (observer *NewHostsObserver) Update(host *scanner.Host) *scanner.Host {
	if len(host.HardwareAddr) == 0 {
		return nil
	}

	fmt.Print(".")
	return nil
}
