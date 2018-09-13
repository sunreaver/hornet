package config

// OrmEngineConfig OrmEngineConfig
type OrmEngineConfig struct {
	// Type
	Type string
	// Dialect 驱动: mysql, sqlite3 ...
	Dialect string
	// Uris db hosts
	Uris []string
}

func (oec *OrmEngineConfig) Verify() error {
	switch oec.Dialect {
	case "mysql", "sqlite3", "postgres":
	default:
		return NoDialect
	}
	if len(oec.Uris) == 0 {
		return NoUris
	}

	return nil
}
