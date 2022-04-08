package route

import (
	"errors"

	"github.com/fushuilu/gfx/cx"
	"github.com/fushuilu/gfx/httpx"
	"github.com/fushuilu/gfx/sharex/perm"

	"github.com/fushuilu/golibrary/lerror"
	"github.com/fushuilu/golibrary/libx/cachex"

	"xorm.io/xorm"
)

func New(gName string, redis *cachex.XRedis, eg *xorm.EngineGroup) (*httpx.XRoute, error) {

	jwtSecret, err := cx.GetString("site.jwt")
	if err != nil {
		return nil, lerror.Wrap(err, "get site.jwt token error")
	}
	if jwtSecret == "" {
		return nil, errors.New("site.jwt secret is empty")
	}

	gcRoute := httpx.NewXRoute()
	gcRoute.IdCardRepo = httpx.NewJwtUserIdCard(jwtSecret, redis)
	gcRoute.PermissionRepo = perm.New(eg, gName)

	return gcRoute, nil
}
