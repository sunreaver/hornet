package gengine

import "github.com/jinzhu/gorm"

// GormEngineChecker GormEngineChecker
type GormEngineChecker struct {
	*gorm.DB
	uri    string
	diaect string
}

// Ping Ping
func (gec *GormEngineChecker) Ping() error {
	if gec != nil {
		if db := gec.DB.DB(); db != nil {
			return db.Ping()
		}
	}
	return nil
}

// ReConnect ReConnect
func (gec *GormEngineChecker) ReConnect() error {
	tmp, e := gorm.Open(gec.diaect, gec.uri)
	if e != nil {
		return e
	}
	gec.DB.Close()
	gec.DB = tmp
	return nil
}
