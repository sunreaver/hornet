package gengine

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sunreaver/hornet/checker"
	"github.com/sunreaver/hornet/config"
	"github.com/sunreaver/hornet/log"
)

// Orm Orm
type Orm struct {
	*gorm.DB
	lock       sync.RWMutex
	dbs        checker.Checkers
	dbsRepair  []int64
	lastRepair int64

	log               *logger
	loggerMode        *bool
	blockGlobalUpdate *bool
	singularTable     *bool

	stop   chan bool
	stoped bool
}

// DestroyEngine DestroyEngine
func (ge *Orm) DestroyEngine() bool {
	ge.lock.Lock()
	defer ge.lock.Unlock()
	if ge.stoped {
		return false
	}
	ge.stoped = true
	close(ge.stop)
	return true
}

func (ge *Orm) repair() {
	log.Logf("repair new db")
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

func (ge *Orm) check() {
	t := time.NewTicker(time.Second)
CHECKING:
	for {
		select {
		case <-t.C:
			ge.dbs.CheckAndReplace(func(newOne int) bool {
				// 有替换发生
				if newOne >= 0 && newOne < len(ge.dbs) {
					newDB := ge.dbs[newOne]
					if db, ok := newDB.(*GormEngineChecker); ok {
						log.Logf("select new db", db.Info())
						ge.DB = db.DB
						if ge.dbsRepair[newOne] != ge.lastRepair {
							ge.repair()
							ge.dbsRepair[newOne] = ge.lastRepair
						}
						return true
					}
				}
				return false
			})
		case <-ge.stop:
			ge.DB.Close()
			for _, d := range ge.dbs {
				if db, ok := d.(*GormEngineChecker); ok {
					err := db.Close()
					log.Logf("gengine close at", db.Info(), "err", err)
				}
			}
			break CHECKING
		}
	}
	t.Stop()
}

// SetLogger 设置logger（拦截logger设置）
func (ge *Orm) SetLogger(log logger) {
	ge.lastRepair = time.Now().UnixNano()
	ge.log = &log
	ge.DB.SetLogger(log)
}

// LogMode LogMode
func (ge *Orm) LogMode(enable bool) *gorm.DB {
	ge.lastRepair = time.Now().UnixNano()
	ge.loggerMode = &enable
	return ge.DB.LogMode(enable)
}

// BlockGlobalUpdate BlockGlobalUpdate
func (ge *Orm) BlockGlobalUpdate(enable bool) *gorm.DB {
	ge.lastRepair = time.Now().UnixNano()
	ge.blockGlobalUpdate = &enable
	return ge.DB.BlockGlobalUpdate(enable)
}

// SingularTable SingularTable
func (ge *Orm) SingularTable(enable bool) {
	ge.lastRepair = time.Now().UnixNano()
	ge.singularTable = &enable
	ge.DB.SingularTable(enable)
}

// NewOrm NewOrm
func NewOrm(cfg config.OrmEngineConfig) (*Orm, error) {
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
		return nil, fmt.Errorf("%v. [%v]", config.ErrNoAvailableHost, err)
	}

	// master保证在0位
	dbs[0], dbs[masterIndex] = dbs[masterIndex], dbs[0]
	out := &Orm{
		lock:      sync.RWMutex{},
		DB:        master,
		dbs:       dbs,
		dbsRepair: make([]int64, len(dbs)),
		stop:      make(chan bool, 1),
	}

	go out.check()

	return out, nil
}
