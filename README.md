# hornet

## 项目介绍

提供了高可用支持的gorm&xorm

## 软件架构

```text
├── README.md
├── checker 抽象的checker接口，提供ping和reconnect
│   └── check.go
├── config 提供配置
│   ├── config.go
│   └── errors.go
├── gengine gorm的engine
│   ├── gorm.go
│   └── gormEngineChecker.go
├── ormengine.go 对外接口，创建engine
└── xengine xorm的engine
    ├── xorm.go
    └── xormEngineChecker.go
```

## 如何获取

`go get -u -v github.com/sunreaver/hornet`

## 使用说明

### 1. gorm

```golang
import (
	"github.com/sunreaver/hornet"
	"github.com/sunreaver/hornet/config"
	"github.com/sunreaver/hornet/gengine"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	// mysql
	dbTmp, e := hornet.NewOrm("gorm", config.OrmEngineConfig{
		Dialect: "mysql",
		Uris: []string{
			"root:123456@tcp(localhost:3307)/t_test?charset=utf8&parseTime=True&loc=Local",
			"root:123456@tcp(localhost:3308)/t_test?charset=utf8&parseTime=True&loc=Local",
		},
	})
	// sqlite3
	//dbTmp, e := hornet.NewOrm("gorm", config.OrmEngineConfig{
	//	Dialect: "sqlite3",
	//	Uris: []string{
	//		"./local.db",
	//	},
	//})
	if e != nil {
		panic(e.Error())
	}
	db := dbTmp.(*gengine.Orm)

	// Do Your Business logic
	// db.FirstOrCreate(...)


	// exit
	db.DestroyEngine()
)
```

### 2. xorm

```golang
import (
	"fmt"
	"time"

	"github.com/sunreaver/hornet"
	"github.com/sunreaver/hornet/config"
	"github.com/sunreaver/hornet/gengine"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
)

func main() {
	dbTmp, e := hornet.NewOrm("xorm", config.OrmEngineConfig{
		Dialect: "mysql",
		Uris: []string{
			"root:123456@tcp(localhost:3307)/t_test?charset=utf8&parseTime=True&loc=Local",
			"root:123456@tcp(localhost:3308)/t_test?charset=utf8&parseTime=True&loc=Local",
		},
	})
	if e != nil {
		panic(e.Error())
	}

	db := dbTmp.(*xengine.Orm)

	// Do Your Business logic
	// db.Insert(...)


	// exit
	db.DestroyEngine()
}
```

## 参与贡献

1. Fork 本项目
2. 新建 Feat_xxx 分支
3. 提交代码
4. 新建 Pull Request
