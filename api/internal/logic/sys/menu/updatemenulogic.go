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

// UpdateMenuLogic 更新菜单业务逻辑
type UpdateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateMenuLogic 创建UpdateMenuLogic实例
func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateMenu 更新菜单
// 业务规则：
// 1. 菜单必须存在
// 2. 父级菜单必须存在（parentId为0时表示根菜单）
// 3. 不能将自己设置为父级（防止循环引用）
func (l *UpdateMenuLogic) UpdateMenu(req *types.UpdateMenuReq) (*types.UpdateMenuResp, error) {
	// 1. 参数验证
	if err := l.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// 2. 检查菜单是否存在
	_, err := l.svcCtx.MenuModel.FindOne(l.ctx, req.Id)
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

	// 3. 验证父级菜单
	if req.ParentId != 0 {
		// 不能将自己设置为父级
		if req.ParentId == req.Id {
			l.Errorw("cannot set self as parent",
				logx.Field("menuId", req.Id),
			)
			return nil, errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "不能将自己设置为父级菜单")
		}

		// 检查父级菜单是否存在
		parentMenu, err := l.svcCtx.MenuModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			if err == menu.ErrNotFound {
				l.Errorw("parent menu not found",
					logx.Field("parentId", req.ParentId),
				)
				return nil, menu.NewParentMenuNotFoundError(req.ParentId)
			}
			return nil, fmt.Errorf("failed to find parent menu: %w", err)
		}
		_ = parentMenu // 确认父级菜单存在
	}

	// 4. 数据转换 (types -> model)
	menuData := &menu.Menu{
		Id:            req.Id,
		ParentId:      req.ParentId,
		Name:          req.Name,
		RoutePath:     req.RoutePath,
		ComponentPath: req.ComponentPath,
		Icon:          req.Icon,
		SortOrder:     req.SortOrder,
		PermTag:       req.PermTag,
		Status:        int(req.Status), // int8 -> int
	}

	// 5. 调用Model层更新数据
	if err := l.svcCtx.MenuModel.Update(l.ctx, menuData); err != nil {
		l.Errorw("failed to update menu",
			logx.Field("menuId", req.Id),
			logx.Field("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to update menu: %w", err)
	}

	// 6. 记录日志
	l.Infow("menu updated successfully",
		logx.Field("menuId", menuData.Id),
		logx.Field("menuName", menuData.Name),
		logx.Field("parentId", menuData.ParentId),
	)

	// 7. 返回结果
	return &types.UpdateMenuResp{
		Id: menuData.Id,
	}, nil
}

// validateUpdateRequest 验证更新请求参数
func (l *UpdateMenuLogic) validateUpdateRequest(req *types.UpdateMenuReq) error {
	if req.Id <= 0 {
		return errorx.New(errorx.ErrCodeParamMissing, "菜单ID不能为空")
	}

	if req.Name == "" {
		return errorx.New(errorx.ErrCodeParamMissing, "菜单名称不能为空")
	}

	if len(req.Name) > 50 {
		return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "菜单名称长度不能超过50个字符")
	}

	if req.ParentId < 0 {
		return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "父级菜单ID不能为负数")
	}

	if req.Status != 0 && req.Status != 1 {
		return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "状态值必须为0或1")
	}

	return nil
}
