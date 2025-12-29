package system

import (
	"net/http"

	"idrm/api/internal/logic/system"
	"idrm/api/internal/svc"
	"idrm/pkg/response"
)

// ORMMetricsHandler ORM指标处理器
type ORMMetricsHandler struct {
	svcCtx *svc.ServiceContext
}

// NewORMMetricsHandler 创建ORM指标处理器
func NewORMMetricsHandler(svcCtx *svc.ServiceContext) *ORMMetricsHandler {
	return &ORMMetricsHandler{
		svcCtx: svcCtx,
	}
}

// ORMMetrics 获取ORM性能指标
func (h *ORMMetricsHandler) ORMMetrics(w http.ResponseWriter, r *http.Request) {
	l := system.NewORMMetricsLogic(r.Context(), h.svcCtx)
	resp, err := l.ORMMetrics()
	if err != nil {
		response.Error(w, err)
	} else {
		response.Success(w, resp)
	}
}

// ResetORMMetrics 重置ORM性能指标
func (h *ORMMetricsHandler) ResetORMMetrics(w http.ResponseWriter, r *http.Request) {
	l := system.NewORMMetricsLogic(r.Context(), h.svcCtx)
	err := l.ResetORMMetrics()
	if err != nil {
		response.Error(w, err)
	} else {
		response.Success(w, "指标重置成功")
	}
}
