package xengine

import (
	"fmt"
	"sync"
	"time"

	"gitee.com/JMArch/hornet/checker"
	"gitee.com/JMArch/hornet/config"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

// XormEngine XormEngine
type XormEngine struct {
	*xorm.Engine
	sync.RWMutex
	dbs checker.Checkers

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

// StopEngine
func (xe *XormEngine) DestroyEngine() {
	if xe.stoped {
		return
	}
	close(xe.stop)
	return
}

func (xe *XormEngine) repair() {
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

func (xe *XormEngine) check() {
	t := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-t.C:
			xe.dbs.CheckAndReplace(func(newOne int) {
				// 有替换发生
				if newOne >= 0 && newOne < len(xe.dbs) {
					newDB := xe.dbs[newOne]
					if db, ok := newDB.(*XormEngineChecker); ok {
						fmt.Println("replace", newOne, db.uri)
						xe.Engine = db.Engine
						xe.repair()
					}
				}
			})
		case <-xe.stop:
			xe.Engine.Close()
			for _, d := range xe.dbs {
				if db, ok := d.(*XormEngineChecker); ok {
					db.Close()
				}
			}
			return
		}
	}
	t.Stop()
}

func (xe *XormEngine) SetLogger(logger core.ILogger) {
	xe.logger = &logger
	xe.Engine.SetLogger(logger)
}

func (xe *XormEngine) SetLogLevel(level core.LogLevel) {
	xe.level = &level
	xe.Engine.SetLogLevel(level)
}

// SetDisableGlobalCache disable global cache or not
func (xe *XormEngine) SetDisableGlobalCache(disable bool) {
	xe.disableGlobalCache = &disable
	xe.Engine.SetDisableGlobalCache(disable)
}

func (xe *XormEngine) SetCacher(tableName string, cacher core.Cacher) {
	xe.cacherLock.Lock()
	xe.cachers[tableName] = cacher
	xe.cacherLock.Unlock()
	xe.Engine.SetCacher(tableName, cacher)
}

// SetMapper set the name mapping rules
func (xe *XormEngine) SetMapper(mapper core.IMapper) {
	xe.tableMapper = &mapper
	xe.columnMapper = &mapper
	xe.Engine.SetMapper(mapper)
}

// SetTableMapper set the table name mapping rule
func (xe *XormEngine) SetTableMapper(mapper core.IMapper) {
	xe.tableMapper = &mapper
	xe.Engine.SetTableMapper(mapper)
}

// SetColumnMapper set the column name mapping rule
func (xe *XormEngine) SetColumnMapper(mapper core.IMapper) {
	xe.columnMapper = &mapper
	xe.Engine.SetColumnMapper(mapper)
}

// SetDefaultCacher set the default cacher. Xorm's default not enable cacher.
func (xe *XormEngine) SetDefaultCacher(cacher core.Cacher) {
	xe.cacher = &cacher
	xe.Engine.SetDefaultCacher(cacher)
}

func (xe *XormEngine) SetMaxOpenConns(conns int) {
	xe.maxOpenConns = &conns
	xe.Engine.SetMaxOpenConns(conns)
}

// SetMaxIdleConns set the max idle connections on pool, default is 2
func (xe *XormEngine) SetMaxIdleConns(conns int) {
	xe.maxIdleConns = &conns
	xe.Engine.SetMaxIdleConns(conns)
}

func (xe *XormEngine) MapCacher(bean interface{}, cacher core.Cacher) error {
	xe.SetCacher(xe.Engine.TableName(bean, true), cacher)
	return nil
}

// SetTZLocation sets time zone of the application
func (xe *XormEngine) SetTZLocation(tz *time.Location) {
	xe.tzLocation = &tz
	xe.Engine.SetTZLocation(tz)
}

// SetTZDatabase sets time zone of the database
func (xe *XormEngine) SetTZDatabase(tz *time.Location) {
	xe.databaseTZ = &tz
	xe.Engine.SetTZDatabase(tz)
}

// SetSchema sets the schema of database
func (xe *XormEngine) SetSchema(schema string) {
	xe.schema = &schema
	xe.Engine.SetSchema(schema)
}

// ShowExecTime show SQL statement and execute time or not on logger if log level is great than INFO
func (xe *XormEngine) ShowExecTime(show ...bool) {
	xe.showExecTime = show
	xe.Engine.ShowExecTime(show...)
}

// ShowSQL show SQL statement or not on logger if log level is great than INFO
func (xe *XormEngine) ShowSQL(show ...bool) {
	xe.showSQL = show
	xe.Engine.ShowSQL(show...)
}

func NewXormEngine(cfg config.OrmEngineConfig) (*XormEngine, error) {
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
		return nil, fmt.Errorf("%v. [%v]", config.NoAvailableHost, err)
	}

	// master保证在0位
	dbs[0], dbs[masterIndex] = dbs[masterIndex], dbs[0]
	out := &XormEngine{
		Engine:     master,
		dbs:        dbs,
		RWMutex:    sync.RWMutex{},
		stop:       make(chan bool, 1),
		cachers:    map[string]core.Cacher{},
		cacherLock: sync.RWMutex{},
	}

	go out.check()

	return out, nil
}
