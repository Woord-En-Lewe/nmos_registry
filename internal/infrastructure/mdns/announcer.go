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
	iface    *net.Interface
}

func NewAnnouncer(cfg *Config) (Announcer, error) {
	iface, err := net.InterfaceByName("")
	if err != nil {
		localIface, err := defaultInterface()
		if err != nil {
			return nil, fmt.Errorf("failed to find network interface: %w", err)
		}
		iface = localIface
	}

	return &announcer{
		config:   cfg,
		services: make([]*zeroconf.Server, 0),
		iface:    iface,
	}, nil
}

func defaultInterface() (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			return &iface, nil
		}
	}
	return nil, fmt.Errorf("no suitable network interface found")
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

	server, err := zeroconf.Register(
		instanceName,
		svc.ServiceType,
		Domain,
		svc.Port,
		txtEntries,
		[]net.Interface{*a.iface},
	)
	if err != nil {
		return fmt.Errorf("failed to register mDNS service: %w", err)
	}

	a.mu.Lock()
	a.services = append(a.services, server)
	a.mu.Unlock()

	log.Printf("mDNS: Announced %s on port %d", svc.ServiceType, svc.Port)
	return nil
}
