package logx

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/fushuilu/gfx/cx"
	"github.com/fushuilu/golibrary"
)

var ctx = context.TODO()

// 打印调试信息
func Info(v ...interface{}) {
	if cx.IsDebug() {
		g.Log().Info(ctx, v...)
	}
}

// 打印调试信息及所在的方法信息
func DebugIfError(err error, v ...interface{}) {
	if err != nil && cx.IsDebug() {
		g.Log().Debug(ctx, v...)
		fmt.Println("error msg:", err)
		file, line, name := golibrary.Caller(2)
		fmt.Printf("\t%s:%d\n\t%s\n", file, line, name)
	}
}

func DebugWithStack(v ...interface{}) {
	if cx.IsDebug() {
		g.Log().Debug(ctx, v...)
		file, line, name := golibrary.Caller(2)
		fmt.Printf("\t%s:%d\n\t%s\n", file, line, name)
	}
}

func Debug(v ...interface{}) {
	if cx.IsDebug() {
		g.Log().Debug(ctx, v...)
	}
}
func Debugf(format string, v ...interface{}) {
	if cx.IsDebug() {
		g.Log().Debugf(ctx, format, v...)
	}
}

func Error(v ...interface{}) {
	g.Log().Error(ctx, v...)
}

func Warning(v ...interface{}) {
	g.Log().Warning(ctx, v...)
}
