package xengine

import (
	"github.com/go-xorm/xorm"
)

// XormEngineChecker XormEngineChecker
type XormEngineChecker struct {
	*xorm.Engine
	uri    string
	diaect string
}

// Ping Ping
func (xec *XormEngineChecker) Ping() (e error) {
	if xec != nil {
		return xec.Engine.Ping()
	}
	return nil
}

// ReConnect ReConnect
func (xec *XormEngineChecker) ReConnect() (e error) {
	tmp, e := xorm.NewEngine(xec.diaect, xec.uri)
	if e != nil {
		return e
	} else if e = tmp.Ping(); e != nil {
		return e
	}
	xec.Close()
	xec.Engine = tmp
	return nil
}
