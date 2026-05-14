package output

import (
	"fmt"
	"strings"

	"mdnsmap/internal/asset"
)

func PrintServices(assets []asset.Asset) {
	if len(assets) == 0 {
		fmt.Println("No services found.")
		return
	}

	serviceMap := make(map[string][]asset.Asset)
	for _, a := range assets {
		portKey := fmt.Sprintf("%d/tcp %s", a.Port, getServiceBaseName(a.Service))
		serviceMap[portKey] = append(serviceMap[portKey], a)
	}

	fmt.Println("services:")
	for key, items := range serviceMap {
		for _, item := range items {
			nameWithNote := item.Name
			if strings.Contains(item.Service, "_afpovertcp") && strings.Contains(item.Name, "AFP") {
			} else if strings.Contains(item.Service, "_device-info") && strings.Contains(item.Name, "AFP") {
				nameWithNote = item.Name + "(AFP)"
			}
			fmt.Printf("%s:\n", key)
			fmt.Printf("Name=%s\n", nameWithNote)
			for _, ipv4 := range item.IPv4 {
				fmt.Printf("IPv4=%s\n", ipv4)
			}
			for _, ipv6 := range item.IPv6 {
				if strings.HasPrefix(ipv6, "fe80:") {
					fmt.Printf("IPv6=%s\n", ipv6)
				}
			}
			fmt.Printf("Hostname=%s\n", item.Hostname)
			if item.TTL > 0 {
				fmt.Printf("TTL=%d\n", item.TTL)
			}
			if len(item.TXT) > 0 {
				fmt.Printf("%s\n", strings.Join(item.TXT, ","))
			}
			fmt.Println()
		}
	}

	fmt.Println("answers:")
	fmt.Println("PTR:")
	ptrSet := make(map[string]bool)
	for _, a := range assets {
		ptrName := getServicePtrName(a.Service)
		if !ptrSet[ptrName] {
			fmt.Printf("%s\n", ptrName)
			ptrSet[ptrName] = true
		}
	}
}

func getServiceBaseName(service string) string {
	if strings.Contains(service, "_workstation") {
		return "workstation"
	} else if strings.Contains(service, "_http") {
		return "http"
	} else if strings.Contains(service, "_https") {
		return "https"
	} else if strings.Contains(service, "_smb") {
		return "smb"
	} else if strings.Contains(service, "_qdiscover") {
		return "qdiscover"
	} else if strings.Contains(service, "_device-info") {
		return "device-info"
	} else if strings.Contains(service, "_afpovertcp") {
		return "afpovertcp"
	} else if strings.Contains(service, "_ssh") {
		return "ssh"
	} else if strings.Contains(service, "_ipp") {
		return "ipp"
	}
	base := strings.TrimSuffix(service, ".local.")
	if idx := strings.Index(base, "._"); idx > 0 {
		base = base[idx+2:]
	}
	return base
}

func getServicePtrName(service string) string {
	service = strings.TrimSuffix(service, ".")
	if strings.HasSuffix(service, ".local") {
		return service + "."
	}
	if idx := strings.Index(service, "._"); idx > 0 {
		return service[idx:] + "."
	}
	return service
}
