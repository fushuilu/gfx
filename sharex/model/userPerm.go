package model

import (
	"errors"
	"time"

	"github.com/fushuilu/golibrary"
	"github.com/fushuilu/golibrary/appx"
	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/fushuilu/golibrary/appx/db"

	"xorm.io/xorm"
)

// 用户权限
type UserPerm struct {
	Id      int64     `json:"id"`
	Created time.Time `json:"created" xorm:"created"` // 创建时间
	Updated time.Time `json:"updated" xorm:"updated"` // 更新时间
	Deleted time.Time `json:"deleted" xorm:"deleted"` // 删除时间

	Gname       string    `json:"gname" xorm:"gname notnull default('')"` // 组名
	UserId      int64     `json:"user_id" xorm:"user_id notnull default(0)"`
	StatusIndex int       `json:"-" xorm:"status_index int notnull default(0)"`
	Account     string    `json:"account" xorm:"account notnull default('')"`
	RealName    string    `json:"real_name" xorm:"real_name notnull default('')"`
	RoleName    string    `json:"role_name" xorm:"role_name notnull default('')"`
	Remark      string    `json:"remark" xorm:"remark notnull default('')"`
	StartAt     time.Time `json:"start_at" xorm:"start_at"`
	EndAt       time.Time `json:"end_at" xorm:"end_at"`

	Status string `json:"status" xorm:"-"`
	Start  string `json:"start" xorm:"-"`
	End    string `json:"end" xorm:"-"`
}

func (mo *UserPerm) IsGet() bool {
	return mo.Id > 0
}

func NewUserPermModelAction(eg *xorm.EngineGroup) db.ModelAction {
	return db.NewModelAction(eg, "用户角色", &UserPerm{})
}

func (mo *UserPerm) Invalid(isCreated bool) error {
	if mo.RoleName == "" {
		return errors.New("必须指定角色")
	}
	if mo.RealName == "" {
		return errors.New("必须填写真实姓名")
	}
	if mo.Account == "" {
		return errors.New("必须填写用户账号")
	} else if mo.UserId < 1 {
		return errors.New("必须指定用户 userId")
	}
	mo.StatusIndex = datax.MapStatusText.GetIntValue(mo.Status)

	var err error
	if mo.StartAt, err = golibrary.DateTimeParse(mo.Start); err != nil {
		return err
	}
	if mo.EndAt, err = golibrary.DateTimeParse(mo.End); err != nil {
		return err
	}
	return appx.InvalidId(isCreated, mo.Id)
}

func (mo *UserPerm) RecordFormat() {
	mo.Start = golibrary.HumanTimeYMDHM(mo.StartAt)
	mo.End = golibrary.HumanTimeYMDHM(mo.EndAt)
	mo.Status = datax.MapStatusText.ToText(mo.StatusIndex)
}

func (UserPerm) TableName() string {
	return "share_userperm"
}
