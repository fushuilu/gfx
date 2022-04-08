package httpx

import (
	"github.com/fushuilu/golibrary/appx/db"
	"github.com/gogf/gf/v2/net/ghttp"
)

// limit ? offset ?
func GetPagination(r *ghttp.Request) (pagination db.Pagination) {
	return QueryPagination(r)
}

func QueryPagination(r *ghttp.Request) (pag db.Pagination) {
	_ = r.GetQueryStruct(&pag)
	if pag.Limit() > db.MaxPageSize {
		pag.LimitSize = db.MaxPageSize
	} else if pag.Limit() < 1 {
		pag.LimitSize = db.DefaultPageSize
	}
	if pag.Page > 0 {
		pag.PageIndex = pag.Page
	}
	if pag.PageIndex < 0 {
		pag.PageIndex = 0
	}
	return
}
