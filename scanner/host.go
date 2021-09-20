package scanner

import (
	"fmt"
	"net"
)

type Host struct {
	HardwareAddr net.HardwareAddr
	IP           net.IP
	Hostname     string
}

func NewHost(ip net.IP) *Host {
	return &Host{
		IP: ip,
	}
}

func (h *Host) String() string {
	return fmt.Sprintf("%s (%s) %s", h.IP.String(), h.HardwareAddr, h.Hostname)
}

func (h *Host) Identifier() string {
	return h.IP.String()
}

func (h *Host) Update(host *Host) *Host {
	var changed bool

	if len(host.HardwareAddr) > 0 {
		if len(h.HardwareAddr) == 0 || h.HardwareAddr.String() != host.HardwareAddr.String() {
			h.HardwareAddr = host.HardwareAddr
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
	HostUpdatePublisher HostUpdatePublisher

	hosts map[string]*Host
}

func NewHostsStorage() HostsStorage {
	return HostsStorage{
		hosts: make(map[string]*Host),
	}
}

func (hs *HostsStorage) Update(h *Host) (retHost *Host) {
	if storedHost, exists := hs.hosts[h.Identifier()]; exists {
		retHost = storedHost.Update(h)
	} else {
		hs.hosts[h.Identifier()] = h
		retHost = h
	}

	if retHost != nil {
		hs.HostUpdatePublisher.NotifySubscribers(retHost)
	}

	return
}

func (hs *HostsStorage) GetHosts() []*Host {
	hosts := make([]*Host, 0, len(hs.hosts))
	for _, host := range hs.hosts {
		hosts = append(hosts, host)
	}

	return hosts
}

type HostUpdateSubscriber interface {
	Update(*Host) *Host
}

type HostUpdatePublisher struct {
	subscribers []HostUpdateSubscriber
}

func (hp *HostUpdatePublisher) Subscribe(s HostUpdateSubscriber) {
	hp.subscribers = append(hp.subscribers, s)
}

func (hp *HostUpdatePublisher) Unsubscribe(s HostUpdateSubscriber) {
	for i, subscriber := range hp.subscribers {
		if subscriber == s {
			hp.subscribers = append(hp.subscribers[:i], hp.subscribers[i+1:]...)
			break
		}
	}
}

func (hp *HostUpdatePublisher) NotifySubscribers(h *Host) {
	for _, subscriber := range hp.subscribers {
		subscriber.Update(h)
	}
}
