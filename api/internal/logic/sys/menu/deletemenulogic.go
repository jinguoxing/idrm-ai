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

// DeleteMenuLogic 删除菜单业务逻辑
type DeleteMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteMenuLogic 创建DeleteMenuLogic实例
func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteMenu 删除菜单
// 业务规则：
// 1. 菜单必须存在
// 2. 如果存在子菜单，则不允许删除（防止孤立节点）
func (l *DeleteMenuLogic) DeleteMenu(req *types.DeleteMenuReq) (*types.DeleteMenuResp, error) {
	// 1. 参数验证
	if req.Id <= 0 {
		return nil, errorx.New(errorx.ErrCodeParamMissing, "菜单ID不能为空")
	}

	// 2. 检查菜单是否存在
	existingMenu, err := l.svcCtx.MenuModel.FindOne(l.ctx, req.Id)
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

	// 3. 检查是否有子菜单
	hasChildren, err := l.svcCtx.MenuModel.HasChildren(l.ctx, req.Id)
	if err != nil {
		l.Errorw("failed to check children",
			logx.Field("menuId", req.Id),
			logx.Field("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to check children: %w", err)
	}

	if hasChildren {
		l.Errorw("cannot delete menu with children",
			logx.Field("menuId", req.Id),
		)
		return nil, menu.NewMenuHasSubEntriesError(req.Id)
	}

	// 4. 调用Model层删除数据（软删除）
	if err := l.svcCtx.MenuModel.Delete(l.ctx, req.Id); err != nil {
		l.Errorw("failed to delete menu",
			logx.Field("menuId", req.Id),
			logx.Field("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to delete menu: %w", err)
	}

	// 5. 记录日志
	l.Infow("menu deleted successfully",
		logx.Field("menuId", req.Id),
		logx.Field("menuName", existingMenu.Name),
	)

	// 6. 返回结果
	return &types.DeleteMenuResp{
		Id: req.Id,
	}, nil
}
