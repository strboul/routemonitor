package main

import (
	"testing"
)

func TestGetCidrIps(t *testing.T) {
	ips, err := getCidrIps("10.0.0.0/24")
	if err != nil {
		t.Fatal(err)
	}
	if len(ips) != 256 {
		t.Errorf("current %v, expected %v", len(ips), 256)
	}
}
