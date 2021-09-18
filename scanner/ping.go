package scanner

import (
	"net"
	"time"

	"github.com/tatsushid/go-fastping"
)

type PingScanner struct {
	ScannerEvents

	fastPing *fastping.Pinger
}

func NewPingScanner() *PingScanner {
	ps := new(PingScanner)
	ps.ScannerEvents.InitEmpty()

	ps.fastPing = fastping.NewPinger()
	ps.SetICMP(true)
	ps.fastPing.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		ps.AddHost(addr)
	}
	ps.fastPing.OnIdle = func() {
		ps.scanEndedHandler()
	}
	return ps
}

func (scanner *PingScanner) SetICMP(useICMP bool) {
	if useICMP {
		scanner.fastPing.Network("ip")
	} else {
		scanner.fastPing.Network("udp")
	}
}

func (scanner *PingScanner) AddHost(ip *net.IPAddr) {
	var host Host
	host.ip = cloneIP(ip.IP)
	scanner.hostFoundHandler(host)
}

func (scanner *PingScanner) Scan(network Range) error {
	for it := network.createIterator(); it.HasNext(); {
		var ipaddr net.IPAddr
		ipaddr.IP = it.GetNext()

		scanner.fastPing.AddIPAddr(&ipaddr)
	}

	err := scanner.fastPing.Run()

	return err
}
