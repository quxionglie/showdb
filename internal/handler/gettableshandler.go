package handler

import (
	xhttp "github.com/zeromicro/x/http"
	"net/http"

	"github.com/quxionglie/showdb/internal/logic"
	"github.com/quxionglie/showdb/internal/svc"
	"github.com/quxionglie/showdb/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func getTablesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTablesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetTablesLogic(r.Context(), svcCtx)
		resp, err := l.GetTables(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
