package menu

import (
	"net/http"

	"idrm/api/internal/logic/sys/menu"
	"idrm/api/internal/svc"
	"idrm/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetMenuHandler 获取菜单详情处理器
type GetMenuHandler struct {
	svcCtx *svc.ServiceContext
}

// NewGetMenuHandler 创建GetMenuHandler实例
func NewGetMenuHandler(svcCtx *svc.ServiceContext) *GetMenuHandler {
	return &GetMenuHandler{
		svcCtx: svcCtx,
	}
}

// GetMenu 处理获取菜单详情请求
func (h *GetMenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	var req types.GetMenuReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 2. 调用Logic层处理业务逻辑
	l := menu.NewGetMenuLogic(r.Context(), h.svcCtx)
	resp, err := l.GetMenu(&req)

	// 3. 返回响应
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
