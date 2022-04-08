package httpx

import (
	"xorm.io/xorm"

	"github.com/gogf/gf/v2/net/ghttp"
)

// 在控制器中使用自定义的身份模型
type Merchant struct {
	Id int64 `json:"id"`
}

// 控制器方法
type MchHandler func(r *ghttp.Request, mch *Merchant) interface{}

type MchRoute struct {
	eg    *xorm.EngineGroup
	route *XRoute
}

func NewRouteMock(route *XRoute, eg *xorm.EngineGroup) MchRoute {
	return MchRoute{
		eg:    eg,
		route: route,
	}
}

func (a *MchRoute) GetMch(userId int64) (mch Merchant, err error) {
	// get mch from user id
	return
}

func (a *MchRoute) Do(method, pattern string, handler MchHandler) {
	a.route.C11(method, "/api/mch/"+pattern, RoleUser, func(r *ghttp.Request, card IdCard) interface{} {
		mch, err := a.GetMch(card.UserId)
		ResponseIfError(r, err)
		return handler(r, &mch)
	})
}
