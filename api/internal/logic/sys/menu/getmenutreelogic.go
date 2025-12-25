package menu

import (
	"context"
	"fmt"

	"idrm/api/internal/svc"
	"idrm/api/internal/types"
	"idrm/model/system/menu"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetMenuTreeLogic 获取菜单树业务逻辑
type GetMenuTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetMenuTreeLogic 创建GetMenuTreeLogic实例
func NewGetMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuTreeLogic {
	return &GetMenuTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetMenuTree 获取菜单树
// 业务规则：
// 1. 支持状态过滤（status参数可选）
// 2. 返回完整的树形结构
// 3. 使用Redis缓存（TTL: 1小时）
func (l *GetMenuTreeLogic) GetMenuTree(req *types.GetMenuTreeReq) (*types.GetMenuTreeResp, error) {
	// 1. 构建缓存key（暂未使用）
	_ = l.buildCacheKey(req.Status)

	// 2. 尝试从Redis缓存获取（如果配置了Redis）
	// TODO: 实现Redis缓存逻辑
	// cachedTree, err := l.getFromCache(cacheKey)
	// if err == nil && cachedTree != nil {
	//     return cachedTree, nil
	// }

	// 3. 从数据库查询所有菜单
	status := -1 // 默认查询所有状态
	if req.Status != 0 {
		status = req.Status
	}

	allMenus, err := l.svcCtx.MenuModel.FindAll(l.ctx, status)
	if err != nil {
		l.Errorw("failed to find all menus",
			logx.Field("status", status),
			logx.Field("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to find all menus: %w", err)
	}

	// 4. 构建树形结构
	tree := l.buildMenuTree(allMenus, 0)

	// 5. 记录日志
	l.Infow("menu tree retrieved successfully",
		logx.Field("totalMenus", len(allMenus)),
		logx.Field("rootNodes", len(tree)),
		logx.Field("status", status),
	)

	// 6. 缓存结果（如果配置了Redis）
	// TODO: 实现Redis缓存逻辑
	// _ = l.saveToCache(cacheKey, tree)

	// 7. 返回结果
	return &types.GetMenuTreeResp{
		Tree: tree,
	}, nil
}

// buildMenuTree 构建菜单树（递归算法）
func (l *GetMenuTreeLogic) buildMenuTree(menus []*menu.Menu, parentId int64) []types.MenuTreeNode {
	var tree []types.MenuTreeNode

	// 遍历所有菜单，找到属于当前父级的子菜单
	for _, m := range menus {
		if m.ParentId == parentId {
			node := l.convertToMenuTreeNode(m)
			// 递归查找子菜单
			node.Children = l.buildMenuTree(menus, m.Id)
			tree = append(tree, node)
		}
	}

	return tree
}

// convertToMenuTreeNode 将Model转换为树节点
func (l *GetMenuTreeLogic) convertToMenuTreeNode(menuData *menu.Menu) types.MenuTreeNode {
	return types.MenuTreeNode{
		Id:            menuData.Id,
		ParentId:      menuData.ParentId,
		Name:          menuData.Name,
		RoutePath:     menuData.RoutePath,
		ComponentPath: menuData.ComponentPath,
		Icon:          menuData.Icon,
		SortOrder:     menuData.SortOrder,
		PermTag:       menuData.PermTag,
		Status:        menuData.Status,
		Children:      []types.MenuTreeNode{}, // 初始化空切片
	}
}

// buildCacheKey 构建缓存key
func (l *GetMenuTreeLogic) buildCacheKey(status int) string {
	if status == 0 {
		return "menu:tree:all"
	}
	return fmt.Sprintf("menu:tree:status:%d", status)
}

// getFromCache 从Redis缓存获取（待实现）
func (l *GetMenuTreeLogic) getFromCache(key string) (*types.GetMenuTreeResp, error) {
	// TODO: 实现Redis缓存读取
	// if l.svcCtx.Redis == nil {
	//     return nil, errors.New("redis not configured")
	// }
	//
	// data, err := l.svcCtx.Redis.Get(key)
	// if err != nil {
	//     return nil, err
	// }
	//
	// var resp types.GetMenuTreeResp
	// if err := json.Unmarshal([]byte(data), &resp); err != nil {
	//     return nil, err
	// }
	//
	// return &resp, nil
	return nil, nil
}

// saveToCache 保存到Redis缓存（待实现）
func (l *GetMenuTreeLogic) saveToCache(key string, tree *types.GetMenuTreeResp) error {
	// TODO: 实现Redis缓存写入
	// if l.svcCtx.Redis == nil {
	//     return nil // Redis未配置，不报错
	// }
	//
	// data, err := json.Marshal(tree)
	// if err != nil {
	//     return err
	// }
	//
	// // 设置1小时过期时间
	// return l.svcCtx.Redis.Setex(key, string(data), int(time.Hour.Seconds()))
	return nil
}
