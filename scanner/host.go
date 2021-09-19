package scanner

import (
	"fmt"
	"net"
)

type HostUpdateSubscriber interface {
	update(*Host)
}

type Host struct {
	MAC      net.HardwareAddr
	IP       net.IP
	Hostname string
}

func NewHost(ip net.IP) *Host {
	return &Host{
		IP: ip,
	}
}

func (h *Host) String() string {
	return fmt.Sprintf("%s (%s) %s", h.IP.String(), h.MAC, h.Hostname)
}

func (h *Host) Identifier() string {
	return h.IP.String()
}

func (h *Host) Update(host *Host) *Host {
	var changed bool

	if len(host.MAC) > 0 {
		if len(h.MAC) == 0 || h.MAC.String() != host.MAC.String() {
			h.MAC = host.MAC
			changed = true
		}
	}

	if len(host.Hostname) > 0 && host.Hostname != h.Hostname {
		h.Hostname = host.Hostname
		changed = true
	}

	if changed {
		return h
	}
	return nil
}

type HostsStorage struct {
	hosts map[string]*Host
}

func NewHostsStorage() HostsStorage {
	return HostsStorage{
		hosts: make(map[string]*Host),
	}
}

func (hs *HostsStorage) Update(h *Host) *Host {
	if storedHost, exists := hs.hosts[h.Identifier()]; exists {
		return storedHost.Update(h)
	} else {
		hs.hosts[h.Identifier()] = h
		return h
	}
}

func (hs *HostsStorage) GetHosts() []*Host {
	hosts := make([]*Host, 0, len(hs.hosts))
	for _, host := range hs.hosts {
		hosts = append(hosts, host)
	}

	return hosts
}
