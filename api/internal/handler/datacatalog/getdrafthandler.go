// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package datacatalog

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"idrm/api/internal/logic/datacatalog"
	"idrm/api/internal/svc"
	"idrm/api/internal/types"
)

// 获取草稿详情
func GetDraftHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDraftReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := datacatalog.NewGetDraftLogic(r.Context(), svcCtx)
		resp, err := l.GetDraft(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
