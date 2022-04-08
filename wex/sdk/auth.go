package sdk

import (
	"fmt"
	"net/url"
)

const (
	authQrRedirectURL  = "https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=%s&agentid=%s&redirect_uri=%s&state=%s"
	authWebRedirectURL = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect"
	getUserinfoURL     = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s"
)

// 构造扫码登录链接
// https://open.work.weixin.qq.com/api/doc/90000/90135/91019
func (w *WorkProgram) GetAuthQrConnectRedirectURL(urlStr, state string) string {
	urlStr = url.QueryEscape(urlStr)
	return fmt.Sprintf(authQrRedirectURL, w.conf.Corpid, w.conf.AgentId, urlStr, state)
}

// 构造网页授权链接
// https://open.work.weixin.qq.com/api/doc/90000/90135/91022
func (w *WorkProgram) GetAuthWebRedirectURL(urlStr, state string) string {
	urlStr = url.QueryEscape(urlStr)
	return fmt.Sprintf(authWebRedirectURL, w.conf.Corpid, urlStr, state)
}

// https://work.weixin.qq.com/api/doc/90000/90135/91437
// 用户通过授权之后，通过 code 获取访问用户 openid
func (w *WorkProgram) GetUserinfo(code string) (user RespGetUserinfo, err error) {
	err = w.api("GetUserinfo", func(token string) string {
		return fmt.Sprintf(getUserinfoURL, token, code)
	}, &user)
	return
}

type RespGetUserinfo struct {
	OpenId string `json:"OpenId"` // 非企业成员的标识，对当前企业唯一
	UserId string `json:"UserId"` // 企业成员 UserID。若需要获得用户详情信息，可调用通讯录接口：读取成员
}
