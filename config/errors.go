package config

import "errors"

var (
	NoUris          = errors.New("no uris")
	NoDialect       = errors.New("no dialect")
	NoAvailableHost = errors.New("no available host")
)
