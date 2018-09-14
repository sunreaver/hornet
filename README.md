# hornet

## 项目介绍
提供了高可用支持的gorm&xorm

## 软件架构

```
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

`go get -u -v gitee.com/JMArch/hornet`

## 使用说明


#### 1. gorm

```golang
import (
	"gitee.com/JMArch/hornet"
	"gitee.com/JMArch/hornet/config"
	"gitee.com/JMArch/hornet/gengine"
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
	db := dbTmp.(*gengine.GormEngine)

	// Do Your Business logic
	// db.FirstOrCreate(...)


	// exit
	db.DestroyEngine()
)
```

#### 2. xorm

```golang
import (
	"fmt"
	"time"

	"gitee.com/JMArch/hornet"
	"gitee.com/JMArch/hornet/config"
	"gitee.com/JMArch/hornet/gengine"
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

	db := dbTmp.(*gengine.GormEngine)

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


## 码云特技

1. 使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2. 码云官方博客 [blog.gitee.com](https://blog.gitee.com)
3. 你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解码云上的优秀开源项目
4. [GVP](https://gitee.com/gvp) 全称是码云最有价值开源项目，是码云综合评定出的优秀开源项目
5. 码云官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6. 码云封面人物是一档用来展示码云会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
