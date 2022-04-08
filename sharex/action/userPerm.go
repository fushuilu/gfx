package action

import (
	"errors"
	"time"

	"github.com/fushuilu/gfx/sharex/model"

	"github.com/fushuilu/golibrary"
	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/fushuilu/golibrary/appx/errorx"
	"github.com/fushuilu/golibrary/lerror"

	"xorm.io/xorm"
)

type UserPermAction struct {
	eg   *xorm.EngineGroup
	name string
}

func NewUserPermAction(eg *xorm.EngineGroup, gName string) UserPermAction {
	return UserPermAction{eg: eg, name: gName}
}

// 用户的应用组权限，不检查 name=''
func (a *UserPermAction) AccessRoles(userId int64, roles []string) (ok bool, err error) {
	if userId < 1 {
		return false, errorx.New("用户 ID 不能为空")
	}

	now := golibrary.HumanTimeYMDHMS(time.Now())
	ok, err = a.eg.
		Where("gname=? AND status_index=? AND user_id=?", a.name, datax.IndexStatusActive, userId).
		And(`(start_at IS NULL AND end_at IS NULL) 
OR (start_at IS NULL AND end_at >?) 
OR (start_at <? AND end_at IS NULL)
OR (start_at <? AND end_at >?)`, now, now, now, now).
		In("role_name", roles).Exist(&model.UserPerm{})
	if err != nil {
		err = lerror.Wrap(err, "查询用户权限错误")

	}
	return
}

// 用户的权限
// name='' 时表示这是一个超级管理员
func (a *UserPermAction) GetRoles(userId int64, time time.Time) (roles []string, err error) {
	if userId < 1 {
		return roles, errors.New("用户 ID 不能为空")
	}

	now := golibrary.HumanTimeYMDHMS(time)
	var rows []model.UserPerm
	if err = a.eg.Cols("role_name").
		Where("(gname='' OR gname=?) AND status_index=? AND user_id=?", a.name, datax.IndexStatusActive, userId).
		And(`(start_at IS NULL AND end_at IS NULL) 
OR (start_at IS NULL AND end_at >?) 
OR (start_at <? AND end_at IS NULL)
OR (start_at <? AND end_at >?)`, now, now, now, now).
		Find(&rows); err != nil {
		return
	}
	roles = make([]string, len(rows))
	for i := range rows {
		roles[i] = rows[i].RoleName
	}
	return
}
