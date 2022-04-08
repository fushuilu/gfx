package httpx

import (
	"errors"
	"fmt"
	"github.com/fushuilu/gfx/logx"
	"github.com/fushuilu/golibrary"
	"github.com/fushuilu/golibrary/lerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

// 本地请求

type LocalClient struct {
	cf   LocalClientConfig
	http golibrary.ClientHttp
}

type LocalClientConfig struct {
	Name   string `json:"name"`
	Token  string `json:"token"`
	Debug  bool   `json:"debug"`
	Origin string `json:"origin"`
}

func NewLocalClient(cf LocalClientConfig) LocalClient {
	return LocalClient{
		cf:   cf,
		http: NewClientHttpWith(cf.Debug),
	}
}

//
//  Request
//  @Description: 发送请求
//  @param name string 当前请求应用的名称
//  @param path string 接口的路径
//  @param data interface{} 请求的数据
//  @param resp interface{} 返回的数据
//  @return error 错误的信息
//
func (c *LocalClient) Request(path string, data interface{}, resp interface{}) error {
	sign, err := golibrary.MD5(fmt.Sprintf("%s:%s:%s", c.cf.Name, path, c.cf.Token))
	//logx.Debug("Request 路径:", path, ";name:", c.cf.Name, ";token:", c.cf.Token)
	if err != nil {
		return lerror.Wrap(err, "生成请求签名错误")
	}


	url := fmt.Sprintf("%s?name=%s&sign=%s", golibrary.HttpURLContact(c.cf.Origin, path), c.cf.Name, sign)
	return c.http.Post(url, data, resp)
}

type LocalServer struct {
	cf LocalServerConfig
}

type LocalServerConfig struct {
	Token string   `json:"token"` // 默认
	Ips   []string `json:"ips"`   // 允许请求的 ip，如果为空，则为 127.0.0.1
}

func NewLocalServer(cf LocalServerConfig) LocalServer {
	if len(cf.Ips) == 0 {
		cf.Ips = []string{"127.0.0.1"}
	}
	return LocalServer{cf: cf}
}

//
//  Receive
//  @Description: 接收请求
//
func (c *LocalServer) Receive(r *ghttp.Request, pd interface{}) error {
	ip := IpClient(r)
	if !golibrary.IsInString(c.cf.Ips, ip) {
		return errors.New("IP 地址不在白名单中")
	}
	sign := r.Get("sign").String()
	if sign == "" {
		return errors.New("接口签名不能为空")
	}
	name := r.Get("name").String()
	encrypt, err := golibrary.MD5(fmt.Sprintf("%s:%s:%s", name, r.Request.URL.Path, c.cf.Token))
	if err != nil {
		return errors.New("生成验证签名错误")
	}
	if encrypt != sign {
		logx.Debug("路径:", r.Request.URL.Path, ";name:", name, ";token:", c.cf.Token)
		logx.Debug("参数 sign:", sign, ";计算结果:", encrypt)
		return errors.New("签名不匹配")
	}

	Request(r, pd)
	return nil
}
