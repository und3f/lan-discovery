package scanner

import (
	"fmt"
	"net"
)

type Host struct {
	IP       net.IP
	Hostname string
}

func (h *Host) String() string {
	return fmt.Sprintf("%s", h.IP.String())
}

type HostFoundHandler func(host Host)
type ScanEndedHandler func()

type Scanner interface {
	Scan(Range) error
	SetHostFoundHandler(handler HostFoundHandler)
	SetScanEndedHandler(handler ScanEndedHandler)
}

type ScannerEvents struct {
	hostFoundHandler HostFoundHandler
	scanEndedHandler ScanEndedHandler
}

func (se *ScannerEvents) SetHostFoundHandler(handler HostFoundHandler) {
	se.hostFoundHandler = handler
}

func (se *ScannerEvents) SetScanEndedHandler(handler ScanEndedHandler) {
	se.scanEndedHandler = handler
}

func (se *ScannerEvents) InitEmpty() {
	se.hostFoundHandler = func(host Host) {}
	se.scanEndedHandler = func() {}
}
