package xengine

import (
	"strings"
	"sync"
	"testing"

	"gitee.com/JMArch/hornet/config"
	_ "github.com/go-sql-driver/mysql"
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
		Convey("ok", func() {
			orm, e := NewOrm(config.OrmEngineConfig{
				Dialect: "mysql",
				Uris:    mysqls,
			})
			So(e, ShouldBeNil)
			So(orm.DestroyEngine(), ShouldBeTrue)
		})
		Convey("fail dialect", func() {
			_, e := NewOrm(config.OrmEngineConfig{
				Dialect: "x",
				Uris:    mysqls,
			})
			So(e, ShouldEqual, config.ErrNoDialect)
		})
		Convey("fail uris", func() {
			_, e := NewOrm(config.OrmEngineConfig{
				Dialect: "mysql",
				Uris:    nil,
			})
			So(e, ShouldEqual, config.ErrNoUris)
		})
		Convey("can not connect", func() {
			_, e := NewOrm(config.OrmEngineConfig{
				Dialect: "mysql",
				Uris:    []string{"a"},
			})
			So(e, ShouldNotBeNil)
			So(strings.HasPrefix(e.Error(), config.ErrNoAvailableHost.Error()), ShouldBeTrue)
		})
	})
}

func TestDestroyEngine(t *testing.T) {
	Convey("destroy", t, func() {
		Convey("ok one goroutine", func() {
			orm, e := NewOrm(config.OrmEngineConfig{
				Dialect: "mysql",
				Uris:    mysqls,
			})
			So(e, ShouldBeNil)
			So(orm.DestroyEngine(), ShouldBeTrue)
			So(orm.DestroyEngine(), ShouldBeFalse)
		})
		Convey("ok more goroutine", func() {
			orm, e := NewOrm(config.OrmEngineConfig{
				Dialect: "mysql",
				Uris:    mysqls,
			})
			So(e, ShouldBeNil)

			var wg sync.WaitGroup
			okCount := 0
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func() {
					ok := orm.DestroyEngine()
					if ok {
						okCount++
					}
					wg.Done()
				}()
			}
			wg.Wait()
			So(orm.DestroyEngine(), ShouldBeFalse)
			So(okCount, ShouldEqual, 1)
		})
	})
}
