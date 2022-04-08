package sdk

import "fmt"

const (
	departmentListURL = "https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=%s&id=%d"
)

// id 否，部门 ID
// https://work.weixin.qq.com/api/doc/90000/90135/90208
func (w *WorkProgram) GetDepartmentList(id int64) (rows []RespDepartmentItem, err error) {

	var resp struct {
		Department []RespDepartmentItem `json:"department"`
	}
	err = w.api("getDepartmentList", func(token string) string {
		return fmt.Sprintf(departmentListURL, token, id)
	}, &resp)
	return resp.Department, err
}

type RespDepartmentItem struct {
	Id       int64  `json:"id"`       // 创建的部门id
	Name     string `json:"name"`     // 部门名称 可能没有
	NameEn   string `json:"name_en"`  // 英文名称 可能没有
	Parentid int64  `json:"parentid"` // 父部门 ID，根部门为1
	Order    int    `json:"order"`
}
