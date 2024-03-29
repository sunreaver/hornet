package xengine

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/sunreaver/hornet/checker"
	"github.com/sunreaver/hornet/config"
	"github.com/sunreaver/hornet/log"
)

// Orm Orm
type Orm struct {
	*xorm.Engine
	lock       sync.RWMutex
	dbs        checker.Checkers
	dbsRepair  []int64
	lastRepair int64

	logger             *core.ILogger
	level              *core.LogLevel
	disableGlobalCache *bool
	cachers            map[string]core.Cacher
	cacherLock         sync.RWMutex
	tableMapper        *core.IMapper
	columnMapper       *core.IMapper
	cacher             *core.Cacher
	maxOpenConns       *int
	maxIdleConns       *int
	tzLocation         **time.Location
	databaseTZ         **time.Location
	schema             *string
	showExecTime       []bool
	showSQL            []bool

	stop   chan bool
	stoped bool
}

// DestroyEngine DestroyEngine
func (xe *Orm) DestroyEngine() bool {
	xe.lock.Lock()
	defer xe.lock.Unlock()
	if xe.stoped {
		return false
	}
	xe.stoped = true
	close(xe.stop)
	return true
}

func (xe *Orm) repair() {
	log.Logf("repair new db")
	if xe.logger != nil {
		xe.Engine.SetLogger(*xe.logger)
	}
	if xe.level != nil {
		xe.Engine.SetLogLevel(*xe.level)
	}
	if xe.disableGlobalCache != nil {
		xe.Engine.SetDisableGlobalCache(*xe.disableGlobalCache)
	}
	for t, c := range xe.cachers {
		xe.Engine.SetCacher(t, c)
	}
	if xe.tableMapper != nil {
		xe.Engine.SetTableMapper(*xe.tableMapper)
	}
	if xe.columnMapper != nil {
		xe.Engine.SetColumnMapper(*xe.columnMapper)
	}
	if xe.cacher != nil {
		xe.Engine.SetDefaultCacher(*xe.cacher)
	}
	if xe.maxOpenConns != nil {
		xe.Engine.SetMaxOpenConns(*xe.maxOpenConns)
	}
	if xe.maxIdleConns != nil {
		xe.Engine.SetMaxIdleConns(*xe.maxIdleConns)
	}
	if xe.tzLocation != nil {
		xe.Engine.SetTZLocation(*xe.tzLocation)
	}
	if xe.databaseTZ != nil {
		xe.Engine.SetTZDatabase(*xe.databaseTZ)
	}
	if xe.schema != nil {
		xe.Engine.SetSchema(*xe.schema)
	}
	if xe.showExecTime != nil {
		xe.Engine.ShowExecTime(xe.showExecTime...)
	}
	if xe.showSQL != nil {
		xe.Engine.ShowSQL(xe.showSQL...)
	}
}

func (xe *Orm) check() {
	t := time.NewTicker(time.Second * 1)
CHECKING:
	for {
		select {
		case <-t.C:
			xe.dbs.CheckAndReplace(func(newOne int) bool {
				// 有替换发生
				if newOne >= 0 && newOne < len(xe.dbs) {
					newDB := xe.dbs[newOne]
					if db, ok := newDB.(*XormEngineChecker); ok {
						log.Logf("select new db", db.Info())
						xe.Engine = db.Engine
						if xe.dbsRepair[newOne] != xe.lastRepair {
							xe.repair()
							xe.dbsRepair[newOne] = xe.lastRepair
						}
						return true
					}
				}
				return false
			})
		case <-xe.stop:
			xe.Engine.Close()
			for _, d := range xe.dbs {
				if db, ok := d.(*XormEngineChecker); ok {
					err := db.Close()
					log.Logf("xengine close at", db.Info(), "err", err)
				}
			}
			break CHECKING
		}
	}
	t.Stop()
}

// SetLogger SetLogger
func (xe *Orm) SetLogger(logger core.ILogger) {
	xe.lastRepair = time.Now().UnixNano()
	xe.logger = &logger
	xe.Engine.SetLogger(logger)
}

// SetLogLevel SetLogLevel
func (xe *Orm) SetLogLevel(level core.LogLevel) {
	xe.lastRepair = time.Now().UnixNano()
	xe.level = &level
	xe.Engine.SetLogLevel(level)
}

// SetDisableGlobalCache disable global cache or not
func (xe *Orm) SetDisableGlobalCache(disable bool) {
	xe.lastRepair = time.Now().UnixNano()
	xe.disableGlobalCache = &disable
	xe.Engine.SetDisableGlobalCache(disable)
}

