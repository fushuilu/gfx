package logic

import (
	"github.com/fushuilu/gfx/httpx"
	"github.com/fushuilu/gfx/sharex/model"

	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/fushuilu/golibrary/appx/db"
	"github.com/fushuilu/golibrary/appx/errorx"

	"xorm.io/xorm"
)

type UserPerm struct {
	ma db.ModelAction
}

func NewUserPerm(eg *xorm.EngineGroup) UserPerm {
	return UserPerm{
		ma: model.NewUserPermModelAction(eg),
	}
}

func (c *UserPerm) Get(id int64) (*model.UserPerm, error) {
	mo := model.UserPerm{}
	err := c.ma.MustGetById(id, &mo)
	return &mo, err
}

func (c *UserPerm) Create(mo *model.UserPerm) error {
	if err := mo.Invalid(true); err != nil {
		return err
	}
	if mo.RoleName == httpx.RoleSuperAdmin {
		return errorx.New("禁止添加超级管理员")
	}

	// 因为当前操作只有管理员和超级管理员能够访问，所以当前用户最低权限也是管理员了
	cond := c.ma.NewWhere().String("gname", mo.Gname).Int64("user_id", mo.UserId).
		String("role_name", mo.RoleName).
		Finish()
	if err := c.ma.MustNotExist(cond); err != nil {
		return err
	}

	_, err := c.ma.Eg().InsertOne(mo)
	return c.ma.InsertRst(mo.Id, err)
}

func (c *UserPerm) Update(mo *model.UserPerm) error {
	if err := mo.Invalid(false); err != nil {
		return err
	}

	if mo.RoleName == httpx.RoleSuperAdmin {
		return errorx.New("禁止修改超级管理员")
	}

	num, err := c.ma.Eg().ID(mo.Id).Cols("start_at", "end_at", "status_index",
		"user_id", "account",
		"real_name", "role_name", "remark").
		Update(&mo)
	return c.ma.UpdateRst(num, err)
}

func (c *UserPerm) Delete(id int64) error {
	mo := model.UserPerm{}
	err := c.ma.GetById(id, &mo)
	if err = c.ma.GetRst(mo.Id, err); err != nil {
		return err
	}

	if mo.RoleName == httpx.RoleSuperAdmin {
		return errorx.New("不可以移除超级管理员")
	}

	return c.ma.DeleteById(id)
}

type UserPermSearchReq struct {
	RealName string `json:"real_name"`
	RoleName string `json:"role_name"`
	Gname    string `json:"gname"`
}

func (c *UserPerm) Search(pd UserPermSearchReq, pag db.Pagination) (*datax.ListResult, error) {
	cond := c.ma.NewWhere().String("gname", pd.Gname).
		Like("real_name", pd.RealName).
		String("role_name", pd.RoleName).Finish()

	var rows []model.UserPerm
	return c.ma.ListResult(cond, pag, &rows)
}
