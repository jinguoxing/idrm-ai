package menu

import (
	"net/http"

	"idrm/api/internal/logic/sys/menu"
	"idrm/api/internal/svc"
	"idrm/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// UpdateMenuHandler 更新菜单处理器
type UpdateMenuHandler struct {
	svcCtx *svc.ServiceContext
}

// NewUpdateMenuHandler 创建UpdateMenuHandler实例
func NewUpdateMenuHandler(svcCtx *svc.ServiceContext) *UpdateMenuHandler {
	return &UpdateMenuHandler{
		svcCtx: svcCtx,
	}
}

// UpdateMenu 处理更新菜单请求
func (h *UpdateMenuHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	var req types.UpdateMenuReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 2. 调用Logic层处理业务逻辑
	l := menu.NewUpdateMenuLogic(r.Context(), h.svcCtx)
	resp, err := l.UpdateMenu(&req)

	// 3. 返回响应
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
