package model

import (
	"errors"
	"github.com/fushuilu/golibrary/appx"
	"github.com/fushuilu/golibrary/appx/db"
	"xorm.io/xorm"
)

type Setting struct {
	Id    int64  `json:"id"`
	Gname string `json:"gname" xorm:"gname notnull default('')"` // 应用组名称

	Title   string `json:"title" xorm:"title notnull default('')"`        // 中文标题
	Name    string `json:"name" xorm:"notnull default('')"`               // 英文名称
	Tag     int    `json:"tag" xorm:"tag notnull default(0)"`             // 分组
	IsOpen  bool   `json:"is_open" xorm:"is_open notnull default(false)"` // 开放
	IsLock  bool   `json:"is_lock" xorm:"is_lock notnull default(false)"` // 锁定(不允许删除)
	Content string `json:"content" xorm:"content notnull default('')"`    // 内容
}

func (mo *Setting) RecordFormat() {

}

func (mo *Setting) IsGet() bool {
	return mo.Id > 0
}

func (mo *Setting) Invalid(isCreated bool) error {
	if mo.Name == "" {
		return errors.New("名称不能为空")
	}

	return appx.InvalidId(isCreated, mo.Id)
}

func NewSettingModelAction(eg *xorm.EngineGroup) db.ModelAction {
	return db.NewModelAction(eg, "设置", &Setting{})
}

func (Setting) TableName() string {
	return "share_setting"
}
