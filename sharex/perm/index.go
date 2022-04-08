package perm

import (
	"time"

	"github.com/fushuilu/gfx/httpx"
	"github.com/fushuilu/gfx/sharex/action"
	"github.com/fushuilu/golibrary"

	"xorm.io/xorm"
)

func New(eg *xorm.EngineGroup, gName string) httpx.RoutePermission {
	return &basePermission{action: action.NewUserPermAction(eg, gName)}
}

type basePermission struct {
	action action.UserPermAction
}

func (a *basePermission) Access(card httpx.IdCard, miniRole httpx.RouteRole) (ok bool) {
	if roles, err := a.action.GetRoles(card.UserId, time.Now()); err != nil {
		return false
	} else if golibrary.IsInString(roles, httpx.RoleSuperAdmin) {
		return true
	} else {
		return golibrary.IsInString(roles, string(miniRole))
	}
}