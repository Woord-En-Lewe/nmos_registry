package mdns

import (
	"time"
)

type Config struct {
	Hostname         string
	RegistrationPort int
	QueryPort        int
	Priority         int
	TTL              time.Duration
}

type Service struct {
	ServiceType  string
	InstanceName string
	Port         int
	Priority     int
	Weight       int
	TTL          uint32
	TXTRecords   map[string]string
}

type DiscoveredService struct {
	InstanceName string
	ServiceType  string
	Domain       string
	Port         int
	TXTRecords   map[string]string
	Addr         string
}

func NewConfig(hostname string, regPort, queryPort int) *Config {
	return &Config{
		Hostname:         hostname,
		RegistrationPort: regPort,
		QueryPort:        queryPort,
		Priority:         DefaultPriority,
		TTL:              DefaultTTL * time.Second,
	}
}

func (c *Config) RegistrationService() *Service {
	return &Service{
		ServiceType:  RegistrationServiceType,
		InstanceName: c.Hostname + "._nmos-registration._tcp." + Domain,
		Port:         c.RegistrationPort,
		Priority:     c.Priority,
		Weight:       DefaultWeight,
		TTL:          uint32(c.TTL.Seconds()),
		TXTRecords: map[string]string{
			"api_ver": RegistryAPIVersion,
			"pri":     formatPriority(c.Priority),
		},
	}
}

func (c *Config) QueryService() *Service {
	return &Service{
		ServiceType:  QueryServiceType,
		InstanceName: c.Hostname + "._nmos-query._tcp." + Domain,
		Port:         c.QueryPort,
		Priority:     c.Priority,
		Weight:       DefaultWeight,
		TTL:          uint32(c.TTL.Seconds()),
		TXTRecords: map[string]string{
			"api_ver": RegistryAPIVersion,
			"pri":     formatPriority(c.Priority),
		},
	}
}

func formatPriority(p int) string {
	if p > 9999 {
		p = 9999
	}
	if p < 0 {
		p = 0
	}
	return string(rune('0'+p/1000)) + string(rune('0'+(p/100)%10)) + string(rune('0'+(p/10)%10)) + string(rune('0'+p%10))
}
