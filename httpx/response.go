package httpx

import (
	"fmt"
	"github.com/fushuilu/gfx"
	"net/http"

	"github.com/fushuilu/gfx/cx"
	"github.com/fushuilu/gfx/logx"
	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/fushuilu/golibrary/lerror"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
)

func response(r *ghttp.Request, httpCode int, errMsg string, body interface{}, isJson bool) {
	hasErr := errMsg != ""

	if cx.IsDebug() {
		logx.Debug("--> [", r.Method, "]", httpCode, r.RequestURI)
		if hasErr {
			logx.Debug("errMsg", errMsg)
		}
		logx.Debug("body", body)
	}

	if httpCode == 0 && hasErr {
		httpCode = http.StatusBadRequest
	}

	if r.Get("error").String() == "0" {
		r.Response.Status = http.StatusOK
		_ = r.Response.WriteJson(datax.CodeResponse{
			Code:    httpCode, // 业务码
			Message: errMsg,
			Data:    body,
		})
	} else {
		r.Response.WriteHeader(httpCode)
		if hasErr {
			r.Response.Write(errMsg)
		} else if isJson {
			_ = r.Response.WriteJson(body)
		} else {
			r.Response.Write(body)
		}
	}
	if hasErr { // ??
		r.Response.CORSDefault()
	}
	r.ExitAll()
	fmt.Println("stop ...")
}

//
//  responseLError
//  @Description: 处理 lerror 类型错误
//
func responseLError(r *ghttp.Request, code int, v *lerror.Error) {
	if v.IsDebug() || cx.IsDebug() {
		logx.Debug(v.Stack())
	} else {
		logx.Error(v.Stack()) // 记录日志
	}
	if v.Code() > 0 {
		code = v.Code()
	}
	response(r, code, v.Error(), nil, false)
}

//
//  responseGError
//  @Description: 处理 gerror 类型错误
//
func responseGError(r *ghttp.Request, code int, v *gerror.Error) {
	if cx.IsDebug() {
		logx.Debug(v.Stack())
	} else {
		logx.Error(v.Stack())
	}
	response(r, code, v.Error(), nil, false)
}

// 用于返回错误信息，并写入日志
func SmartResponseError(r *ghttp.Request, code int, data interface{}) {
	if code == 0 {
		code = http.StatusBadRequest
	}

	switch v := data.(type) {
	case *lerror.Error:
		responseLError(r, code, v)
	case lerror.Error:
		responseLError(r, code, &v)
	case *gerror.Error:
		responseGError(r, code, v)
	case gerror.Error:
		responseGError(r, code, &v)
	case *gvalid.Error:
		response(r, code, (*v).Current().Error(), nil, false)
	case gvalid.Error:
		response(r, code, v.Current().Error(), nil, false)
	case error:
		response(r, code, v.Error(), nil, false)
	default:
		response(r, code, gconv.String(data), nil, false)
	}
}

//
//  SmartResponse
//  @Description: 智能响应
//  @param r
//  @param data "响应信息内容body"，"响应码 httpCode"
//
func SmartResponse(r *ghttp.Request, d interface{}) {
	if d == nil {
		response(r, 0, "", "", false)
		return
	}
	if d == datax.Exit {
		r.Exit()
		return
	}
	if gfx.IsError(d) {
		SmartResponseError(r, 400, d)
		return
	}

	switch d.(type) {
	case string, bool, int, int64, uint, uint64, float32, float64:
		response(r, 0, "", d, false)
	case datax.Location307Response:
		if cx.IsDebug() {
			logx.Debug(">>>>>> Location", d)
		}
		res := d.(datax.Location307Response)
		r.Response.Header().Set("Location", res.Location)
		if res.AccessControlAllowOrigin != "" {
			r.Response.Header().Set("Access-Control-Allow-Origin", res.AccessControlAllowOrigin)
		}

		r.Response.WriteHeader(http.StatusTemporaryRedirect)
	default:
		response(r, 0, "", d, true)
	}
}
