package sdk

import (
	"fmt"
)

const (
	getUserDetailURL            = "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s"
	departmentUserSimpleListURL = "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist?access_token=%s&department_id=%d&fetch_child=%d"
	departmentUserDetailListURL = "https://qyapi.weixin.qq.com/cgi-bin/user/list?access_token=%s&department_id=%d&fetch_child=%d"
)

/*
读取企业成员的详情信息
https://work.weixin.qq.com/api/doc/90000/90135/90196
userid 成员 UserID。对应管理端的帐号，企业内必须唯一。
[{"id":1,"name":"东莞市博与科技有限公司","name_en":"","parentid":0,"order":100000000},
{"id":2,"name":"蓝海","name_en":"","parentid":1,"order":100000000}]
*/
func (w *WorkProgram) GetUserDetail(userid string) (info RespWorkUserDetail, err error) {
	err = w.api("GetUserDetail", func(token string) string {
		return fmt.Sprintf(getUserDetailURL, token, userid)
	}, &info)
	return
}

type RespWorkUserDetail struct {
	UserId         string  `json:"UserId"`            // 成员UserID。对应管理端的帐号，企业内必须唯一。不区分大小写，长度为1~64个字节
	Name           string  `json:"name"`              // 2020年6月30日起，对所有历史第三方应用不再返回真实name，使用userid代替name
	Department     []int64 `json:"department"`        // 成员所属部门 ID 列表
	Order          []int   `json:"order"`             // 部门内排序，默认为0，数量必须和department一致，数值越大排序越前面。值范围是[0, 2^32)
	Position       string  `json:"position"`          // 职务信息
	Mobile         string  `json:"mobile"`            // 手机号码，第三方仅通讯录应用可获取
	Gender         string  `json:"gender"`            // 性别。0表示未定义，1表示男性，2表示女性
	Email          string  `json:"email"`             // 邮箱，第三方仅通讯录应用可获取
	IsLeaderInDept []int   `json:"is_leader_in_dept"` // 表示在所在的部门内是否为上级。；第三方仅通讯录应用可获取
	Avatar         string  `json:"avatar"`            // 头像url。 第三方仅通讯录应用可获取
	ThumbAvatar    string  `json:"thumb_avatar"`      // 头像缩略图url。第三方仅通讯录应用可获取
	Telephone      string  `json:"telephone"`
	Alias          string  `json:"alias"`
	Address        string  `json:"address"`
	OpenUserid     string  `json:"open_userid"` // 全局唯一。对于同一个服务商，不同应用获取到企业内同一个成员的open_userid是相同的
	MailDepartment int     `json:"mail_department"`
	Extattr        struct {
		Attrs []RespExternalAttr `json:"attrs"`
	} `json:"extattr"`
	Status           int    `json:"status"` // 激活状态: 1=已激活，2=已禁用，4=未激活，5=退出企业。
	QrCode           string `json:"qr_code"`
	ExternalPosition string `json:"external_position"`
	ExternalProfile  struct {
		ExternalCorpName string             `json:"external_corp_name"` // 企业简称
		ExternalAttr     []RespExternalAttr `json:"external_attr"`
	} `json:"external_profile"`
}

type RespExternalAttr struct {
	Type int      `json:"type"`
	Name string   `json:"name"`
	Text struct { // type = 0
		Value string `json:"value"`
	} `json:"text"`
	Web struct { // type = 1
		Url   string `json:"url"`
		Title string `json:"title"`
	} `json:"web"`
	Miniprogram struct {
		Appid    string `json:"appid"`
		Pagepath string `json:"pagepath"`
		Title    string `json:"title"`
	} `json:"miniprogram"`
}

/*
https://work.weixin.qq.com/api/doc/90000/90135/90200
获取部门成员
[{"userid":"LvShuTao","name":"吕树涛","department":[1],"open_userid":""},
{"userid":"CeShiXiaoHao","name":"测试小号","department":[1],"open_userid":""}]
*/
func (w *WorkProgram) GetDepartmentUserSimpleList(departmentId int64) (rows []RespDepartmentUserSimpleItem, err error) {
	var resp struct {
		Userlist []RespDepartmentUserSimpleItem `json:"userlist"`
	}
	err = w.api("GetDepartmentUserList", func(token string) string {
		return fmt.Sprintf(departmentUserSimpleListURL, token, departmentId, 0)
	}, &resp)
	return resp.Userlist, err
}

type RespDepartmentUserSimpleItem struct {
	Userid     string  `json:"userid"`      // 成员 ID
	Name       string  `json:"name"`        // 成员名称，可使用 userid 代替 name
	Department []int64 `json:"department"`  // 成员所属部门列表
	OpenUserid string  `json:"open_userid"` // 全局唯一
}

/*
https://work.weixin.qq.com/api/doc/90000/90135/90201
获取部门成员详情列表

[{"UserId":"LvShuTao","name":"吕树涛","department":[1],"order":[0],
	"position":"",
	"mobile":"13400010001",
	"gender":"1",
	"email":"",
	"is_leader_in_dept":[0],
	"avatar":"https://wework.qpic.cn/bizmail/Rh9ic2HUZBvfrWRkpnPsriavtaKP2u3SdibpRVuJPPffbHOpXN2N9g/0",
	"thumb_avatar":"https://wework.qpic.cn/bizmail/Rh9ic2HUZBvfmYFrWRkpnPsriavtaKP2u3SdibpRVuJPPffbHOpXN2N9g/100",
	"telephone":"",
	"alias":"吕生",
	"address":"",
	"open_userid":"",
	"mail_depart0,
	"extattr":{"attrs":[]},
	"status":1,
	"qr_code":"https://open.work.weixin.qq.com/wwopen/userQRCode?vcode=vc0b2e0274fa899670",
	"external_position":"",
	"external_profile":{"external_corp_name":"","external_attr":null}
}]
*/
func (w *WorkProgram) GetDepartmentUserDetailList(departmentId int64) (rows []RespWorkUserDetail, err error) {
	var resp struct {
		Userlist []RespWorkUserDetail `json:"userlist"`
	}
	err = w.api("GetDepartmentUserDetailList", func(token string) string {
		return fmt.Sprintf(departmentUserDetailListURL, token, departmentId, 0)
	}, &resp)
	return resp.Userlist, err
}
