package netutil

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ParseCIDR(cidr string) (*net.IPNet, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	return ipnet, nil
}

func IPInCIDR(ipStr string, ipnet *net.IPNet) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ipnet.Contains(ip)
}

func ParsePorts(portsStr string) (map[int]bool, error) {
	ports := make(map[int]bool)
	parts := strings.Split(portsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		hyphenIndex := strings.Index(part, "-")
		if hyphenIndex > 0 {
			startStr := strings.TrimSpace(part[:hyphenIndex])
			endStr := strings.TrimSpace(part[hyphenIndex+1:])
			start, err := strconv.Atoi(startStr)
			if err != nil {
				return nil, fmt.Errorf("invalid port range start: %s", startStr)
			}
			end, err := strconv.Atoi(endStr)
			if err != nil {
				return nil, fmt.Errorf("invalid port range end: %s", endStr)
			}
			if start < 1 || end > 65535 || start > end {
				return nil, fmt.Errorf("invalid port range: %d-%d", start, end)
			}
			for p := start; p <= end; p++ {
				ports[p] = true
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", part)
			}
			if p < 1 || p > 65535 {
				return nil, fmt.Errorf("port out of range: %d", p)
			}
			ports[p] = true
		}
	}
	if len(ports) == 0 {
		return nil, fmt.Errorf("no ports specified")
	}
	return ports, nil
}

func ParseServices(servicesStr string) []string {
	if servicesStr == "" {
		return nil
	}
	parts := strings.Split(servicesStr, ",")
	var res []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			res = append(res, p)
		}
	}
	return res
}
