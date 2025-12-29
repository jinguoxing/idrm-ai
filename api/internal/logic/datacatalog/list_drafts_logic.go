package datacatalog

import (
	"context"
	"fmt"

	"idrm/api/internal/svc"
	"idrm/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ListDraftsLogic 草稿列表
type ListDraftsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListDraftsLogic 创建ListDraftsLogic实例
func NewListDraftsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDraftsLogic {
	return &ListDraftsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListDrafts 查询草稿列表
func (l *ListDraftsLogic) ListDrafts(req *types.ListDraftsReq) (resp *types.ListDraftsResp, err error) {
	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 查询草稿列表
	drafts, total, err := l.svcCtx.DataCatalogModel.FindDrafts(l.ctx, req.SourceDeptId, req.Status, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("查询草稿列表失败: %w", err)
	}

	// 构建响应
	items := make([]types.GetDraftResp, 0, len(drafts))
	for _, draft := range drafts {
		items = append(items, types.GetDraftResp{
			CatalogId:       draft.Id,
			Title:           draft.Title,
			Description:     draft.Description,
			SourceDeptId:    draft.SourceDeptId,
			SourceDeptName:  draft.OwnerName,
			DataDomain:      draft.DataDomain,
			PublishStatus:   string(draft.PublishStatus),
			CreatedAt:       draft.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       formatTimePtr(draft.UpdatedAt),
		})
	}

	return &types.ListDraftsResp{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
