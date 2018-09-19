package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVerify(t *testing.T) {
	Convey("verify", t, func() {
		c := &OrmEngineConfig{
			Dialect: "x",
			Uris:    []string{""},
		}
		So(c.Verify(), ShouldEqual, ErrNoDialect)

		c = &OrmEngineConfig{
			Dialect: "mysql",
			Uris:    []string{},
		}
		So(c.Verify(), ShouldEqual, ErrNoUris)

		c = &OrmEngineConfig{
			Dialect: "mysql",
			Uris:    []string{""},
		}
		So(c.Verify(), ShouldBeNil)

		c = &OrmEngineConfig{
			Dialect: "sqlite3",
			Uris:    []string{""},
		}
		So(c.Verify(), ShouldBeNil)

		c = &OrmEngineConfig{
			Dialect: "postgres",
			Uris:    []string{""},
		}
		So(c.Verify(), ShouldBeNil)
	})
}
