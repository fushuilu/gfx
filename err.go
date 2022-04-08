package gfx

import (
	"github.com/fushuilu/golibrary/lerror"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gvalid"
)

// 是否为错误类型
func IsError(resp interface{}) bool {
	switch resp.(type) {
	case error,
		*gerror.Error, gerror.Error,
		*gvalid.Error, gvalid.Error,
		*lerror.Error, lerror.Error:
		return true
	default:
		return false
	}
}
