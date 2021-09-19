package scanner

import (
	"github.com/prometheus/procfs"
)

type ARPDiscovery struct {
	procfs procfs.FS
	arp    []procfs.ARPEntry

	ScannerEvents
}

const PROC_PATH = "/proc"

func NewARPDiscovery() (ARPDiscovery, error) {
	discovery := ARPDiscovery{}
	discovery.InitEmpty()

	if fs, err := procfs.NewFS(PROC_PATH); err != nil {
		return discovery, err
	} else {
		discovery.procfs = fs
	}
	return discovery, nil
}

func (discovery *ARPDiscovery) Discover() error {
	if arp, err := discovery.procfs.GatherARPEntries(); err == nil {
		discovery.arp = make([]procfs.ARPEntry, 0)
		for _, entry := range arp {
			if entry.IsComplete() {
				host := NewHost(entry.IPAddr)
				host.MAC = entry.HWAddr
				discovery.hostFoundHandler(host)
			}
		}
	} else {
		return err
	}

	discovery.scanEndedHandler()
	return nil
}
