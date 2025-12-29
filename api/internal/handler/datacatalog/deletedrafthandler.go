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

// 删除草稿
func DeleteDraftHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteDraftReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := datacatalog.NewDeleteDraftLogic(r.Context(), svcCtx)
		resp, err := l.DeleteDraft(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
