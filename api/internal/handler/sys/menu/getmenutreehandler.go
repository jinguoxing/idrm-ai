package menu

import (
	"net/http"

	"idrm/api/internal/logic/sys/menu"
	"idrm/api/internal/svc"
	"idrm/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetMenuTreeHandler 获取菜单树处理器
type GetMenuTreeHandler struct {
	svcCtx *svc.ServiceContext
}

// NewGetMenuTreeHandler 创建GetMenuTreeHandler实例
func NewGetMenuTreeHandler(svcCtx *svc.ServiceContext) *GetMenuTreeHandler {
	return &GetMenuTreeHandler{
		svcCtx: svcCtx,
	}
}

// GetMenuTree 处理获取菜单树请求
func (h *GetMenuTreeHandler) GetMenuTree(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	var req types.GetMenuTreeReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 2. 调用Logic层处理业务逻辑
	l := menu.NewGetMenuTreeLogic(r.Context(), h.svcCtx)
	resp, err := l.GetMenuTree(&req)

	// 3. 返回响应
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
