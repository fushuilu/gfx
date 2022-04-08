package bootx

import (
	"fmt"

	"github.com/fushuilu/gfx/cx"
	"github.com/fushuilu/golibrary"
	"github.com/fushuilu/golibrary/lerror"

	"github.com/gogf/gf/v2/os/glog"
)

func Start() error {
	// 读取开发模式
	mode, err := cx.GetString("mode")
	if err != nil {
		return lerror.Wrap(err, "read config.mode failed")
	}
	logPath, err := cx.GetString("log.path")
	if err != nil {
		return lerror.Wrap(err, "read config.log.path failed")
	}

	switch mode {
	case "dev", "local", "debug":
		cx.SetDebug(true)
		fmt.Println("run app in develop mode")

		glog.SetStdoutPrint(true)
		glog.SetLevel(glog.LEVEL_ALL)
	default:
		fmt.Println("run app in product mode")

		glog.SetStdoutPrint(false)
		glog.SetLevel(glog.LEVEL_PROD)
	}
	// 日志
	if logPath != "" {
		var exist bool
		if exist, err = golibrary.FileExist(logPath); err != nil {
			return lerror.Wrap(err, "检测日志目录路径时错误")
		} else if exist {
			return glog.SetPath(logPath)
		}
	}

	return nil
}
