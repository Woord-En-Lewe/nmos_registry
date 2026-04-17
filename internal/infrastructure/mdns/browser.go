package mdns

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
)

type Browser interface {
	DiscoverRegistrationAPI(ctx context.Context) ([]DiscoveredService, error)
	DiscoverQueryAPI(ctx context.Context) ([]DiscoveredService, error)
}

type browser struct {
	ifaces []net.Interface
}

func NewBrowser() (Browser, error) {
	ifaces, err := allInterfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to find network interfaces: %w", err)
	}

	return &browser{
		ifaces: ifaces,
	}, nil
}

func (b *browser) DiscoverRegistrationAPI(ctx context.Context) ([]DiscoveredService, error) {
	return b.discover(ctx, RegistrationServiceType)
}

func (b *browser) DiscoverQueryAPI(ctx context.Context) ([]DiscoveredService, error) {
	return b.discover(ctx, QueryServiceType)
}

func (b *browser) discover(ctx context.Context, serviceType string) ([]DiscoveredService, error) {
	resolver, err := zeroconf.NewResolver(zeroconf.SelectIfaces(b.ifaces))
	if err != nil {
		return nil, fmt.Errorf("failed to create resolver: %w", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(entryCh <-chan *zeroconf.ServiceEntry, service string) {
		for entry := range entryCh {
			addr := ""
			if len(entry.AddrIPv4) > 0 {
				addr = entry.AddrIPv4[0].String()
			}
			log.Printf("mDNS: Discovered %s: %s (%s:%d)", service, entry.Instance, addr, entry.Port)
		}
	}(entries, serviceType)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := resolver.Browse(ctx, serviceType, Domain, entries); err != nil {
		return nil, fmt.Errorf("failed to browse for %s: %w", serviceType, err)
	}

	<-ctx.Done()

	var services []DiscoveredService
	for entry := range entries {
		addr := ""
		if len(entry.AddrIPv4) > 0 {
			addr = entry.AddrIPv4[0].String()
		}
		svc := DiscoveredService{
			InstanceName: entry.Instance,
			ServiceType:  entry.Service,
			Domain:       entry.Domain,
			Port:         entry.Port,
			Addr:         addr,
			TXTRecords:   parseTXTRecords(entry.Text),
		}
		services = append(services, svc)
	}

	return services, nil
}

func parseTXTRecords(txt []string) map[string]string {
	records := make(map[string]string)
	for _, entry := range txt {
		pair := strings.SplitN(entry, "=", 2)
		if len(pair) == 2 {
			records[pair[0]] = pair[1]
		}
	}
	return records
}
