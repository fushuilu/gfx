package wex

import (
	"errors"
	"github.com/fushuilu/gfx/wex/sdk"
	"github.com/fushuilu/golibrary/appx/errorx"
	"github.com/fushuilu/golibrary/lerror"
	"github.com/fushuilu/golibrary/libx/cachex"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/auth"
	"github.com/silenceper/wechat/v2/miniprogram/encryptor"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"

	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
)

// 公众号
func GetGzhInstance(c WechatConfig, cache *cachex.XRedis) (*officialaccount.OfficialAccount, error) {

	if c.AppID == "" {
		return nil, errors.New("公众号配置 appid 为空")
	} else if c.Token == "" {
		return nil, errors.New("公众号配置 token 为空")
	} else if c.AppSecret == "" {
		return nil, errors.New("公众号配置 appsecret 为空")
	}

	wc := wechat.NewWechat()
	return wc.GetOfficialAccount(&offConfig.Config{
		AppID:          c.AppID,
		AppSecret:      c.AppSecret,
		Token:          c.Token,
		EncodingAESKey: c.AesKey,
		Cache:          cache,
	}), nil
}

// 小程序
func GetMiniInstance(c WechatConfig, cache *cachex.XRedis) (*miniprogram.MiniProgram, error) {
	if c.AppID == "" {
		return nil, errors.New("小程序 appid 为空")
	} else if c.AppSecret == "" {
		return nil, errors.New("小程序 appsecret 为空")
	}

	wc := wechat.NewWechat()
	return wc.GetMiniProgram(&miniConfig.Config{
		AppID:     c.AppID,
		AppSecret: c.AppSecret,
		Cache:     cache,
	}), nil
}

// 企业微信
func GetWorkInstance(c WechatConfig, cache *cachex.XRedis) (*sdk.WorkProgram, error) {
	if c.AppID == "" {
		return nil, errors.New("企业微信 AgentId|AppID 为空")
	} else if c.CropId == "" {
		return nil, errors.New("企业微信 CropId 为空")
	} else if c.AppSecret == "" {
		return nil, errors.New("企业微信 appsecret 为空")
	}

	if program, err := sdk.NewWorkProgram(&sdk.WorkConfig{
		Corpid:  c.CropId,
		AgentId: c.AppID,
		Secret:  c.AppSecret,
		Cache:   cache,
	}); err != nil {
		return nil, lerror.Wrap(err, "初始化企业微信错误")
	} else {
		return &program, nil
	}
}

func GetWorkCpt(c *WechatConfig) (wk *sdk.WxCpt, err error) {
	if c.AppID == "" {
		return nil, errors.New("企业微信 cpt appid 为空")
	} else if c.CropId == "" {
		return nil, errors.New("企业微信 cpt CropId 为空")
	}
	cpt := sdk.NewWxCpt(c.Token, c.AesKey, c.CropId)
	return &cpt, nil
}

// 小程序
type Mini interface {
	Code2Session(mini *miniprogram.MiniProgram, code string) (rst auth.ResCode2Session, err error)
	Decrypt(mini *miniprogram.MiniProgram, sessionKey, encryptedData, iv string) (info *encryptor.PlainData, err error)
}

type mini struct {
}

func NewMini() Mini {
	return &mini{}
}

func (r *mini) Code2Session(mini *miniprogram.MiniProgram, code string) (rst auth.ResCode2Session, err error) {
	if code == "" {
		err = errorx.New("code should not empty")
		return
	}
	if rst, err = mini.GetAuth().Code2Session(code); err != nil {
		return rst, lerror.Wrap(err, "小程序 code 换取 session 错误")
	}
	return
}

func (r *mini) Decrypt(mini *miniprogram.MiniProgram, sessionKey, encryptedData, iv string) (info *encryptor.PlainData, err error) {
	info, err = mini.GetEncryptor().Decrypt(sessionKey, encryptedData, iv)
	if err != nil {
		err = lerror.Wrap(err, "解密用户信息失败")
	}
	return
}

//
type Web interface {
	GetRedirectURL(wc WechatConfig, redirectURI, scope, state string) (string, error)
}

func NewWeb(cache *cachex.XRedis) Web {
	return &web{
		cache: cache,
	}
}

type web struct {
	cache *cachex.XRedis
}

func (w *web) GetRedirectURL(wc WechatConfig, redirectURI, scope, state string) (string, error) {
	if wc.IsGzh() || wc.IsWeb() {
		wx, err := GetGzhInstance(wc, w.cache)
		if err != nil {
			return "", err
		}
		oa := wx.GetOauth()
		if wc.IsGzh() {
			return oa.GetRedirectURL(redirectURI, scope, state)
		} else if wc.IsWeb() {
			return oa.GetWebAppRedirectURL(redirectURI, scope, state)
		}
	} else if wc.IsWork() {
		wk, err := GetWorkInstance(wc, w.cache)
		if err != nil {
			return "", err
		}
		if "snsapi_base" == scope {
			return wk.GetAuthWebRedirectURL(redirectURI, state), nil
		} else {
			return wk.GetAuthQrConnectRedirectURL(redirectURI, state), nil
		}
	}
	return "", errorx.New("不支持的公众号类型，无法生成回调地址", wc)
}

type WxUserInfo interface {
	// 获取(公众号/网页应用)用户信息
	GetUserInfo(oa *officialaccount.OfficialAccount, accessToken, openid string) (info oauth.UserInfo, err error)
	// 获取(公众号/网页应用) 用户 AccessToken
	GetUserAccessToken(oa *officialaccount.OfficialAccount, code string) (rst oauth.ResAccessToken, err error)
	// 企业微信员工信息
	GetWorkerInfo(wk *sdk.WorkProgram, code string) (rst sdk.RespGetUserinfo, err error)
}

func NewWxUserInfo() WxUserInfo {
	return &userInfo{}
}

type userInfo struct {
}

func (r *userInfo) GetUserInfo(oa *officialaccount.OfficialAccount, accessToken, openid string) (info oauth.UserInfo, err error) {
	info, err = oa.GetOauth().GetUserInfo(accessToken, openid, "")
	if err != nil {
		err = lerror.Wrap(err, "获取公众号用户信息失败")
	}
	return
}

func (r *userInfo) GetUserAccessToken(oa *officialaccount.OfficialAccount, code string) (rst oauth.ResAccessToken, err error) {
	rst, err = oa.GetOauth().GetUserAccessToken(code)

	if err != nil {
		err = lerror.Wrap(err, "获取 UserAccessToken 信息失败")
	}
	return
}

func (r *userInfo) GetWorkerInfo(wk *sdk.WorkProgram, code string) (rst sdk.RespGetUserinfo, err error) {
	return wk.GetUserinfo(code)
}
