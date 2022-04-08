package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/fushuilu/golibrary/appx/errorx"
	"io/ioutil"
	"net/http"

	"github.com/fushuilu/gfx/logx"
	"github.com/fushuilu/golibrary"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

type clientHttp struct {
	isDebug bool
}

// 请求的域名
func NewClientHttp() golibrary.ClientHttp {
	return &clientHttp{}
}

func (ch *clientHttp) Debug(yes bool) {
	ch.isDebug = yes
}

func NewClientHttpWith(isDebug bool) golibrary.ClientHttp {
	return &clientHttp{
		isDebug: isDebug,
	}
}

func (ch *clientHttp) Post(url string, data interface{}, resp interface{}) error {
	ctx := gctx.New()
	rst, err := g.Client().Post(ctx, url, data)
	if err != nil {
		logx.Debug("g.http post error:", err)
		logx.Error("g.http post 失败", url, data)
		return errorx.New("http post 请求错误")
	}
	if ch.isDebug {
		rst.RawDump()
	}
	if rst.StatusCode != http.StatusOK {
		body := rst.ReadAllString()
		return errorx.New(body)
	}
	body := rst.ReadAll()

	//logx.Debug("body:",string(body))

	return json.Unmarshal(body, resp)
}

func (ch *clientHttp) Get(url string, params interface{}, resp interface{}) error {
	response, err := g.Client().Get(gctx.New(), url, params)
	if err != nil {
		logx.Debug("g.http get error:", err)
		logx.Error("g.http get 失败", url, params)
		return errorx.New("http GET 请求错误")
	}
	if ch.isDebug {
		response.RawDump()
	}
	if response.StatusCode != http.StatusOK {
		body := response.ReadAllString()
		return errors.New(body)
	}
	body := response.ReadAll()
	return json.Unmarshal(body, resp)
}

func (ch *clientHttp) PostByte(url string, data []byte, rst interface{}) error {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var body []byte
	if ch.isDebug {
		logx.Debug("http response:", string(body))
	}
	if body, err = ioutil.ReadAll(resp.Body); err != nil {

		logx.Debug("g.http PostByte error:", err)
		logx.Error("g.http PostByte 失败", url, data)
		return errorx.New("post 请求错误")
	}

	return json.Unmarshal(body, rst)
}
