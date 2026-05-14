package discovery

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"mdnsmap/internal/asset"
	"mdnsmap/internal/config"
	"mdnsmap/internal/netutil"

	"github.com/grandcat/zeroconf"
)

func Scan(ctx context.Context, cfg *config.Config, ipnet *net.IPNet, allowedPorts map[int]bool, services []string) ([]asset.Asset, error) {
	if services == nil || len(services) == 0 {
		services = config.DefaultServices
	}
	resultsCh := make(chan *zeroconf.ServiceEntry, 100)
	var wg sync.WaitGroup
	for _, svc := range services {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			browseOne(ctx, s, resultsCh)
		}(svc)
	}
	go func() {
		wg.Wait()
		close(resultsCh)
	}()
	seen := make(map[string]bool)
	var assets []asset.Asset
	for entry := range resultsCh {
		for _, ip4 := range entry.AddrIPv4 {
			ipStr := ip4.String()
			if !netutil.IPInCIDR(ipStr, ipnet) {
				continue
			}
			if !allowedPorts[entry.Port] {
				continue
			}
			svcType := entry.Service
			if !strings.HasSuffix(svcType, ".") {
				svcType += "."
			}
			key := fmt.Sprintf("%s|%s|%s|%d|%s", svcType, entry.Instance, entry.HostName, entry.Port, ipStr)
			if seen[key] {
				continue
			}
			seen[key] = true
			a := asset.Asset{
				IP:       ipStr,
				Port:     entry.Port,
				Host:     entry.HostName,
				Service:  svcType,
				Name:     entry.Instance,
				IPv4:     ipsToStrings(entry.AddrIPv4),
				IPv6:     ipsToStrings(entry.AddrIPv6),
				Hostname: entry.HostName,
				TTL:      entry.TTL,
				TXT:      entry.Text,
				Banner:   asset.BuildBanner(entry.Text, svcType),
			}
			assets = append(assets, a)
		}
	}
	return assets, nil
}

func ipsToStrings(ips []net.IP) []string {
	res := make([]string, 0, len(ips))
	for _, ip := range ips {
		res = append(res, ip.String())
	}
	return res
}

func browseOne(ctx context.Context, service string, resultsCh chan<- *zeroconf.ServiceEntry) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return
	}
	entries := make(chan *zeroconf.ServiceEntry)
	go func() {
		for entry := range entries {
			select {
			case resultsCh <- entry:
			case <-ctx.Done():
				return
			}
		}
	}()
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_ = resolver.Browse(timeoutCtx, service, "local.", entries)
}
