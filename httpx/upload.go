package httpx

import (
	"errors"

	"github.com/gogf/gf/v2/net/ghttp"
)

// 随机文件名
func UploadOne(r *ghttp.Request, path string, invalid func(file *ghttp.UploadFile) error) (filename string, err error) {

	file := r.GetUploadFile("file")
	if file == nil {
		return "", errors.New("上传文件[file]为空")
	}
	if err = invalid(file); err != nil {
		return
	}

	filename, err = file.Save(path, true)
	return
}
