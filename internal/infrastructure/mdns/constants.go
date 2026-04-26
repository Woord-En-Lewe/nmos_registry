package mdns

const (
	RegistrationServiceType = "_nmos-register._tcp"
	QueryServiceType        = "_nmos-query._tcp"
	RegistryServiceType     = "_nmos-registry._tcp"
	Domain                  = "local"

	RegistryAPIVersion = "v1.3"
	DefaultPriority    = 100
	DefaultWeight      = 0
	DefaultTTL         = 120
)