// SetCacher SetCacher
func (xe *Orm) SetCacher(tableName string, cacher core.Cacher) {
	xe.lastRepair = time.Now().UnixNano()
	xe.cacherLock.Lock()
	xe.cachers[tableName] = cacher
	xe.cacherLock.Unlock()
	xe.Engine.SetCacher(tableName, cacher)
}

// SetMapper set the name mapping rules
func (xe *Orm) SetMapper(mapper core.IMapper) {
	xe.lastRepair = time.Now().UnixNano()
	xe.tableMapper = &mapper
	xe.columnMapper = &mapper
	xe.Engine.SetMapper(mapper)
}

// SetTableMapper set the table name mapping rule
func (xe *Orm) SetTableMapper(mapper core.IMapper) {
	xe.lastRepair = time.Now().UnixNano()
	xe.tableMapper = &mapper
	xe.Engine.SetTableMapper(mapper)
}

// SetColumnMapper set the column name mapping rule
func (xe *Orm) SetColumnMapper(mapper core.IMapper) {
	xe.lastRepair = time.Now().UnixNano()
	xe.columnMapper = &mapper
	xe.Engine.SetColumnMapper(mapper)
}

// SetDefaultCacher set the default cacher. Xorm's default not enable cacher.
func (xe *Orm) SetDefaultCacher(cacher core.Cacher) {
	xe.lastRepair = time.Now().UnixNano()
	xe.cacher = &cacher
	xe.Engine.SetDefaultCacher(cacher)
}

// SetMaxOpenConns is only available for go 1.2+
func (xe *Orm) SetMaxOpenConns(conns int) {
	xe.lastRepair = time.Now().UnixNano()
	xe.maxOpenConns = &conns
	xe.Engine.SetMaxOpenConns(conns)
}

// SetMaxIdleConns set the max idle connections on pool, default is 2
func (xe *Orm) SetMaxIdleConns(conns int) {
	xe.lastRepair = time.Now().UnixNano()
	xe.maxIdleConns = &conns
	xe.Engine.SetMaxIdleConns(conns)
}

// MapCacher MapCacher
func (xe *Orm) MapCacher(bean interface{}, cacher core.Cacher) error {
	xe.lastRepair = time.Now().UnixNano()
	xe.SetCacher(xe.Engine.TableName(bean, true), cacher)
	return nil
}

// SetTZLocation sets time zone of the application
func (xe *Orm) SetTZLocation(tz *time.Location) {
	xe.lastRepair = time.Now().UnixNano()
	xe.tzLocation = &tz
	xe.Engine.SetTZLocation(tz)
}

// SetTZDatabase sets time zone of the database
func (xe *Orm) SetTZDatabase(tz *time.Location) {
	xe.lastRepair = time.Now().UnixNano()
	xe.databaseTZ = &tz
	xe.Engine.SetTZDatabase(tz)
}

// SetSchema sets the schema of database
func (xe *Orm) SetSchema(schema string) {
	xe.lastRepair = time.Now().UnixNano()
	xe.schema = &schema
	xe.Engine.SetSchema(schema)
}

// ShowExecTime show SQL statement and execute time or not on logger if log level is great than INFO
func (xe *Orm) ShowExecTime(show ...bool) {
	xe.lastRepair = time.Now().UnixNano()
	xe.showExecTime = show
	xe.Engine.ShowExecTime(show...)
}

// ShowSQL show SQL statement or not on logger if log level is great than INFO
func (xe *Orm) ShowSQL(show ...bool) {
	xe.lastRepair = time.Now().UnixNano()
	xe.showSQL = show
	xe.Engine.ShowSQL(show...)
}

// NewOrm NewOrm
func NewOrm(cfg config.OrmEngineConfig) (*Orm, error) {
	if e := cfg.Verify(); e != nil {
		return nil, e
	}

	var master *xorm.Engine
	masterIndex := -1
	var err error
	dbs := make(checker.Checkers, len(cfg.Uris))
	for index, uri := range cfg.Uris {
		db, e := xorm.NewEngine(cfg.Dialect, uri)
		dbs[index] = &XormEngineChecker{
			Engine: db,
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
		Engine:     master,
		lock:       sync.RWMutex{},
		dbs:        dbs,
		dbsRepair:  make([]int64, len(dbs)),
		stop:       make(chan bool, 1),
		cachers:    map[string]core.Cacher{},
		cacherLock: sync.RWMutex{},
	}

	go out.check()

	return out, nil
}
