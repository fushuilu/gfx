package httpx

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/fushuilu/golibrary/appx/datax"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"

	_ "github.com/mattn/go-sqlite3"
)

var cacheClient = map[int]*gclient.Client{}

// 测试专用客户端
func NewTestClient(conf ...interface{}) *gclient.Client {
	var port int
	if len(conf) > 0 {
		port = gconv.Int(conf[0])
	}
	if port < 1 {
		port = grand.N(10000, 11000) // 随机端口
	}
	if cli, ok := cacheClient[port]; ok {
		return cli
	}
	var client *gclient.Client
	s := g.Server()
	s.SetPort(port)
	s.SetDumpRouterMap(true)
	s.SetErrorLogEnabled(true)
	_ = s.Start()
	defer s.Shutdown()
	//time.Sleep(200 * time.Millisecond)

	client = gclient.New()
	client.SetBrowserMode(true)
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", port))
	cacheClient[port] = client
	return client
}

/*
func TestHandleGet(t *testing.T) {
  // 为被测试的处理器创建相应的请求
  request, _ := http.NewRequest("GET", "/post/1", nil)
  // 向被测试的处理器发送请求
  mux.ServeHTTP(writer, request)
  // 验证结果
  if writer.Code != 200 {
    t.Errorf("Response code is %v", writer.Code)
  }
  var p post
  json.Unmarshal(writer.Body.Bytes(), &p)
  if p.Id != 1 {
    t.Error("Cannot retrieve JSON post")
  }
}

func TestHandlePut(t *testing.T) {
  json := strings.NewReader(`{"content":"Updated post", "author":"Sau"}`)
  request, _ := http.NewRequest("PUT", "/post/1", json)
  mux.ServeHTTP(writer, request)

  if writer.Code != 200 {
    t.Errorf("Response code is %v", writer.Code)
  }
}
*/

// 制造一个 post 请求，通常用于控制器的测试
func NewTestPostRequest(postData interface{}, conf ...interface{}) *ghttp.Request {
	data, err := gjson.Encode(postData)
	gtest.Assert(err, nil)
	u := "/"
	if len(conf) > 0 {
		u = gconv.String(conf[0])
	}
	request, err := http.NewRequest("POST", u, strings.NewReader(string(data)))
	gtest.Assert(err, nil)
	return &ghttp.Request{
		Request: request,
	}
}

// 制造一个 get 请求，通常用于 get 请求
func NewTestGetRequest(path string) *ghttp.Request {
	request2, err := http.NewRequest("GET", path, nil)
	gtest.Assert(err, nil)
	return &ghttp.Request{
		Request: request2,
	}
}

func NewTestPostFormRequest(pd url.Values) ghttp.Request {
	r, _ := http.NewRequest("POST", "/project", strings.NewReader(pd.Encode()))
	r.Header.Add("Content-Type", datax.MIMEApplicationForm)
	r.Header.Add("Content-Length", strconv.Itoa(len(pd.Encode())))

	return ghttp.Request{
		Request: r,
	}
}

func NewTestResponse(req *http.Request, body string) http.Response {
	return http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       req,
		Header:        make(http.Header, 0),
	}
}
