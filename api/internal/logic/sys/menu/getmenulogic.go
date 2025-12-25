package menu

import (
	"context"
	"fmt"

	"idrm/api/internal/svc"
	"idrm/api/internal/types"
	"idrm/model/system/menu"
	"idrm/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetMenuLogic 获取菜单详情业务逻辑
type GetMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetMenuLogic 创建GetMenuLogic实例
func NewGetMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuLogic {
	return &GetMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetMenu 获取菜单详情
func (l *GetMenuLogic) GetMenu(req *types.GetMenuReq) (*types.GetMenuResp, error) {
	// 1. 参数验证
	if req.Id <= 0 {
		return nil, errorx.New(errorx.ErrCodeParamMissing, "菜单ID不能为空")
	}

	// 2. 调用Model层查询菜单
	menuData, err := l.svcCtx.MenuModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == menu.ErrNotFound {
			l.Errorw("menu not found",
				logx.Field("menuId", req.Id),
			)
			return nil, errorx.NewWithMsg(errorx.ErrCodeNotFound, "菜单不存在")
		}
		l.Errorw("failed to find menu",
			logx.Field("menuId", req.Id),
			logx.Field("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to find menu: %w", err)
	}

	// 3. 记录日志
	l.Infow("menu retrieved successfully",
		logx.Field("menuId", menuData.Id),
		logx.Field("menuName", menuData.Name),
	)

	// 4. 数据转换 (model -> types)
	return l.convertToGetMenuResp(menuData), nil
}

// convertToGetMenuResp 将Model转换为Response
func (l *GetMenuLogic) convertToGetMenuResp(menuData *menu.Menu) *types.GetMenuResp {
	return &types.GetMenuResp{
		Id:            menuData.Id,
		ParentId:      menuData.ParentId,
		Name:          menuData.Name,
		RoutePath:     menuData.RoutePath,
		ComponentPath: menuData.ComponentPath,
		Icon:          menuData.Icon,
		SortOrder:     menuData.SortOrder,
		PermTag:       menuData.PermTag,
		Status:        menuData.Status,
		CreatedAt:     menuData.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     menuData.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
