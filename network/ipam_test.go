package network

import (
	"net"
	"testing"
)


func TestIPAM_Allocate(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("192.168.7.1/24")
	ip, _ := ipAllocator.Allocate(ipnet)
	t.Logf("alloc ip: %v", ip)
}

func TestIPAM_Release(t *testing.T) {
	ip, ipnet, _ := net.ParseCIDR("192.168.7.1/24")
	ipAllocator.Release(ipnet, &ip)
}
