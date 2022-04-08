package httpx

import (
	"github.com/fushuilu/golibrary/appx/datax"

	"github.com/gogf/gf/v2/net/ghttp"
)

// 标准返回结果数据结构封装。
func Json(r *ghttp.Request, code int, message string, d ...interface{}) {
	responseData := interface{}(nil)
	if len(d) > 0 {
		responseData = d[0]
	}
	_ = r.Response.WriteJson(datax.CodeResponse{
		Code:    code,
		Message: message,
		Data:    responseData,
	})
}

// 返回JSON数据并退出当前HTTP执行函数。
func JsonExit(r *ghttp.Request, err int, msg string, data ...interface{}) {
	Json(r, err, msg, data...)
	r.Exit()
}
