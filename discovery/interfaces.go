package discovery

import (
	"log"
	"net"

	"github.com/und3f/lan-discovery/scanner"
)

func Interfaces() ([]*scanner.Host, scanner.Range, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}

	mr := scanner.MultipleRanges{}
	var hosts []*scanner.Host

	for _, interf := range interfaces {
		if interf.Flags&net.FlagUp == 0 || interf.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := interf.Addrs()
		if err != nil {
			log.Printf("Failed to obtain address for interface %v: %v\n", interf, err)
			continue
		}
		for _, addr := range addrs {
			if ip, _, err := net.ParseCIDR(addr.String()); err != nil {
				log.Printf("Range parsing of address %s failed: %v\n", addr, err)
				continue
			} else {
				host := scanner.NewHost(ip)
				host.HardwareAddr = interf.HardwareAddr
				hosts = append(hosts, host)
			}

			if r, err := scanner.ParseCIDR(addr.String()); err != nil {
				log.Printf("Range parsing of address %s failed: %v\n", addr, err)
				continue
			} else {
				mr.Append(r)
			}
		}
	}

	return hosts, &mr, nil
}
