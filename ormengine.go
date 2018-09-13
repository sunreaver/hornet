package hornet

import (
	"errors"

	"gitee.com/JMArch/hornet/config"
	"gitee.com/JMArch/hornet/gengine"
	"gitee.com/JMArch/hornet/xengine"
)

// NewGorm NewGorm
func newGorm(cfg config.OrmEngineConfig) (*gengine.GormEngine, error) {
	return gengine.NewGormEngine(cfg)
}

// NewXorm NewXorm
func newXorm(cfg config.OrmEngineConfig) (*xengine.XormEngine, error) {
	return xengine.NewXormEngine(cfg)
}

func NewOrm(t string, cfg config.OrmEngineConfig) (interface{}, error) {
	switch cfg.Type {
	case "gorm":
		return newGorm(cfg)
	case "xorm":
		return newXorm(cfg)
	}
	return nil, errors.New("type must be in [\"gorm\", \"xorm\"]")
}
