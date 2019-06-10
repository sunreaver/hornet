package hornet

import (
	"testing"

	"github.com/sunreaver/hornet/config"
	"github.com/sunreaver/hornet/gengine"
	"github.com/sunreaver/hornet/xengine"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	mysqls = []string{
		"root:123456@tcp(localhost:3307)/t_test?charset=utf8&parseTime=True&loc=Local",
		"root:123456@tcp(localhost:3308)/t_test?charset=utf8&parseTime=True&loc=Local",
		"root:123456@tcp(localhost:3309)/t_test?charset=utf8&parseTime=True&loc=Local",
		"root:123456@tcp(localhost:3310)/t_test?charset=utf8&parseTime=True&loc=Local",
		"root:123456@tcp(localhost:3311)/t_test?charset=utf8&parseTime=True&loc=Local",
		"root:123456@tcp(localhost:3312)/t_test?charset=utf8&parseTime=True&loc=Local",
	}
)

func TestNewOrm(t *testing.T) {
	Convey("test NewOrm", t, func() {
		Convey("gorm", func() {
			_, e := NewOrm("gorm", config.OrmEngineConfig{})
			So(e, ShouldNotBeNil)
			_, e = NewOrm("gorm", config.OrmEngineConfig{
				Dialect: "mysql",
				Uris:    []string{},
			})
			So(e, ShouldNotBeNil)

			dbEngine, e := NewOrm("gorm", config.OrmEngineConfig{
				Dialect: "mysql",
				Uris:    mysqls,
			})
			So(e, ShouldBeNil)

			db, ok := dbEngine.(*gengine.Orm)
			So(ok, ShouldBeTrue)

			destroy := db.DestroyEngine()
			So(destroy, ShouldBeTrue)
			destroy = db.DestroyEngine()
			So(destroy, ShouldBeFalse)
		})
		Convey("xorm", func() {
			_, e := NewOrm("xorm", config.OrmEngineConfig{})
			So(e, ShouldNotBeNil)
			_, e = NewOrm("xorm", config.OrmEngineConfig{
				Dialect: "mysql",
				Uris:    []string{},
			})
			So(e, ShouldNotBeNil)

			dbEngine, e := NewOrm("xorm", config.OrmEngineConfig{
				Dialect: "mysql",
				Uris:    mysqls,
			})
			So(e, ShouldBeNil)

			db, ok := dbEngine.(*xengine.Orm)
			So(ok, ShouldBeTrue)

			db.SetLogLevel(core.LOG_WARNING)

			destroy := db.DestroyEngine()
			So(destroy, ShouldBeTrue)
			destroy = db.DestroyEngine()
			So(destroy, ShouldBeFalse)
		})
	})
}
