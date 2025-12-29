package menu

import (
	"context"
	"fmt"

	"idrm/api/internal/svc"
	"idrm/api/internal/types"
	"idrm/model/system/menu"

	"github.com/zeromicro/go-zero/core/logx"
)

// ListMenusLogic 获取菜单列表业务逻辑
type ListMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListMenusLogic 创建ListMenusLogic实例
func NewListMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenusLogic {
	return &ListMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListMenus 获取菜单列表
// 业务规则：
// 1. 支持按父级ID过滤（parentId参数可选）
// 2. 支持按状态过滤（status参数可选）
// 3. 返回分页结果
func (l *ListMenusLogic) ListMenus(req *types.ListMenusReq) (*types.ListMenusResp, error) {
	// 1. 参数验证
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 2. 查询菜单列表
	var menus []*menu.Menu
	var err error

	if req.ParentId != 0 {
		// 按父级ID查询
		menus, err = l.svcCtx.MenuModel.FindByParentId(l.ctx, req.ParentId)
	} else {
		// 查询所有菜单（带状态过滤）
		status := -1
		if req.Status != 0 {
			status = int(req.Status)
		}
		menus, err = l.svcCtx.MenuModel.FindAll(l.ctx, status)
	}

	if err != nil {
		l.Errorw("failed to find menus",
			logx.Field("parentId", req.ParentId),
			logx.Field("status", req.Status),
			logx.Field("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to find menus: %w", err)
	}

	// 3. 计算总数
	total := int64(len(menus))

	// 4. 分页处理（内存分页，因为Model层返回全部数据）
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	if start >= len(menus) {
		// 超出范围，返回空列表
		l.Infow("no menus found for pagination",
			logx.Field("page", req.Page),
			logx.Field("pageSize", req.PageSize),
		)
		return &types.ListMenusResp{
			Items:    []*types.GetMenuResp{},
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		}, nil
	}

	if end > len(menus) {
		end = len(menus)
	}

	pagedMenus := menus[start:end]

	// 5. 数据转换 (model -> types)
	items := make([]*types.GetMenuResp, 0, len(pagedMenus))
	for _, m := range pagedMenus {
		items = append(items, &types.GetMenuResp{
			Id:            m.Id,
			ParentId:      m.ParentId,
			Name:          m.Name,
			RoutePath:     m.RoutePath,
			ComponentPath: m.ComponentPath,
			Icon:          m.Icon,
			SortOrder:     m.SortOrder,
			PermTag:       m.PermTag,
			Status:        int8(m.Status), // int -> int8
			CreatedAt:     m.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     m.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 6. 记录日志
	l.Infow("menus listed successfully",
		logx.Field("total", total),
		logx.Field("page", req.Page),
		logx.Field("pageSize", req.PageSize),
		logx.Field("returned", len(items)),
	)

	// 7. 返回结果
	return &types.ListMenusResp{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
