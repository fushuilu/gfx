package httpx

import (
	"errors"
	"fmt"

	"github.com/fushuilu/gfx/logx"

	"github.com/fushuilu/golibrary"
	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/fushuilu/golibrary/lerror"

	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type RouteRole string

const (
	RoleIgnore       RouteRole = "ignore"      // 忽略，不做权限检测
	RoleMustNotLogin           = "not-login"   // 必须未登录
	RoleUser                   = "user"        // 必须登录
	RoleUserTry                = "user-try"    // 尝试登录
	RoleAdmin                  = "admin"       // 管理员
	RoleSuperAdmin             = "super-admin" // 超级管理员
)

type IdCard struct {
	UserId int64  `json:"user_id"` // 用户身份（主键）
	Uid    string `json:"uid"`     // 用户身份（字符串）
	Name   string `json:"name"`    // 应用名称
}

// C: controller 控制器
// R: role 用户角色 0 不需要，1 需要
// I: interface 返回值 0 不需要，1 需要
type (
	CR0I1HandlerFunc func(r *ghttp.Request) interface{}
	CR1I0HandlerFunc func(r *ghttp.Request, card IdCard)
	CR1I1HandlerFunc func(r *ghttp.Request, card IdCard) interface{}
)

const (
	HttpMethodGet = "GET"
	HttpMethodPos = "POST"
	HttpMethodPut = "PUT"
	HttpMethodDel = "DELETE"
)

// (无分组)权限路由
type XRoute struct {
	S        *ghttp.Server
	Response func(r *ghttp.Request, body interface{})

	IdCardRepo     UserIdCard      // 检查用户是否登录
	PermissionRepo RoutePermission // 默认权限判断

	check bool
}

// 检查是否设置正确
func (x *XRoute) invalid() error {
	if x.PermissionRepo == nil {
		return errors.New("route PermissionRepo is nil")
	} else if x.IdCardRepo == nil {
		return errors.New("route IdCaredRepo is nil")
	}
	x.check = true
	return nil
}

// 获取当前用户的身份，并进行校验

func (x *XRoute) GetIdCard(r *ghttp.Request, role RouteRole) (idCard IdCard) {
	if role == RoleIgnore {
		return
	}

	if !x.check {
		golibrary.PanicIfError(x.invalid())
	}
	var err error
	idCard, err = x.IdCardRepo.GetIdCard(r.GetHeader(datax.Authorization))
	logx.DebugIfError(err, "route getIdCard err")

	switch role {
	case RoleUserTry:
		return
	case RoleMustNotLogin:
		if idCard.UserId > 0 {
			x.Response(r, lerror.NewCode(403, "只允许访客访问"))
		}
	case RoleUser:
		if idCard.UserId < 1 {
			x.Response(r, lerror.NewCode(401, "您还没有登录"))
		}
	default:
		if idCard.UserId < 1 {
			x.Response(r, lerror.NewCode(500, "账户未登录"))
		}
		if !x.PermissionRepo.Access(idCard, role) {
			x.Response(r, lerror.NewCode(403, fmt.Sprintf("没有访问角色 %s 的权限", role)))
		}
	}
	return
}

func NewXRoute() *XRoute {
	return &XRoute{
		S:        g.Server(),
		Response: SmartResponse,
	}
}

func (x *XRoute) BindHandler(method, pattern string, handler ghttp.HandlerFunc) {
	if method == "" {
		panic("空的请求方式:" + pattern)
	}
	x.S.BindHandler(method+":"+pattern, handler)
}

func (x *XRoute) C01Get(pattern string, handler CR0I1HandlerFunc) {
	x.C01(HttpMethodGet, pattern, handler)
}

func (x *XRoute) C01Pos(pattern string, handler CR0I1HandlerFunc) {
	x.C01(HttpMethodPos, pattern, handler)
}

func (x *XRoute) C01Put(pattern string, handler CR0I1HandlerFunc) {
	x.C01(HttpMethodPut, pattern, handler)
}

func (x *XRoute) C01Del(pattern string, handler CR0I1HandlerFunc) {
	x.C01(HttpMethodDel, pattern, handler)
}

func (x *XRoute) C01(method, pattern string, handler CR0I1HandlerFunc) {
	logRoute(pattern, method, gdebug.FuncPath(handler))
	x.BindHandler(method, pattern, func(r *ghttp.Request) {
		x.Response(r, handler(r))
	})
}

func (x *XRoute) C11Get(pattern string, role RouteRole, handler CR1I1HandlerFunc) {
	x.C11(HttpMethodGet, pattern, role, handler)
}
func (x *XRoute) C11Pos(pattern string, role RouteRole, handler CR1I1HandlerFunc) {
	x.C11(HttpMethodPos, pattern, role, handler)
}
func (x *XRoute) C11Put(pattern string, role RouteRole, handler CR1I1HandlerFunc) {
	x.C11(HttpMethodPut, pattern, role, handler)
}
func (x *XRoute) C11Del(pattern string, role RouteRole, handler CR1I1HandlerFunc) {
	x.C11(HttpMethodDel, pattern, role, handler)
}

func (x *XRoute) C11(method, pattern string, role RouteRole, handler CR1I1HandlerFunc) {
	logRoute(pattern, method, gdebug.FuncPath(handler))
	x.BindHandler(method, pattern, func(r *ghttp.Request) {
		x.Response(r, handler(r, x.GetIdCard(r, role)))
	})
}

func (x *XRoute) C10Get(pattern string, role RouteRole, handler CR1I0HandlerFunc) {
	x.c10(HttpMethodGet, pattern, role, handler)
}
func (x *XRoute) C10Pos(pattern string, role RouteRole, handler CR1I0HandlerFunc) {
	x.c10(HttpMethodPos, pattern, role, handler)
}
func (x *XRoute) C10Put(pattern string, role RouteRole, handler CR1I0HandlerFunc) {
	x.c10(HttpMethodPut, pattern, role, handler)
}
func (x *XRoute) C10Del(pattern string, role RouteRole, handler CR1I0HandlerFunc) {
	x.c10(HttpMethodDel, pattern, role, handler)
}

func (x *XRoute) c10(method, pattern string, role RouteRole, handler CR1I0HandlerFunc) {
	logRoute(pattern, method, gdebug.FuncPath(handler))
	x.BindHandler(method, pattern, func(r *ghttp.Request) {
		handler(r, x.GetIdCard(r, role))
	})
}

// 路由组
type XGroupRoute struct {
	prefix string    // 路由前缀
	role   RouteRole // 最小访问角色
	xRoute *XRoute   // 权限路由
}

func NewXGroupRoute(prefix string, role RouteRole, route *XRoute, handle func(r XGroupRoute)) {
	xr := XGroupRoute{
		prefix: prefix,
		role:   role,
		xRoute: route,
	}
	handle(xr)
}

// 简单模式，没有权限判断，也不需要设置路由
func NewGroupRoute(prefix string, handle func(r XGroupRoute)) {
	NewXGroupRoute(prefix, RoleIgnore, NewXRoute(), handle)
}

func (xr *XGroupRoute) Get(pattern string, handler ghttp.HandlerFunc) {
	xr.C00(HttpMethodGet, pattern, handler)
}
func (xr *XGroupRoute) Pos(pattern string, handler ghttp.HandlerFunc) {
	xr.C00(HttpMethodPos, pattern, handler)
}
func (xr *XGroupRoute) Put(pattern string, handler ghttp.HandlerFunc) {
	xr.C00(HttpMethodPut, pattern, handler)
}
func (xr *XGroupRoute) Del(pattern string, handler ghttp.HandlerFunc) {
	xr.C00(HttpMethodDel, pattern, handler)
}

// 原生保持一致
func (xr *XGroupRoute) C00(method, pattern string, handler ghttp.HandlerFunc) {
	logRoute(xr.prefix+pattern, method, gdebug.FuncPath(handler))
	xr.xRoute.BindHandler(method, xr.prefix+pattern, func(r *ghttp.Request) {
		xr.xRoute.GetIdCard(r, xr.role)
		handler(r)
	})
}

// 使用路径，不使用前缀
func (xr *XGroupRoute) C01PathGet(path string, handler CR0I1HandlerFunc) {
	xr.C01(HttpMethodGet, path, handler)
}
func (xr *XGroupRoute) C01PathPos(path string, handler CR0I1HandlerFunc) {
	xr.C01(HttpMethodPos, path, handler)
}

func (xr *XGroupRoute) C01Get(pattern string, handler CR0I1HandlerFunc) {
	xr.C01(HttpMethodGet, pattern, handler)
}
func (xr *XGroupRoute) C01Pos(pattern string, handler CR0I1HandlerFunc) {
	xr.C01(HttpMethodPos, pattern, handler)
}
func (xr *XGroupRoute) C01Put(pattern string, handler CR0I1HandlerFunc) {
	xr.C01(HttpMethodPut, pattern, handler)
}
func (xr *XGroupRoute) C01Del(pattern string, handler CR0I1HandlerFunc) {
	xr.C01(HttpMethodDel, pattern, handler)
}

func (xr *XGroupRoute) C01(method, pattern string, handler CR0I1HandlerFunc) {
	logRoute(xr.prefix+pattern, method, gdebug.FuncPath(handler))

	xr.xRoute.BindHandler(method, xr.prefix+pattern, func(r *ghttp.Request) {
		xr.xRoute.GetIdCard(r, xr.role)
		xr.xRoute.Response(r, handler(r))
	})
}

func (xr *XGroupRoute) C10Get(pattern string, handler CR1I0HandlerFunc) {
	xr.xRoute.C10Get(xr.prefix+pattern, xr.role, handler)
}
func (xr *XGroupRoute) C10Pos(pattern string, handler CR1I0HandlerFunc) {
	xr.xRoute.C10Pos(xr.prefix+pattern, xr.role, handler)
}
func (xr *XGroupRoute) C10Put(pattern string, handler CR1I0HandlerFunc) {
	xr.xRoute.C10Put(xr.prefix+pattern, xr.role, handler)
}
func (xr *XGroupRoute) C10Del(pattern string, handler CR1I0HandlerFunc) {
	xr.xRoute.C10Del(xr.prefix+pattern, xr.role, handler)
}

func (xr *XGroupRoute) C11(method string, pattern string, handler CR1I1HandlerFunc) {
	switch method {
	case HttpMethodGet:
		xr.C11Get(pattern, handler)
	case HttpMethodPos:
		xr.C11Pos(pattern, handler)
	case HttpMethodPut:
		xr.C11Put(pattern, handler)
	case HttpMethodDel:
		xr.C11Del(pattern, handler)
	default:
		panic("unknown http method:" + method)
	}
}

func (xr *XGroupRoute) C11Get(pattern string, handler CR1I1HandlerFunc) {
	xr.xRoute.C11Get(xr.prefix+pattern, xr.role, handler)
}
func (xr *XGroupRoute) C11Pos(pattern string, handler CR1I1HandlerFunc) {
	xr.xRoute.C11Pos(xr.prefix+pattern, xr.role, handler)
}
func (xr *XGroupRoute) C11Put(pattern string, handler CR1I1HandlerFunc) {
	xr.xRoute.C11Put(xr.prefix+pattern, xr.role, handler)
}
func (xr *XGroupRoute) C11Del(pattern string, handler CR1I1HandlerFunc) {
	xr.xRoute.C11Del(xr.prefix+pattern, xr.role, handler)
}

type routeInfo struct {
	path   string
	method string
	name   string
}

// 打印路由
var routes = make([]routeInfo, 0)

func logRoute(path, method, name string) {
	routes = append(routes, routeInfo{path: path, method: method, name: name})
}

// 注意：使用原生绑定的路由不会出现在这里（包含退出登录等……）
func PrintRoutes() {
	for _, v := range routes {
		logx.Debugf("%-50s%-6s %s", v.path, v.method, v.name)
	}
	ClearRoutes()
}

func ClearRoutes() {
	routes = make([]routeInfo, 0)
}
