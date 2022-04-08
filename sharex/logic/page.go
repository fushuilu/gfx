package logic

import (
	"errors"

	"github.com/fushuilu/gfx/sharex/model"

	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/fushuilu/golibrary/appx/db"
	"github.com/fushuilu/golibrary/appx/dbx"

	"xorm.io/xorm"
)

type Page struct {
	xMa dbx.ModelAction
	ma  db.ModelAction
}

func NewPage(eg *xorm.EngineGroup) Page {
	ma := model.NewPageModelAction(eg)
	return Page{
		ma:  ma,
		xMa: dbx.NewModelAction(ma),
	}
}

type PageSearchReq struct {
	Gname   string `json:"gname"`
	Tag     int64  `json:"tag"`
	Title   string `json:"title"`
	Name    string `json:"name"`
	Keyword string `json:"keyword"`
	Status  string `json:"status"`
}

func (c *Page) Search(pd PageSearchReq, pag db.Pagination) (*datax.ListResult, error) {

	var (
		statusIndex int
	)
	datax.MapStatusText.GetValueIfNotEmpty(pd.Status, func(i int) {
		statusIndex = i
	})

	cond := c.ma.NewWhere().String("gname", pd.Gname).
		Int64("tag", pd.Tag).
		Like("title", pd.Title).
		String("name", pd.Name).
		Or(pd.Keyword, "title", "name").
		Int("status_index", statusIndex).Finish()

	var rows []model.Page
	return c.ma.ListResultWith(cond, pag, &rows, func(se *xorm.Session) {
		se.Desc("sort_id", "id")
	})
}

type PageGetReq struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Gname       string `json:"gname"`
	StatusIndex int    `json:"-"`
}

// 通常用于开放接口
func (c *Page) GetBy(pd PageGetReq) (*model.Page, error) {
	if pd.Id < 1 && pd.Name == "" {
		return nil, errors.New("必须指定查询的单页 id 或 名称")
	}

	cond := c.ma.NewWhere().Int64("id", pd.Id).
		String("name", pd.Name).
		String("gname", pd.Gname).
		Int("status_index", pd.StatusIndex).Finish()

	mo := model.Page{}
	if err := c.ma.MustGet(cond, &mo); err != nil {
		return nil, err
	}
	return &mo, nil
}

func (c *Page) Delete(id int64) error {
	return c.ma.DeleteById(id)
}

func (c *Page) Create(mo *model.Page) error {

	if err := mo.Invalid(true); err != nil {
		return err
	}
	cond := c.ma.NewWhere().String("name", mo.Name).String("gname", mo.Gname).Finish()
	if err := c.ma.MustNotExist(cond); err != nil {
		return err
	}

	_, err := c.ma.Eg().InsertOne(mo)
	return c.ma.InsertRst(mo.Id, err)
}

func (c *Page) Update(mo *model.Page) error {
	if err := mo.Invalid(false); err != nil {
		return err
	}
	cond := c.ma.NewWhere().String("name", mo.Name).
		String("gname", mo.Gname).NotEqual("id", mo.Id).Finish()

	if err := c.ma.MustNotExist(cond); err != nil {
		return err
	}

	num, err := c.ma.Eg().ID(mo.Id).Cols("sort_id", "tag", "title", "name", "content", "status_index").
		Update(mo)
	return c.ma.UpdateRst(num, err)
}

func (c *Page) ChangeStatus(data datax.ChangeStatus) error {
	return c.xMa.ChangeStatus(data)
}
