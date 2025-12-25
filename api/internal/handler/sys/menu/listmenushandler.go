package menu

import (
	"net/http"

	"idrm/api/internal/logic/sys/menu"
	"idrm/api/internal/svc"
	"idrm/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ListMenusHandler 获取菜单列表处理器
type ListMenusHandler struct {
	svcCtx *svc.ServiceContext
}

// NewListMenusHandler 创建ListMenusHandler实例
func NewListMenusHandler(svcCtx *svc.ServiceContext) *ListMenusHandler {
	return &ListMenusHandler{
		svcCtx: svcCtx,
	}
}

// ListMenus 处理获取菜单列表请求
func (h *ListMenusHandler) ListMenus(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	var req types.ListMenusReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 2. 调用Logic层处理业务逻辑
	l := menu.NewListMenusLogic(r.Context(), h.svcCtx)
	resp, err := l.ListMenus(&req)

	// 3. 返回响应
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
