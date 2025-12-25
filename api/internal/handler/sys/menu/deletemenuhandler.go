package menu

import (
	"net/http"

	"idrm/api/internal/logic/sys/menu"
	"idrm/api/internal/svc"
	"idrm/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// DeleteMenuHandler 删除菜单处理器
type DeleteMenuHandler struct {
	svcCtx *svc.ServiceContext
}

// NewDeleteMenuHandler 创建DeleteMenuHandler实例
func NewDeleteMenuHandler(svcCtx *svc.ServiceContext) *DeleteMenuHandler {
	return &DeleteMenuHandler{
		svcCtx: svcCtx,
	}
}

// DeleteMenu 处理删除菜单请求
func (h *DeleteMenuHandler) DeleteMenu(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	var req types.DeleteMenuReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 2. 调用Logic层处理业务逻辑
	l := menu.NewDeleteMenuLogic(r.Context(), h.svcCtx)
	resp, err := l.DeleteMenu(&req)

	// 3. 返回响应
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
