package httpx

import (
	"net/http"

	"github.com/fushuilu/gfx/cx"
	"github.com/fushuilu/gfx/logx"
	"github.com/fushuilu/golibrary/lerror"

	"github.com/gogf/gf/v2/net/ghttp"
)

// 控制器中使用，但是会导致无法测试
//
//  ResponseIfError
//  @Description: 响应错误的信息
//  @param err 真实的错误
//  @param msg 如果不为空，则将 msg[0] 信息返回给客户端
//
func ResponseIfError(r *ghttp.Request, err error, msg ...string) {
	if err != nil {
		if msg != nil {
			if e, ok := err.(*lerror.Error); ok {
				if cx.IsDebug() {
					logx.Debug(e.Stack())
				}
			}
			response(r, http.StatusBadRequest, msg[0], nil, false)
		} else {
			SmartResponseError(r, http.StatusBadRequest, err)
		}
	}
}

func ResponseError(r *ghttp.Request, msg string) {
	response(r, http.StatusBadRequest, msg, nil, false)
}

func ResponseAssertError(r *ghttp.Request, condition bool, msg string) {
	if condition {
		ResponseError(r, msg)
	}
}
