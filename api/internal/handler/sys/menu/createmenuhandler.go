package menu

import (
	"net/http"

	"idrm/api/internal/logic/sys/menu"
	"idrm/api/internal/svc"
	"idrm/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// CreateMenuHandler 创建菜单处理器
type CreateMenuHandler struct {
	svcCtx *svc.ServiceContext
}

// NewCreateMenuHandler 创建CreateMenuHandler实例
func NewCreateMenuHandler(svcCtx *svc.ServiceContext) *CreateMenuHandler {
	return &CreateMenuHandler{
		svcCtx: svcCtx,
	}
}

// CreateMenu 处理创建菜单请求
func (h *CreateMenuHandler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	var req types.CreateMenuReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 2. 调用Logic层处理业务逻辑
	l := menu.NewCreateMenuLogic(r.Context(), h.svcCtx)
	resp, err := l.CreateMenu(&req)

	// 3. 返回响应
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
