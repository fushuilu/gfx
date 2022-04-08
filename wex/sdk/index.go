package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/fushuilu/gfx/logx"
	"github.com/fushuilu/golibrary/appx/errorx"
	"github.com/fushuilu/golibrary/lerror"
	"github.com/fushuilu/golibrary/libx/cachex"

	"github.com/silenceper/wechat/v2/util"
)

const (
	// https://work.weixin.qq.com/api/doc/90000/90135/91039
	// 应用 access_token 缓存在 redis 中
	accessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
)

// https://work.weixin.qq.com/api/doc/90000/90135/90665#corpid
type WorkConfig struct {
	Corpid  string // 企业ID
	AgentId string // 应用 ID
	Secret  string // 应用密钥，每个应用都有自己的密钥
	Cache   *cachex.XRedis
}

type WorkProgram struct {
	conf *WorkConfig
}

func NewWorkProgram(conf *WorkConfig) (w WorkProgram, err error) {
	w = WorkProgram{conf: conf}

	if w.conf.Cache == nil {
		err = errorx.New("请设置企业微信缓存 redis")
		return
	}
	if w.conf.Corpid == "" {
		err = errorx.New("请设置企业微信 CropId")
		return
	}
	if w.conf.AgentId == "" {
		err = errorx.New("请提供企业微信应用 AgentId")
	}
	if w.conf.Secret == "" {
		err = errorx.New("请提供企业微信应用 Secret")
		return
	}
	_, err = w.GetAccessToken()
	return
}

// 获取 access_token 必须
func (w *WorkProgram) GetAccessToken() (token string, err error) {
	key := "work_" + w.conf.AgentId
	if w.conf.Cache.IsExist(key) {
		token, err = w.conf.Cache.GetString(key)
		if err != nil {
			return
		}
	} else {
		// 获取 accessToken
		urlStr := fmt.Sprintf(accessTokenURL, w.conf.Corpid, w.conf.Secret)
		var response []byte
		response, err = util.HTTPGet(urlStr)
		if err != nil {
			return
		}
		rst := workAccessTokenResponse{}
		err = json.Unmarshal(response, &rst)
		if err != nil {
			return
		}
		if rst.Errcode != 0 {
			err = fmt.Errorf("get work accessToken error : errcode=%v , errmsg=%v", rst.Errcode, rst.Errmsg)
			return
		}
		// 缓存
		if err = w.conf.Cache.SetString(key, rst.AccessToken, 7100); err != nil {
			err = lerror.Wrap(err, "设置企业微信 accessToken 错误")
			logx.Debug("work accessToken", rst.AccessToken)
			return
		}
		return rst.AccessToken, nil
	}
	return
}

type workAccessTokenResponse struct {
	Errcode     int    `json:"errcode"`      //	出错返回码，为0表示成功，非0表示调用失败
	Errmsg      string `json:"errmsg"`       // 返回码提示语
	AccessToken string `json:"access_token"` // 获取到的凭证，最长为512字节
	ExpiresIn   int64  `json:"expires_in"`   // 凭证的有效时间（秒）默认 7200
}

type RespWorkCommon struct {
	Errcode int    `json:"errcode"` // 出错返回码，为0表示成功，非0表示调用失败
	Errmsg  string `json:"errmsg"`  // 返回码提示语
}

func (w *WorkProgram) api(name string, url func(token string) string, v interface{}) (err error) {

	var accessToken string
	accessToken, err = w.GetAccessToken()
	if err != nil {
		return
	}
	reqURL := url(accessToken)
	var response []byte
	response, err = util.HTTPGet(reqURL)

	if err != nil {
		return
	}
	logx.Info("|<==== ", reqURL)
	logx.Info("work response:", string(response))

	rc := RespWorkCommon{}
	if err = json.Unmarshal(response, &rc); err != nil {
		return
	}

	if rc.Errcode != 0 {
		return fmt.Errorf("work (%s) error : errcode=%v , errmsg=%v",
			name, rc.Errcode, rc.Errmsg)
	}
	err = json.Unmarshal(response, v)
	return
}
