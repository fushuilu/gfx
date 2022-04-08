package httpx

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fushuilu/gfx/logx"
	"github.com/fushuilu/golibrary"
	"github.com/fushuilu/golibrary/appx"
	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/fushuilu/golibrary/appx/errorx"
	"github.com/fushuilu/golibrary/lerror"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Ip 限制
func IpLimit(r *ghttp.Request, name string, maxNum, durationSeconds int) (err error) {
	ip := IpClient(r)
	if ip != "::1" && ip != "127.0.01" {
		now := int(time.Now().Unix())
		key := ip + ":" + name

		get, err := r.Session.Get(key+":c", 0)
		if err != nil {
			return lerror.Wrap(err, "get session IpLimit error 1")
		}

		latest := get.Int()               // 最近创建的时间
		if now-latest < durationSeconds { // 在间隔时间之内

			get, err := r.Session.Get(key, 0)
			if err != nil {
				return lerror.Wrap(err, "get session IpLimit error 2")
			}
			num := get.Int()
			if num >= maxNum {
				return errorx.New("请求过于频率，请稍候再试")
			}
			_ = r.Session.Set(key, strconv.Itoa(num+1))
			_ = r.Session.Set(key+":c", strconv.Itoa(now))
		} else {
			_ = r.Session.Set(key, "1")
			_ = r.Session.Set(key+":c", strconv.Itoa(now))
		}
	}
	return
}

func IpClient(r *ghttp.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// 请求 header
// Access-Control-Request-Method:  该字段的值对应当前请求类型，例如 GET、POST、PUT等等。浏览器会自动处理。
// Access-Control-Request-Headers: 该字段的值对应当前请求可能会携带的额外的自定义 header 字段名，如标识请求流水的 x-request-id，用于 Auth 鉴权的 Authorization 字段。
// 如果服务端支持该跨域请求，建议返回 204 状态码（返回 200 也可以
// 响应 header
// Access-Control-Allow-Origin: 允许哪些域被允许跨域，例如 http://qq.com 或 https://qq.com，或者设置为 * ，即允许所有域访问（通常见于 CDN ）
// Access-Control-Allow-Credentials: 是否携带票据访问（对应 fetch 方法中 credentials），当该值为 true 时，Access-Control-Allow-Origin 不允许设置为 *
// 跨域如果想带 cookie，Access-Control-Allow-Origin就不能设置为*，需要指定具体域名。换言之：Access-Control-Allow-Credentials: true 和 Access-Control-Allow-Origin: *不能同时使用。
// Access-Control-Allow-Methods: 标识该资源支持哪些方法，例如：POST, GET, PUT, DELETE
// Access-Control-Allow-Headers: 标识允许哪些额外的自定义 header 字段和非简单值的字段
// Access-Control-Expose-Headers: 通过该字段指出哪些额外的 header 可以被支持。正常情况下只能读取 Cache-Control/Content-Language/Content-Type/Expires/Last-Modified/Pragma
// Access-Control-Max-Age: 表示可以缓存 Access-Control-Allow-Methods 和 Access-Control-Allow-Headers 提供的信息多长时间，单位秒，一般为10分钟。

// 节省 OPTIONS
// 方法一：服务器端设置 Access-Control-Max-Age 字段
// 该缓存只针对这一个请求 URL 和相同的 header，无法针对整个域或者模糊匹配 URL 做缓存。
/*
server {
	listen 80;
	listen [::]:80;

	# 使用 xxx.test chrome 才不会自动跳转 https
  	server_name pp.test;
	# 前端如果自行设置了 credentials:true，则会导致请求失败一半

	add_header 'Access-Control-Allow-Origin' '*';
	add_header 'Access-Control-Allow-Credentials' 'true';
	add_header 'Access-Control-Allow-Headers' 'Authorization,Accept,Origin,DNT,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range, X-CSRF-Token';
	add_header 'Access-Control-Allow-Methods' 'GET,POST,OPTIONS,PUT,DELETE,PATCH';
	add_header 'Access-Control-Max-Age' 600;


	location ^~ /api/base/ {
		if ($request_method = 'OPTIONS') {
			return 204;
		}
		proxy_pass  http://127.0.0.1:8010; # 注意不要以 / 结尾防止拼接错误
	}

}
*/
// Deprecated: 使用 nginx proxy 代替
func Cors(s *ghttp.Server) {
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		option := r.Response.DefaultCORSOptions()
		r.Response.CORS(option)
		if r.Method == "OPTIONS" {
			r.Response.WriteStatus(http.StatusNoContent)
			r.ExitAll()
		}
		r.Middleware.Next()
	})
}

