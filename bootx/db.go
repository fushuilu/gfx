package bootx

import (
	"fmt"

	"github.com/fushuilu/gfx/cx"
	"github.com/fushuilu/golibrary/appx/db"
	"github.com/fushuilu/golibrary/lerror"
	"github.com/fushuilu/golibrary/libx/cachex"

	"xorm.io/xorm"
)

/*
[db]
    name = "postgres"
    master = "postgres://用户名:密码@localhost/库?sslmode=disable"
    slaver = "postgres://用户名:密码@localhost/库?sslmode=disable"
    maxIdleConn = 0
    maxOpenConn = 10
*/
func LoadDB(name string) (eg *xorm.EngineGroup, err error) {
	dc := db.Config{}
	if err = cx.GetStruct(name, &dc); err != nil {
		return eg, lerror.Wrap(err, fmt.Sprintf("could not find %s in config", name))
	}
	eg = db.CreateEngineGroup(dc)
	return
}

/*
[xredis]
    host = "127.0.0.1:6379"
    password = ""
    database = 10
*/
func LoadRedis(name string) (redis *cachex.XRedis, err error) {
	redisOpts := cachex.RedisOpts{}
	if err := cx.GetStruct(name, &redisOpts); err != nil {
		return nil, lerror.Wrap(err, fmt.Sprintf("could not find %s in config file", name))
	}
	return cachex.NewXRedis(&redisOpts), nil
}
