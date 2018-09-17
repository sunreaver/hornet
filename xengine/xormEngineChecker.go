package xengine

import (
	"fmt"

	"github.com/go-xorm/xorm"
)

// XormEngineChecker XormEngineChecker
type XormEngineChecker struct {
	*xorm.Engine
	uri    string
	diaect string
}

// Info Info
func (xec *XormEngineChecker) Info() string {
	return fmt.Sprintf("uri:%s, diaect:%s", xec.uri, xec.diaect)
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
