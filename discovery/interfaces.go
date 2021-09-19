package discovery

import (
	"log"
	"net"

	"github.com/und3f/lanscan/scanner"
)

func Interfaces() (scanner.Range, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	mr := scanner.MultipleRanges{}

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
			r, err := scanner.ParseCIDR(addr.String())
			if err != nil {
				log.Printf("Range parsing of address %s failed: %v\n", addr, err)
				continue
			}
			mr.Append(r)
		}
	}

	return &mr, nil
}
