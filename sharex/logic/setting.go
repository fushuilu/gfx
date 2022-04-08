package logic

import (
	"errors"

	"github.com/fushuilu/gfx/sharex/model"

	"github.com/fushuilu/golibrary"
	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/fushuilu/golibrary/appx/db"

	"xorm.io/xorm"
)

type Setting struct {
	ma db.ModelAction
}

func NewSetting(eg *xorm.EngineGroup) Setting {

	return Setting{
		ma: model.NewSettingModelAction(eg),
	}
}

type SettingSearchReq struct {
	Title  string `json:"title"`
	IsOpen string `json:"is_open"`
	IsLock string `json:"is_lock"`
	Tag    int    `json:"tag"`
	Tags   string `json:"tags"`
	Gname  string `json:"gname"`
}

func (c *Setting) Search(pd SettingSearchReq, pag db.Pagination) (*datax.ListResult, error) {
	where := c.ma.NewWhere().String("gname", pd.Gname).
		Like("title", pd.Title).
		MockBoolean("is_open", pd.IsOpen).
		MockBoolean("is_lock", pd.IsLock).Int("tag", pd.Tag)
	if pd.Tags != "" {
		if tags, err := golibrary.SplitToInts(pd.Tags); err != nil {
			return nil, err
		} else if len(tags) > 0 {
			where.In("tag", tags)
		}
	}
	var rows []model.Setting

	return c.ma.ListResult(where.Finish(), pag, &rows)
}

type SettingGroupListReq struct {
	Tag   int    `json:"tag"`
	Gname string `json:"gname"`
	Open  string `json:"open"`
}

func (c *Setting) GroupList(pd SettingGroupListReq) (map[string]string, error) {

	cond := c.ma.NewWhere().String("gname", pd.Gname).
		Int("tag", pd.Tag).MockBoolean("is_open", pd.Open).Finish()

	var rows []model.Setting

	if err := c.ma.ListWith(cond, &rows, func(se *xorm.Session) {
		se.Cols("name", "content")
	}); err != nil {
		return nil, err
	}

	dict := make(map[string]string, len(rows))
	for i := range rows {
		dict[rows[i].Name] = rows[i].Content
	}
	return dict, nil
}

type SettingGetReq struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Gname string `json:"gname"`
	Open  string `json:"open"`
}

func (c *Setting) Get(pd SettingGetReq) (*model.Setting, error) {

	if pd.Id < 1 && pd.Name == "" {
		return nil, errors.New("必须指定 id 或 name")
	}

	where := c.ma.NewWhere().Int64("id", pd.Id).
		String("name", pd.Name).
		String("gname", pd.Gname).MockBoolean("is_open", pd.Open)

	mo := model.Setting{}
	if err := c.ma.MustGet(where.Finish(), &mo); err != nil {
		return nil, err
	}
	return &mo, nil
}

func (c *Setting) Delete(id int64) error {
	cond := c.ma.NewWhere().Int64("id", id).And("is_lock", false).Finish()
	return c.ma.DeleteWith(cond)
}

func (c *Setting) Create(mo *model.Setting) error {

	cond := c.ma.NewWhere().String("name", mo.Name).
		String("gname", mo.Gname).Finish()

	if err := c.ma.MustNotExist(cond); err != nil {
		return err
	}

	_, err := c.ma.Eg().InsertOne(mo)
	return c.ma.InsertRst(mo.Id, err)
}

func (c *Setting) Update(mo *model.Setting) error {

	cond := c.ma.NewWhere().String("name", mo.Name).
		String("gname", mo.Gname).NotEqual("id", mo.Id).
		Finish()
	if err := c.ma.MustNotExist(cond); err != nil {
		return err
	}

	num, err := c.ma.Eg().ID(mo.Id).Cols("title", "name", "tag", "is_open", "is_lock", "content").Update(&mo)
	return c.ma.UpdateRst(num, err)
}
