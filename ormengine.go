package hornet

import (
	"errors"

	"gitee.com/JMArch/hornet/config"
	"gitee.com/JMArch/hornet/gengine"
	"gitee.com/JMArch/hornet/xengine"
)

// NewGorm NewGorm
func newGorm(cfg config.OrmEngineConfig) (*gengine.Orm, error) {
	return gengine.NewOrm(cfg)
}

// NewXorm NewXorm
func newXorm(cfg config.OrmEngineConfig) (*xengine.Orm, error) {
	return xengine.NewOrm(cfg)
}

// NewOrm 创建一个Orm
// t: Orm类型，有gorm、xorm两种
// cfg: orm的配置
func NewOrm(t string, cfg config.OrmEngineConfig) (interface{}, error) {
	switch t {
	case "gorm":
		return newGorm(cfg)
	case "xorm":
		return newXorm(cfg)
	}
	return nil, errors.New("type must be in [\"gorm\", \"xorm\"]")
}
