package model

import (
	"errors"
	"time"

	"github.com/fushuilu/golibrary/appx"
	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/fushuilu/golibrary/appx/db"

	"xorm.io/xorm"
)

type Page struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"-" xorm:"created"`
	UpdatedAt time.Time `json:"-" xorm:"updated"`
	DeletedAt time.Time `json:"-" xorm:"deleted"`
	UserId    int64     `json:"user_id" xorm:"user_id notnull default(0)"` // 用户

	Gname   string `json:"gname" xorm:"gname notnull default('')"`    // 应用组名称
	SortId  int    `json:"sort_id" xorm:"sort_id notnull default(0)"` // 排序 ID
	Tag     int64  `json:"tag" xorm:"tag notnull default(0)"`         // 绑定的标识
	Title   string `json:"title" xorm:"title notnull default('')"`    // 中文标题
	Name    string `json:"name" xorm:"notnull default('')"`           // 英文名称
	Content string `json:"content" xorm:"text"`                       // 单页内容

	StatusIndex int    `json:"-" xorm:"int notnull default(1)"` // 状态
	Status      string `json:"status" xorm:"-"`
}

func (mo *Page) RecordFormat() {
	mo.Status = datax.MapStatusText.ToText(mo.StatusIndex)
}

func (mo *Page) IsGet() bool {
	return mo.Id > 0
}

func (mo *Page) Invalid(isCreated bool) error {
	if mo.Title == "" {
		return errors.New("标题不能为空")
	}

	datax.MapStatusText.GetValue(mo.Status, func(i int) {
		mo.StatusIndex = i
	})

	return appx.InvalidId(isCreated, mo.Id)
}

func NewPageModelAction(eg *xorm.EngineGroup) db.ModelAction {
	return db.NewModelAction(eg, "单页", &Page{})
}

func (Page) TableName() string {
	return "share_page"
}
