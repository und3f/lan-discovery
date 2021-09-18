package scanner

import (
	"testing"
)

const subnetworkString = "127.0.0.1/32"

func TestPingScanner(t *testing.T) {
	scanRange, err := ParseCIDR(subnetworkString)
	if err != nil {
		t.Fatalf("Failed to parse: %s", err)
	}

	scanner := NewPingScanner()
	if scanner == nil {
		t.Fatal("Failed to initialize ping scanner")
	}

	if err := scanner.Scan(scanRange); err != nil {
		t.Errorf("Failed to scan network: %s", err)
	}

	/*
		if result == nil {
			t.Errorf("Scan returned nil")
		}
	*/
}
