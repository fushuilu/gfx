package httpx


// 当前用户角色权限判断
type RoutePermission interface {
	// 是否大于给定的最小权限，会访问下面的 RoleInterface
	Access(card IdCard, miniRole RouteRole) (ok bool)
}

