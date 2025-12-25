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

// CreateMenuLogic 创建菜单业务逻辑
type CreateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateMenuLogic 创建CreateMenuLogic实例
func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateMenu 创建菜单
// 业务规则：
// 1. 菜单名称不能为空
// 2. 父级菜单必须存在（parentId为0时表示根菜单）
// 3. 排序号默认为0
func (l *CreateMenuLogic) CreateMenu(req *types.CreateMenuReq) (*types.CreateMenuResp, error) {
	// 1. 参数验证
	if err := l.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// 2. 验证父级菜单是否存在
	if req.ParentId != 0 {
		parentMenu, err := l.svcCtx.MenuModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			if err == menu.ErrNotFound {
				l.Errorw("parent menu not found",
					logx.Field("parentId", req.ParentId),
				)
				return nil, menu.NewParentMenuNotFoundError(req.ParentId)
			}
			l.Errorw("failed to find parent menu",
				logx.Field("parentId", req.ParentId),
				logx.Field("error", err.Error()),
			)
			return nil, fmt.Errorf("failed to find parent menu: %w", err)
		}
		_ = parentMenu // 确认父级菜单存在
	}

	// 3. 数据转换 (types -> model)
	menuData := &menu.Menu{
		ParentId:      req.ParentId,
		Name:          req.Name,
		RoutePath:     req.RoutePath,
		ComponentPath: req.ComponentPath,
		Icon:          req.Icon,
		SortOrder:     req.SortOrder,
		PermTag:       req.PermTag,
		Status:        req.Status,
	}

	// 4. 调用Model层插入数据
	result, err := l.svcCtx.MenuModel.Insert(l.ctx, menuData)
	if err != nil {
		l.Errorw("failed to create menu",
			logx.Field("name", req.Name),
			logx.Field("parentId", req.ParentId),
			logx.Field("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to create menu: %w", err)
	}

	// 5. 记录日志
	l.Infow("menu created successfully",
		logx.Field("menuId", result.Id),
		logx.Field("menuName", result.Name),
		logx.Field("parentId", result.ParentId),
	)

	// 6. 返回结果 (model -> types)
	return &types.CreateMenuResp{
		Id:   result.Id,
		Name: result.Name,
	}, nil
}

// validateCreateRequest 验证创建请求参数
func (l *CreateMenuLogic) validateCreateRequest(req *types.CreateMenuReq) error {
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

	// 权限标识长度验证
	if len(req.PermTag) > 100 {
		return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "权限标识长度不能超过100个字符")
	}

	// 路由路径长度验证
	if len(req.RoutePath) > 200 {
		return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "路由路径长度不能超过200个字符")
	}

	// 组件路径长度验证
	if len(req.ComponentPath) > 200 {
		return errorx.NewWithMsg(errorx.ErrCodeParamInvalid, "组件路径长度不能超过200个字符")
	}

	return nil
}
