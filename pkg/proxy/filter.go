package proxy

import (
	"net"
)

// ShouldBypassSamePodProxy detects if destination IP is on local loopback or local container pod network
func ShouldBypassSamePodProxy(destIP net.IP, destPort uint16, bypassPorts []uint16) bool {
	if destIP.IsLoopback() {
		return true
	}
	for _, port := range bypassPorts {
		if port == destPort {
			return true
		}
	}
	return false
}
