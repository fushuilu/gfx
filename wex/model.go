package wex

import (
	"errors"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
	"github.com/silenceper/wechat/v2/officialaccount/user"
)

type WechatConfig struct {
	AppID     string `json:"app_id,omitempty"`     // 开发者 ID, 企业应用的 agentId
	AppSecret string `json:"app_secret,omitempty"` // 开发者密钥
	Token     string `json:"token,omitempty"`      // 令牌
	AesKey    string `json:"aes_key,omitempty"`    // 消息加密密钥
	Url       string `json:"url,omitempty"`        // 服务器地址
	EncMethod string `json:"enc_method,omitempty"` // 加密方式
	Domain    string `json:"domain,omitempty"`     // 服务器域名，如 http://demo.fushuilu.com，用于接收微信回传的 code
	Callback  string `json:"callback,omitempty"`   // 回调地址，用于将授权信息携带到业务应用的授权页面地址，通常只有 web 应用才需要
	Kind      string `json:"kind,omitempty"`       // 类型
	CropId    string `json:"crop_id,omitempty"`    // 企业微信 CropId (相当于 appid)
}

func (c *WechatConfig) Invalid() error {
	if len(c.AppID) < 3 {
		return errors.New("app id 不能为空")
	}
	if len(c.AppSecret) < 1 {
		return errors.New("app secret 不能为空")
	}
	return nil
}

func (c *WechatConfig) IsMini() bool {
	return c.Kind == "mini"
}
func (c *WechatConfig) IsGzh() bool {
	return c.Kind == "gzh"
}
func (c *WechatConfig) IsWeb() bool {
	return c.Kind == "web"
}
func (c *WechatConfig) IsWork() bool {
	return c.Kind == "work"
}

type AccountInfo struct {
	From       string `json:"from"`  // gzh | mini | web | work
	Appid      string `json:"appid"` // 微信应用 id
	OpenId     string `json:"open_id"`
	UnionId    string `json:"union_id"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	City       string `json:"city"`
	Province   string `json:"province"`
	Country    string `json:"country"`
	Language   string `json:"language"`
	Headimgurl string `json:"headimgurl"` // 头像
	Phone      string `json:"phone"`      // 电话号码
	Email      string `json:"email"`      // 邮箱地址
	Workname   string `json:"workname"`   // 企业微信用户名称
	// 补充的信息
	UserId int64  `json:"user_id"`
	Uid    string `json:"uid"`
}

func (ai *AccountInfo) FromUserInfo(info *user.Info, appid, kind string) {
	ai.From = kind
	ai.Appid = appid
	ai.OpenId = info.OpenID
	ai.UnionId = info.UnionID
	ai.Nickname = info.Nickname
	ai.Sex = gconv.Int(info.Sex)
	ai.Province = info.Province
	ai.City = info.City
	ai.Country = info.Country
	ai.Language = info.Language
	ai.Headimgurl = info.Headimgurl
}

func (ai *AccountInfo) FromAccessToken(token *oauth.ResAccessToken) {
	ai.OpenId = token.OpenID
	ai.UnionId = token.UnionID
}

func (ai *AccountInfo) FromOauthUserInfo(info *oauth.UserInfo) {
	ai.OpenId = info.OpenID
	ai.UnionId = info.Unionid
	ai.Nickname = info.Nickname
	ai.Sex = gconv.Int(info.Sex)
	ai.Province = info.Province
	ai.City = info.City
	ai.Country = info.Country
	ai.Headimgurl = info.HeadImgURL
}

type WorkAccountInfo struct {
	CropId     string `json:"crop_id"`
	AgentId    string `json:"agent_id"`
	WorkUserid string `json:"work_userid"`
	UserId     int64  `json:"user_id"`
	Uid        string `json:"uid"`
}

func (pd *WorkAccountInfo) Invalid() error {
	if pd.CropId == "" {
		return errors.New("企业微信 id 为空")
	} else if pd.AgentId == "" {
		return errors.New("企业微信应用 id 为空")
	} else if pd.WorkUserid == "" {
		return errors.New("企业微信用户 id 为空")
	}
	return nil
}
