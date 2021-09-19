package scanner

import (
	"net"
	"sort"
	"testing"
)

const SampleHostname = "Somehostname"

func TestHost(t *testing.T) {
	ip1 := net.ParseIP("127.0.0.1")
	// ip2 := net.ParseIP("127.0.0.2")
	ip3 := net.ParseIP("127.0.0.1")

	host := NewHost(ip1)

	if host == nil {
		t.Errorf("failed to create host")
	}

	host2 := NewHost(ip3)
	host2.Hostname = SampleHostname
	host.Update(host2)
	if want, got := SampleHostname, host.Hostname; want != got {
		t.Errorf("wanted %s, got %s", want, got)
	}
}

func TestHostsStorage(t *testing.T) {
	ip1 := net.ParseIP("127.0.0.1")
	ip2 := net.ParseIP("127.0.0.2")

	host1 := NewHost(ip1)
	host2 := NewHost(ip1)
	host2.Hostname = SampleHostname

	hs := NewHostsStorage()

	// Add single host to storage
	hs.Update(host1)

	// Retrieve hosts
	hosts := hs.GetHosts()
	if want, got := 1, len(hosts); want != got {
		t.Errorf("wanted %d, got %d", want, got)
	}
	if want, got := host1.Identifier(), hosts[0].Identifier(); want != got {
		t.Errorf("wanted %s, got %s", want, got)
	}
	if want, got := "", hosts[0].Hostname; want != got {
		t.Errorf("wanted %s, got %s", want, got)
	}

	// Update existing host info
	hs.Update(host2)

	// Check host info updated
	hosts = hs.GetHosts()
	if want, got := 1, len(hosts); want != got {
		t.Errorf("wanted %d, got %d", want, got)
	}
	if want, got := SampleHostname, hosts[0].Hostname; want != got {
		t.Errorf(`wanted "%s", got "%s"`, want, got)
	}

	// Add second host
	host3 := NewHost(ip2)
	hs.Update(host3)

	// Check host info updated
	hosts = hs.GetHosts()
	if want, got := 2, len(hosts); want != got {
		t.Errorf("wanted %d, got %d", want, got)
	}

	sort.Slice(hosts, func(i, j int) bool { return hosts[i].IP.String() < hosts[j].IP.String() })
	if want, got := host1.Identifier(), hosts[0].Identifier(); want != got {
		t.Errorf("wanted %s, got %s", want, got)
	}
	if want, got := host3.Identifier(), hosts[1].Identifier(); want != got {
		t.Errorf("wanted %s, got %s", want, got)
	}
}
