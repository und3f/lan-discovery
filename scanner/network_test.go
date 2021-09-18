package scanner

import (
	"net"
	"testing"
)

func TestRangeParser(t *testing.T) {
	if r, err := ParseCIDR("127.17.239.1/14"); err != nil {
		t.Errorf("Failed to parse: %s", err)
	} else {
		if !r.start.Equal(net.ParseIP("127.16.0.0")) {
			t.Errorf("Wrong start ip address %s", r.start)
		}
		if !r.end.Equal(net.ParseIP("127.19.255.255")) {
			t.Errorf("Wrong end ip address %s", r.end)
		}
	}
}

func TestRangeIterator(t *testing.T) {
	if r, err := ParseCIDR("127.0.0.1/24"); err != nil {
		t.Errorf("Failed to parse: %s", err)
	} else {
		count := 1
		it := r.createIterator()
		firstIP := it.GetNext()
		var lastIP net.IP
		for it.HasNext() {
			lastIP = it.GetNext()
			count++
		}
		if count != 256 {
			t.Errorf("Iterated %d times while expected 256", count)
		}
		if !firstIP.Equal(net.ParseIP("127.0.0.0")) {
			t.Errorf("Wrong first address %s", firstIP)
		}
		if !lastIP.Equal(net.ParseIP("127.0.0.255")) {
			t.Errorf("Wrong last address %s", lastIP)
		}
	}
}
