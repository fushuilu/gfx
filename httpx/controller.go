package httpx

import (
	"github.com/fushuilu/golibrary/appx"
	"github.com/fushuilu/golibrary/appx/datax"
	"github.com/gogf/gf/v2/net/ghttp"
)

func ActionChangeStatus(r *ghttp.Request) (data datax.ChangeStatus) {
	Request(r, &data)
	if data.Sid != "" {
		data.Id = appx.ExplodePrefixNum(data.Sid)
	}
	if data.Id < 1 {
		ResponseError(r, "id is empty")
	}
	data.StatusIndex = datax.MapStatusText.GetIntValue(data.Status)
	return
}

func ActionChangeState(r *ghttp.Request) (data datax.ChangeState) {
	Request(r, &data)

	if data.Sid != "" {
		data.Id = appx.ExplodePrefixNum(data.Sid)
	}

	if data.Id < 1 {
		ResponseError(r, "id is empty")
	}
	data.StateIndex = datax.MapStatusText.GetIntValue(data.State)
	return
}
