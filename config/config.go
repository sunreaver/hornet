package config

// OrmEngineConfig OrmEngineConfig
type OrmEngineConfig struct {
	// Dialect 驱动: mysql, sqlite3 ...
	Dialect string
	// Uris db hosts
	Uris []string
}

// Verify Verify
func (oec *OrmEngineConfig) Verify() error {
	switch oec.Dialect {
	case "mysql", "sqlite3", "postgres":
	default:
		return ErrNoDialect
	}
	if len(oec.Uris) == 0 {
		return ErrNoUris
	}

	return nil
}