// https://goframe.org/pages/viewpage.action?pageId=1114483
// 通常用于 PUT/DELETE 请求，要求提交的是 json 字符串
// 在 gf 中，对于复杂的数据类型（结构嵌套, map[string]interface{} 嵌套），必须以 json 方式提交
// 可以使用 gjson.Encode(marchData)，在测试中 ghttp.Client 会自动添加 application/json
// https://goframe.org/net/ghttp/client
func RequestData(r *ghttp.Request, v interface{}) (err error) {
	if "GET" == r.Request.Method {
		return r.GetQueryStruct(v)
	}
	// 被读过两次后就没有数据了
	// https://goframe.org/pages/viewpage.action?pageId=1114483
	ct := r.Header.Get("Content-Type")
	switch ct {
	case datax.MIMEApplicationForm,
		datax.MIMEApplicationFormData,
		datax.MIMEMultipartForm,
		datax.MIMEMultipartMixed:
		return r.GetFormStruct(v)
	case datax.MIMEApplicationJSON:
		return r.Parse(v)
	default:
		return r.GetRequestStruct(v)
	}
}

func Request(r *ghttp.Request, v interface{}) {
	if err := RequestData(r, v); err != nil {
		logx.Debug("request data error:", err)
	}
}

// 读取参数并对其进行验证，通常用在实现了 Model 的模型上
// 注意：因为赋值，所以 pd 必须是一个 point, 例如 &param.Page{}
func RInvalid(r *ghttp.Request, pd interface{}, isCreate bool) {
	if !golibrary.IsPointer(pd) {
		ResponseIfError(r, errors.New("RInvalid param failed, not a pointer"))
	}

	Request(r, pd)
	//glog.Debug("RInvalid:", pd)
	err := pd.(appx.ParamInvalid).Invalid(isCreate)
	ResponseIfError(r, err)
}

func Hostname(r *ghttp.Request) (host string) {
	host = r.Request.Host
	colon := strings.LastIndexByte(host, ':')
	if colon != -1 && validOptionalPort(host[colon:]) {
		//host, port = host[:colon], host[colon+1:]
		host = host[:colon]
	}

	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}

	return
}

// from url
func validOptionalPort(port string) bool {
	if port == "" {
		return true
	}
	if port[0] != ':' {
		return false
	}
	for _, b := range port[1:] {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

func MustNotEmptyInt64(r *ghttp.Request, name string) (id int64) {
	q := r.Get(name).Int64()
	if q < 1 {
		response(r, http.StatusBadRequest, fmt.Sprintf("参数 %s 错误.", name), nil, false)
		r.ExitAll()
	}
	return q
}

//
//  MustNotEmptyString
//  @Description: 获取查询参数
//  @param r
//  @param name key 名
//  @param strict 如果为 true，则将 null, undefined 作为 空错误
//  @return string
//
func MustNotEmptyString(r *ghttp.Request, name string, strict bool) string {
	q := r.Get(name).String()
	if q == "" || strings.TrimSpace(q) == "" {
		response(r, http.StatusBadRequest, fmt.Sprintf("参数 %s 错误..", name), nil, false)
		r.ExitAll()
	}
	if strict {
		switch q {
		case "null", "undefined":
			response(r, http.StatusBadRequest, fmt.Sprintf("参数 %s 不能为空", name), nil, false)
			r.ExitAll()
		}
	}
	return q
}

func ChangeStatus(r *ghttp.Request) datax.ChangeStatus {
	data := datax.ChangeStatus{}
	if data.Id == 0 && data.Sid != "" {
		data.Id = appx.ExplodePrefixNum(data.Sid)
	}
	if data.Id < 1 {
		ResponseError(r, "id 不能为空")
	}
	return data
}
