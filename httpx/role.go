package httpx

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// 当前接口的角色
type Role int

const (
	_     Role = iota
	Admin      // 管理员
	Mch        // 商户
	User       // 用户
	Open       // 公开

)

// 当前接口的配置
type ApiInfo struct {
	Name  string // 群名字（英文）
	Title string // 中文标题
	Role  Role
}

func AdminApi() ApiInfo {
	return ApiInfo{Role: Admin}
}
func MchApi() ApiInfo {
	return ApiInfo{Role: Mch}
}
func UserApi() ApiInfo {
	return ApiInfo{Role: User}
}
func OpenApi() ApiInfo {
	return ApiInfo{Role: Open}
}

func (c *ApiInfo) IsAdmin() bool {
	return c.Role == Admin
}
func (c *ApiInfo) IsUser() bool {
	return c.Role == User
}
func (c *ApiInfo) IsOpen() bool {
	return c.Role == Open
}

func (c *ApiInfo) IsMch() bool {
	return c.Role == Mch
}

func (c *ApiInfo) MustLogin(r *ghttp.Request) {
	if c.IsOpen() {
		ResponseError(r, "非开放接口，不支持调用")
	}
	return
}

func (c *ApiInfo) MustAdmin(r *ghttp.Request) {
	if c.IsAdmin() {
		return
	}
	ResponseError(r, "admin 管理接口，没有权限调用")
}

func (c *ApiInfo) MustUser(r *ghttp.Request) {
	if c.IsUser() {
		return
	}
	ResponseError(r, "user 用户接口，没有权限调用")
}

func (c *ApiInfo) MustMch(r *ghttp.Request) {
	if c.IsMch() {
		return
	}
	ResponseError(r, "mch 商户接口，没有权限调用")
}

func (c *ApiInfo) MustRoles(r *ghttp.Request, roles ...Role) {
	for _, v := range roles {
		if c.Role == v {
			return
		}
	}
	ResponseError(r, "没有访问当前接口的权限")
}

func (c *ApiInfo) SetName(name string) {
	if name != "" {
		c.Name = name
	}
}
