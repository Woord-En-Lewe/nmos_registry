package mdns

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/grandcat/zeroconf"
)

type Announcer interface {
	Start(ctx context.Context) error
	Stop() error
	AnnounceRegistrationAPI() error
	AnnounceQueryAPI() error
}

type announcer struct {
	config   *Config
	services []*zeroconf.Server
	mu       sync.Mutex
	ifaces   []net.Interface
}

func NewAnnouncer(cfg *Config) (Announcer, error) {
	ifaces, err := allInterfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to find network interfaces: %w", err)
	}

	return &announcer{
		config:   cfg,
		services: make([]*zeroconf.Server, 0),
		ifaces:   ifaces,
	}, nil
}

func allInterfaces() ([]net.Interface, error) {
	var ifaces []net.Interface
	addrs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range addrs {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		if !hasIPv4Address(&iface) {
			continue
		}

		ifaces = append(ifaces, iface)
	}
	if len(ifaces) == 0 {
		return nil, fmt.Errorf("no suitable network interface with IPv4 found")
	}
	return ifaces, nil
}

func hasIPv4Address(iface *net.Interface) bool {
	addrs, err := iface.Addrs()
	if err != nil {
		return false
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() {
				return true
			}
		}
	}
	return false
}

func (a *announcer) Start(ctx context.Context) error {
	if err := a.AnnounceRegistrationAPI(); err != nil {
		return fmt.Errorf("failed to announce registration API: %w", err)
	}
	if err := a.AnnounceQueryAPI(); err != nil {
		return fmt.Errorf("failed to announce query API: %w", err)
	}
	return nil
}

func (a *announcer) Stop() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, svc := range a.services {
		svc.Shutdown()
	}
	a.services = nil

	return nil
}

func (a *announcer) AnnounceRegistrationAPI() error {
	return a.registerService(a.config.RegistrationService())
}

func (a *announcer) AnnounceQueryAPI() error {
	return a.registerService(a.config.QueryService())
}

func (a *announcer) registerService(svc *Service) error {
	var txtEntries []string
	for k, v := range svc.TXTRecords {
		txtEntries = append(txtEntries, fmt.Sprintf("%s=%s", k, v))
	}

	instanceName := svc.InstanceName
	if instanceName == "" {
		instanceName = fmt.Sprintf("NMOS Registry (%s)", a.config.Hostname)
	}

	for _, iface := range a.ifaces {
		server, err := zeroconf.Register(
			instanceName,
			svc.ServiceType,
			Domain,
			svc.Port,
			txtEntries,
			[]net.Interface{iface},
		)
		if err != nil {
			log.Printf("mDNS: Failed to announce %s on interface %s: %v", svc.ServiceType, iface.Name, err)
			continue
		}

		a.mu.Lock()
		a.services = append(a.services, server)
		a.mu.Unlock()

		log.Printf("mDNS: Announced %s on interface %s port %d", svc.ServiceType, iface.Name, svc.Port)
	}

	return nil
}
