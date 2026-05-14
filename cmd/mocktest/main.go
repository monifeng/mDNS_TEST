package main

import (
	"fmt"
	"mdnsmap/internal/asset"
	"mdnsmap/internal/output"
)

func main() {
	mockAssets := []asset.Asset{
		{
			IP:       "192.168.1.10",
			Port:     9,
			Host:     "slw-nas.local.",
			Service:  "_workstation._tcp.local.",
			Name:     "slw-nas",
			IPv4:     []string{"192.168.1.10"},
			IPv6:     []string{"fe80::265e:beff:fe69:a313"},
			Hostname: "slw-nas.local.",
			TTL:      10,
			TXT:      []string{},
			Banner:   asset.BuildBanner([]string{}, "_workstation._tcp.local."),
		},
		{
			IP:       "192.168.1.10",
			Port:     5000,
			Host:     "slw-nas.local.",
			Service:  "_http._tcp.local.",
			Name:     "slw-nas",
			IPv4:     []string{"192.168.1.10"},
			IPv6:     []string{"fe80::265e:beff:fe69:a313"},
			Hostname: "slw-nas.local.",
			TTL:      10,
			TXT:      []string{"path=/"},
			Banner:   asset.BuildBanner([]string{"path=/"}, "_http._tcp.local."),
		},
		{
			IP:       "192.168.1.10",
			Port:     445,
			Host:     "slw-nas.local.",
			Service:  "_smb._tcp.local.",
			Name:     "slw-nas",
			IPv4:     []string{"192.168.1.10"},
			IPv6:     []string{"fe80::265e:beff:fe69:a313"},
			Hostname: "slw-nas.local.",
			TTL:      10,
			TXT:      []string{},
			Banner:   asset.BuildBanner([]string{}, "_smb._tcp.local."),
		},
		{
			IP:       "192.168.1.10",
			Port:     5000,
			Host:     "slw-nas.local.",
			Service:  "_qdiscover._tcp.local.",
			Name:     "slw-nas",
			IPv4:     []string{"192.168.1.10"},
			IPv6:     []string{"fe80::265e:beff:fe69:a313"},
			Hostname: "slw-nas.local.",
			TTL:      10,
			TXT:      []string{"accessType=https", "accessPort=86", "model=TS-X64", "displayModel=TS-464C", "fwVer=5.2.9", "fwBuildNum=20260214"},
			Banner:   asset.BuildBanner([]string{"accessType=https", "accessPort=86", "model=TS-X64", "displayModel=TS-464C", "fwVer=5.2.9", "fwBuildNum=20260214"}, "_qdiscover._tcp.local."),
		},
		{
			IP:       "192.168.1.10",
			Port:     0,
			Host:     "slw-nas.local.",
			Service:  "_device-info._tcp.local.",
			Name:     "slw-nas(AFP)",
			IPv4:     []string{"192.168.1.10"},
			IPv6:     []string{"fe80::265e:beff:fe69:a313"},
			Hostname: "slw-nas.local.",
			TTL:      10,
			TXT:      []string{"model=Xserve"},
			Banner:   asset.BuildBanner([]string{"model=Xserve"}, "_device-info._tcp.local."),
		},
		{
			IP:       "192.168.1.10",
			Port:     548,
			Host:     "slw-nas.local.",
			Service:  "_afpovertcp._tcp.local.",
			Name:     "slw-nas(AFP)",
			IPv4:     []string{"192.168.1.10"},
			IPv6:     []string{"fe80::265e:beff:fe69:a313"},
			Hostname: "slw-nas.local.",
			TTL:      10,
			TXT:      []string{},
			Banner:   asset.BuildBanner([]string{}, "_afpovertcp._tcp.local."),
		},
	}

	fmt.Println("=== 测试 services 格式输出 ===")
	output.PrintServices(mockAssets)
	fmt.Println()

	fmt.Println("=== 测试 JSON 格式输出 ===")
	_ = output.PrintJSON(mockAssets, true)
	fmt.Println()

	fmt.Println("=== 测试数据集格式输出 ===")
	_ = output.PrintDataset(mockAssets, true)
}
