package gengine

import (
	"fmt"
	"sync"
	"time"

	"gitee.com/JMArch/hornet/checker"
	"gitee.com/JMArch/hornet/config"
	"github.com/jinzhu/gorm"
)

// GormEngine GormEngine
type GormEngine struct {
	*gorm.DB
	sync.RWMutex
	dbs checker.Checkers

	log               *logger
	loggerMode        *bool
	blockGlobalUpdate *bool
	singularTable     *bool

	stop   chan bool
	stoped bool
}

// StopEngine
func (ge *GormEngine) DestroyEngine() {
	ge.Lock()
	defer ge.Unlock()
	if ge.stoped {
		return
	}
	ge.stoped = true
	close(ge.stop)
	return
}

func (ge *GormEngine) repair() {
	if ge.log != nil {
		ge.DB.SetLogger(*ge.log)
	}
	if ge.loggerMode != nil {
		ge.DB.LogMode(*ge.loggerMode)
	}
	if ge.blockGlobalUpdate != nil {
		ge.DB.BlockGlobalUpdate(*ge.blockGlobalUpdate)
	}
	if ge.singularTable != nil {
		ge.DB.SingularTable(*ge.singularTable)
	}
}

type logger interface {
	Print(v ...interface{})
}

func (ge *GormEngine) check() {
	t := time.NewTicker(time.Millisecond * 10)
	for {
		select {
		case <-t.C:
			ge.dbs.CheckAndReplace(func(newOne int) bool {
				// 有替换发生
				if newOne >= 0 && newOne < len(ge.dbs) {
					newDB := ge.dbs[newOne]
					if db, ok := newDB.(*GormEngineChecker); ok {
						fmt.Println("replace", newOne, db.uri)
						ge.DB = db.DB
						ge.repair()
						return true
					}
				}
				return false
			})
		case <-ge.stop:
			ge.DB.Close()
			for _, d := range ge.dbs {
				if db, ok := d.(*GormEngineChecker); ok {
					db.Close()
				}
			}
			return
		}
	}
	t.Stop()
}

// SetLogger 设置logger（拦截logger设置）
func (ge *GormEngine) SetLogger(log logger) {
	ge.log = &log
	ge.DB.SetLogger(log)
}

func (ge *GormEngine) LogMode(enable bool) *gorm.DB {
	ge.loggerMode = &enable
	return ge.DB.LogMode(enable)
}

func (ge *GormEngine) BlockGlobalUpdate(enable bool) *gorm.DB {
	ge.blockGlobalUpdate = &enable
	return ge.DB.BlockGlobalUpdate(enable)
}

func (ge *GormEngine) SingularTable(enable bool) {
	ge.singularTable = &enable
	ge.DB.SingularTable(enable)
}

func NewGormEngine(cfg config.OrmEngineConfig) (*GormEngine, error) {
	if e := cfg.Verify(); e != nil {
		return nil, e
	}

	var master *gorm.DB
	masterIndex := -1
	var err error
	dbs := make(checker.Checkers, len(cfg.Uris))
	for index, uri := range cfg.Uris {
		db, e := gorm.Open(cfg.Dialect, uri)
		dbs[index] = &GormEngineChecker{
			DB:     db,
			uri:    uri,
			diaect: cfg.Dialect,
		}
		if masterIndex == -1 {
			if e != nil {
				err = e
			} else {
				e := dbs[index].Ping()
				if e != nil {
					err = e
				} else {
					masterIndex = index
					master = db
				}
			}
		}
	}
	if masterIndex == -1 {
		return nil, fmt.Errorf("%v. [%v]", config.NoAvailableHost, err)
	}

	// master保证在0位
	dbs[0], dbs[masterIndex] = dbs[masterIndex], dbs[0]
	out := &GormEngine{
		DB:      master,
		dbs:     dbs,
		RWMutex: sync.RWMutex{},
		stop:    make(chan bool, 1),
	}

	go out.check()

	return out, nil
}
