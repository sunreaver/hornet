package config

import "errors"

var (
	// ErrNoUris uri为空
	ErrNoUris = errors.New("no uris")
	// ErrNoDialect db驱动不支持
	ErrNoDialect = errors.New("no dialect")
	// ErrNoAvailableHost 所有uri都不能连接
	ErrNoAvailableHost = errors.New("no available host")
)
