package httpx

import (
	"github.com/fushuilu/gfx/bootx"
	"github.com/fushuilu/golibrary/libx/cachex"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"
)

/*
config.toml
------------------------------
[db]
    name = "postgres"
    master = "postgres://用户名:密码@localhost/库?sslmode=disable"
    slaver = "postgres://用户名:密码@localhost/库?sslmode=disable"
    maxIdleConn = 0
    maxOpenConn = 10

[redis]
    default = "127.0.0.1:6379,10"
[xredis]
    host = "127.0.0.1:6379"
    password = ""
    database = 10
*/

type App struct {
	Eg    *xorm.EngineGroup
	Cache *caches.LRUCacher
	Redis *cachex.XRedis
	Data  map[string]interface{} // 其它追加的数据
}

func InitApp() (*App, error) {
	var (
		app App
		err error
	)
	if app.Eg, err = bootx.LoadDB("db"); err != nil {
		return nil, err
	}
	if app.Redis, err = bootx.LoadRedis("xredis"); err != nil {
		return nil, err
	}
	app.Cache = caches.NewLRUCacher(caches.NewMemoryStore(), 1000)
	app.Data = map[string]interface{}{}

	return &app, nil
}

func (rc *App) GetData(name string) (interface{}, bool) {
	if e, ok := rc.Data[name]; ok {
		return e, true
	}
	return nil, false
}

func (rc *App) SetData(name string, d interface{}) {
	rc.Data[name] = d
}
