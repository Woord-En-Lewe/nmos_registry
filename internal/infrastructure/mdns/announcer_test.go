package mdns

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig("test-host", 8080, 8081)

	if cfg.Hostname != "test-host" {
		t.Errorf("expected hostname 'test-host', got '%s'", cfg.Hostname)
	}
	if cfg.RegistrationPort != 8080 {
		t.Errorf("expected registration port 8080, got %d", cfg.RegistrationPort)
	}
	if cfg.QueryPort != 8081 {
		t.Errorf("expected query port 8081, got %d", cfg.QueryPort)
	}
	if cfg.Priority != DefaultPriority {
		t.Errorf("expected priority %d, got %d", DefaultPriority, cfg.Priority)
	}
}

func TestConfigRegistrationService(t *testing.T) {
	cfg := NewConfig("test-host", 8080, 8081)
	svc := cfg.RegistrationService()

	if svc.ServiceType != RegistrationServiceType {
		t.Errorf("expected service type '%s', got '%s'", RegistrationServiceType, svc.ServiceType)
	}
	if svc.Port != 8080 {
		t.Errorf("expected port 8080, got %d", svc.Port)
	}
	if svc.TXTRecords["api_ver"] != RegistryAPIVersion {
		t.Errorf("expected api_ver '%s', got '%s'", RegistryAPIVersion, svc.TXTRecords["api_ver"])
	}
	if svc.TXTRecords["pri"] == "" {
		t.Error("expected pri TXT record to be set")
	}
}

func TestConfigQueryService(t *testing.T) {
	cfg := NewConfig("test-host", 8080, 8081)
	svc := cfg.QueryService()

	if svc.ServiceType != QueryServiceType {
		t.Errorf("expected service type '%s', got '%s'", QueryServiceType, svc.ServiceType)
	}
	if svc.Port != 8081 {
		t.Errorf("expected port 8081, got %d", svc.Port)
	}
}

func TestParseTXTRecords(t *testing.T) {
	txt := []string{"api_ver=v1.3", "pri=0100"}
	records := parseTXTRecords(txt)

	if records["api_ver"] != "v1.3" {
		t.Errorf("expected api_ver 'v1.3', got '%s'", records["api_ver"])
	}
	if records["pri"] != "0100" {
		t.Errorf("expected pri '0100', got '%s'", records["pri"])
	}
}

func TestParseTXTRecordsEmpty(t *testing.T) {
	txt := []string{}
	records := parseTXTRecords(txt)

	if len(records) != 0 {
		t.Errorf("expected empty map, got %d entries", len(records))
	}
}
